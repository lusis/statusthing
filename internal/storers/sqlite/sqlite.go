package sqlite

import (
	"database/sql"

	_ "github.com/lusis/statusthing/internal/storers/sqlite/driver" // sql driver
	"github.com/lusis/statusthing/internal/storers/unimplemented"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3" // goqu dialect
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
