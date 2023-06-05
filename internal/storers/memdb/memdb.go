package memdb

import (
	"context"
	"os"

	"github.com/lusis/statusthing/internal/storers/sqlite"
	"github.com/lusis/statusthing/migrations"

	"golang.org/x/exp/slog"
)

// Store is an in-memory store
type Store = sqlite.Store

// New returns a new [Store]
func New() (*Store, error) {
	db, err := migrations.MigrateDatabase(context.TODO(), "sqlite3", ":memory:", slog.NewTextHandler(os.Stdout, nil))
	if err != nil {
		return nil, err
	}

	return sqlite.New(db)
}
