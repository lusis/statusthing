package sqlite

import (
	"context"
	"database/sql"

	"github.com/lusis/statusthing/internal/serrors"
	"github.com/lusis/statusthing/internal/storers/unimplemented"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3" // goqu dialect
	_ "modernc.org/sqlite"                             // sql driver
)

// Store stores statusthing data
type Store struct {
	*unimplemented.StatusThingStore
	db     *sql.DB
	goqudb *goqu.Database
}

// New returns a new [Store]
func New(db *sql.DB) (*Store, error) {
	goqu.SetIgnoreUntaggedFields(true)
	gdb := goqu.New("sqlite3", db)
	return &Store{db: db, goqudb: gdb}, nil
}

func createTables(ctx context.Context, db *sql.DB) error {
	txn, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return serrors.NewWrappedError("start-transaction", serrors.ErrUnrecoverable, err)
	}
	if txn == nil {
		return serrors.NewError("txn", serrors.ErrUnrecoverable)
	}

	// status -> items -> notes
	if _, err := txn.ExecContext(ctx, stmtCreateStatusTable); err != nil {
		if rerr := txn.Rollback(); rerr != nil {
			return serrors.NewWrappedError("rollback", serrors.ErrUnrecoverable, rerr)
		}
		return serrors.NewWrappedError("create-table", serrors.ErrUnrecoverable, err)
	}
	if _, err := txn.ExecContext(ctx, stmtCreateItemsTable); err != nil {
		if rerr := txn.Rollback(); rerr != nil {
			return serrors.NewWrappedError("rollback", serrors.ErrUnrecoverable, rerr)
		}
		return serrors.NewWrappedError("create-table", serrors.ErrUnrecoverable, err)
	}
	if _, err := txn.ExecContext(ctx, stmtCreateNotesTable); err != nil {
		if rerr := txn.Rollback(); rerr != nil {
			return serrors.NewWrappedError("rollback", serrors.ErrUnrecoverable, rerr)
		}
		return serrors.NewWrappedError("create-table", serrors.ErrUnrecoverable, err)
	}
	if err := txn.Commit(); err != nil {
		return serrors.NewWrappedError("commit", serrors.ErrUnrecoverable, err)
	}
	return nil
}
