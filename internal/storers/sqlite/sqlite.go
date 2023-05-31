package sqlite

import (
	"context"
	"database/sql"
	"database/sql/driver"

	"github.com/lusis/statusthing/internal/serrors"
	"github.com/lusis/statusthing/internal/storers/unimplemented"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3" // goqu dialect
	"modernc.org/sqlite"                               // sql driver
)

// Store stores statusthing data
type Store struct {
	*unimplemented.StatusThingStore
	db     *sql.DB
	goqudb *goqu.Database
}

// we need to register our own wrapper over modernc to enable fks reliably
// this pattern is documented here: https://github.com/ent/ent/discussions/1667#discussioncomment-1132296
// however we may want to wrap other drivers or add additional pragmas here
func init() {
	sql.Register("sqlite3", fkDriver{Driver: &sqlite.Driver{}})
}

type fkDriver struct {
	*sqlite.Driver
}

// Open opens a connection with the PRAGMA set
func (d fkDriver) Open(name string) (driver.Conn, error) {
	conn, err := d.Driver.Open(name)
	if err != nil {
		return conn, err
	}
	c := conn.(interface {
		Exec(stmt string, args []driver.Value) (driver.Result, error)
	})
	if _, err := c.Exec("PRAGMA foreign_keys = on;", nil); err != nil {
		conn.Close()
		return nil, serrors.NewWrappedError("driver", serrors.ErrStoreUnavailable, err)
	}
	return conn, nil
}

// New returns a new [Store]
func New(db *sql.DB) (*Store, error) {
	goqu.SetIgnoreUntaggedFields(true)
	gdb := goqu.New("sqlite3", db)
	return &Store{db: db, goqudb: gdb}, nil
}

// CreateTables creates the required tables for the sqlite3 store
func CreateTables(ctx context.Context, db *sql.DB) error {
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
	// if _, err := txn.ExecContext(ctx, stmtCreateUsersTable); err != nil {
	// 	if rerr := txn.Rollback(); rerr != nil {
	// 		return serrors.NewWrappedError("rollback", serrors.ErrUnrecoverable, rerr)
	// 	}
	// 	return serrors.NewWrappedError("create-table", serrors.ErrUnrecoverable, err)
	// }
	if err := txn.Commit(); err != nil {
		return serrors.NewWrappedError("commit", serrors.ErrUnrecoverable, err)
	}
	return nil
}
