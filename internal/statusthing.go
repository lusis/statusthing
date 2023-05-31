// Package internal contains internal code
package internal

import (
	"context"
	"net/http"

	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
	"github.com/go-chi/chi"
	v1connect "github.com/lusis/statusthing/gen/go/statusthing/v1/statusthingv1connect"
	"github.com/lusis/statusthing/internal/handlers"
	"github.com/lusis/statusthing/internal/serrors"
	"github.com/lusis/statusthing/internal/services"
	"github.com/lusis/statusthing/internal/storers"
	"github.com/lusis/statusthing/internal/validation"
	"golang.org/x/exp/slog"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// StatusThing is a statuspage application
type StatusThing struct {
	apiHandler   *handlers.APIHandler
	adminHandler *handlers.AdminHandler
	svc          *services.StatusThingService
	store        storers.StatusThingStorer
	mux          chi.Router
	httpServer   *http.Server
}

// New returns a new StatusThing
func New(store storers.StatusThingStorer, listenAddress string, logHandler slog.Handler, devMode bool) (*StatusThing, error) {
	if !validation.ValidString(listenAddress) {
		return nil, serrors.NewError("listenAddress", serrors.ErrEmptyString)
	}
	if store == nil {
		return nil, serrors.NewError("store", serrors.ErrNilVal)
	}
	svc, err := services.NewStatusThingService(store)
	if err != nil {
		return nil, serrors.NewWrappedError("service", serrors.ErrDependencyMissing, err)
	}
	mux := chi.NewRouter()
	// session.NewSession()
	// mux.Use(session.Sessions.LoadAndSave)
	if err := registerAPIHandler(mux, svc, devMode); err != nil {
		return nil, err
	}
	adminHandler, err := handlers.NewAdminHandler(svc, mux, devMode)
	if err != nil {
		return nil, serrors.NewWrappedError("adminhandler", serrors.ErrDependencyMissing, err)
	}
	server := &http.Server{
		Addr:     listenAddress,
		Handler:  h2c.NewHandler(requestLogger(mux), &http2.Server{}),
		ErrorLog: slog.NewLogLogger(logHandler, slog.LevelError),
	}
	st := &StatusThing{
		adminHandler: adminHandler,
		store:        store,
		mux:          mux,
		svc:          svc,
		httpServer:   server,
	}

	return st, nil
}

// Start starts the server
func (st *StatusThing) Start() error {
	if err := st.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Stop stops the server
func (st *StatusThing) Stop(ctx context.Context) error {
	return st.httpServer.Shutdown(ctx)
}

// Mux returns the configured mux for adding additiona routes
func (st *StatusThing) Mux() chi.Router {
	return st.mux
}

func registerAPIHandler(mux chi.Router, svc *services.StatusThingService, reflect bool) error {
	apiHandler, err := handlers.NewAPIHandler(svc)
	if err != nil {
		return serrors.NewWrappedError("apihandler", serrors.ErrDependencyMissing, err)
	}
	if reflect {
		reflector := grpcreflect.NewStaticReflector(
			"statusthing.v1.ItemsService",
			"statusthing.v1.StatusService",
			"statusthing.v1.NotesService",
		)
		mux.Mount(grpcreflect.NewHandlerV1(reflector))
		mux.Mount(grpcreflect.NewHandlerV1Alpha(reflector))
	}
	mux.Mount(v1connect.NewItemsServiceHandler(apiHandler))
	mux.Mount(v1connect.NewNotesServiceHandler(apiHandler))
	mux.Mount(v1connect.NewStatusServiceHandler(apiHandler))

	return nil
}

func requestLogger(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("handling request", "http.path", r.URL.Path, "http.host", r.Host, "http.method", r.Method, "http.client", r.Header.Get("User-Agent"), "content-type", r.Header.Get("content-type"))
		slog.Debug("all headers", "headers", r.Header)
		next.ServeHTTP(w, r)
	}
}
