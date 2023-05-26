package memdb

import (
	"context"
	"database/sql"

	"github.com/lusis/statusthing/internal/serrors"
	"github.com/lusis/statusthing/internal/storers/sqlite"
)

// Store is an in-memory store
type Store = sqlite.Store

// New returns a new [Store]
func New() (*Store, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	if err := sqlite.CreateTables(context.TODO(), db); err != nil {
		return nil, serrors.NewWrappedError("create-tables", serrors.ErrStoreUnavailable, err)
	}

	return sqlite.New(db)
}
