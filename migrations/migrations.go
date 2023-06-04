// Package migrations stores migrations for databases
package migrations

import (
	"context"
	"database/sql"
	"embed"
	"io/fs"

	"github.com/lusis/statusthing/internal/serrors"
	"github.com/lusis/statusthing/internal/validation"
	"github.com/lusis/statusthing/migrations/sqlite3"
	"golang.org/x/exp/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	_ "github.com/go-sql-driver/mysql"                      // import mysql driver
	_ "github.com/golang-migrate/migrate/v4/database/mysql" // mysql migrate driver
	_ "github.com/golang-migrate/migrate/v4/source/file"    // for loading migrations
)

// migrationFS is the filesystem storing migrations
//
//go:embed all:data
var migrationFS embed.FS

// MigrateDatabase migrates the database of type dbType using the provided db
// Caller is responsible for closing the db connection after
// TODO: Don't know how I feel about returning the db here but it works for now
func MigrateDatabase(_ context.Context, driver string, dsn string, handler slog.Handler) (*sql.DB, error) {
	if !validation.ValidString(driver) {
		return nil, serrors.NewError("driver", serrors.ErrEmptyString)
	}
	if !validation.ValidString(dsn) {
		return nil, serrors.NewError("dsn", serrors.ErrEmptyString)
	}
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, serrors.NewWrappedError("driver", serrors.ErrStoreUnavailable, err)
	}
	if err := db.Ping(); err != nil {
		return nil, serrors.NewWrappedError("db", serrors.ErrStoreUnavailable, err)
	}
	var databaseDriver database.Driver
	var migrationDir string
	var subfs fs.FS
	migrationDir = "."
	switch driver {
	case "sqlite", "sqlite3":
		dst, err := sqlite3.WithInstance(db, &sqlite3.Config{})
		if err != nil {
			return nil, serrors.NewWrappedError("destination", serrors.ErrUnrecoverable, err)
		}
		databaseDriver = dst
		myfs, err := fs.Sub(migrationFS, "data/sqlite")
		if err != nil {
			return nil, serrors.NewWrappedError("migrations-subfs", serrors.ErrUnrecoverable, err)
		}
		subfs = myfs
	default:
		return nil, serrors.NewError("driver", serrors.ErrNotImplemented)
	}

	if !validation.ValidString(migrationDir) {
		return nil, serrors.NewError("driver", serrors.ErrNotImplemented)
	}
	if subfs == nil {
		return nil, serrors.NewError("driver", serrors.ErrNotImplemented)
	}
	migrationfs, err := iofs.New(subfs, migrationDir)
	if err != nil {
		return nil, serrors.NewWrappedError("migrationfs-iofs", serrors.ErrUnrecoverable, err)
	}
	defer migrationfs.Close()
	migrator, err := migrate.NewWithInstance("iofs", migrationfs, driver, databaseDriver)
	if err != nil {
		return nil, serrors.NewWrappedError("migrations", serrors.ErrUnrecoverable, err)
	}
	migrator.Log = &migrationSlogger{handler: handler}
	if merr := migrator.Up(); merr != nil {
		if merr != migrate.ErrNoChange {
			return nil, serrors.NewWrappedError("migrations", serrors.ErrUnrecoverable, merr)
		}
	}
	return db, nil
}

type migrationSlogger struct {
	handler slog.Handler
}

// Printf satisfies an interface
func (l *migrationSlogger) Printf(format string, v ...any) {
	ll := slog.NewLogLogger(l.handler, slog.LevelInfo)
	ll.Printf(format, v...)
}

// Verbose satisfies an interface
func (l *migrationSlogger) Verbose() bool {
	return true
}
