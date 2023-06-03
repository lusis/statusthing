package driver

import (
	"database/sql"
	"database/sql/driver"

	"github.com/lusis/statusthing/internal/serrors"

	"modernc.org/sqlite"
)

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
