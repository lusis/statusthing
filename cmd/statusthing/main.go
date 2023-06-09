// Package main ...
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	flag "github.com/spf13/pflag"

	"github.com/lusis/statusthing/internal"
	"github.com/lusis/statusthing/internal/storers/sqlite"
	"github.com/lusis/statusthing/migrations"

	"golang.org/x/exp/slog"
)

var (
	apiAddr *string = flag.String("api-addr", "127.0.0.1:9000", "address to serve the api")
	devMode *bool   = flag.Bool("devmode", false, "enables grpc reflection and template reloading for development")
)

func main() {
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)
	flag.Parse()
	db, err := migrations.MigrateDatabase(context.TODO(), "sqlite3", "statusthing.db", logHandler)
	if err != nil {
		logger.Error("error migrating database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	store, err := sqlite.New(db)
	if err != nil {
		logger.Error("unable to create store", "error", err)
		os.Exit(1)
	}
	if err := db.Ping(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	server, err := internal.New(store, *apiAddr, logHandler, *devMode)
	if err != nil {
		slog.Error("cannot create statusthing", "error", err)
		os.Exit(1)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		if err := server.Stop(context.TODO()); err != nil {
			slog.Error("error shutting down", "error", err)
		}
	}()

	slog.Info("starting statusthing")
	if err := server.Start(); err != nil {
		slog.Error("error starting", "error", err)
		os.Exit(1)
	}
	slog.Info("statusthing stopped")
}
