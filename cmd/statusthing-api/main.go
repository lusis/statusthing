// Package main ...
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
	flag "github.com/spf13/pflag"

	"github.com/lusis/statusthing/gen/go/statusthing/v1/statusthingv1connect"
	"github.com/lusis/statusthing/internal/handlers"
	"github.com/lusis/statusthing/internal/services"
	"github.com/lusis/statusthing/internal/storers/memdb"

	"golang.org/x/exp/slog"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var (
	apiAddr        *string = flag.String("api-addr", "127.0.0.1:9000", "address to serve the api")
	reflectionFlag *bool   = flag.Bool("enable-grpc-reflection", true, "enables grpc reflection")
)

func main() {
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)
	flag.Parse()
	store, err := memdb.New()
	if err != nil {
		slog.Error("unable to create store", "error", err)
		os.Exit(1)
	}
	svc, err := services.NewStatusThingService(store, services.WithDefaults())
	if err != nil {
		slog.Error("unable to create service", "error", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	if *reflectionFlag {
		reflector := grpcreflect.NewStaticReflector(
			"statusthing.v1.ItemsService",
			"statusthing.v1.StatusService",
			"statusthing.v1.NotesService",
		)
		mux.Handle(grpcreflect.NewHandlerV1(reflector))
		mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
	}

	apiHandler, err := handlers.NewAPIHandler(svc)
	if err != nil {
		slog.Error("unable to create api handler", "error", err)
		os.Exit(1)
	}
	itempath, itemhandler := statusthingv1connect.NewItemsServiceHandler(apiHandler)
	mux.Handle(itempath, itemhandler)
	notepath, notehandler := statusthingv1connect.NewNotesServiceHandler(apiHandler)
	mux.Handle(notepath, notehandler)
	statuspath, statushandler := statusthingv1connect.NewStatusServiceHandler(apiHandler)
	mux.Handle(statuspath, statushandler)
	server := &http.Server{
		Addr:     *apiAddr,
		Handler:  h2c.NewHandler(requestLogger(mux), &http2.Server{}),
		ErrorLog: slog.NewLogLogger(logHandler, slog.LevelError),
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		if err := server.Shutdown(context.TODO()); err != nil {
			slog.Error("error shutting down server", "error", err)
		}
	}()

	slog.Info("starting server")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("error starting server", "error", err)
		os.Exit(1)
	}
	slog.Info("server stopped")
}

func requestLogger(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("handling request", "http.path", r.URL.Path, "http.host", r.Host, "http.client", r.Header.Get("User-Agent"))
		next.ServeHTTP(w, r)
	}
}
