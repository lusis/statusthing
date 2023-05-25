package sqlite

import (
	"context"
	"database/sql"
	"io/ioutil"
	"testing"

	"github.com/lusis/statusthing/internal/storers"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite" // sql driver
)

func makeTempFilename(t testing.TB) string {
	f, err := ioutil.TempFile("", "statusthing-storer-sqlite-tests-")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func makeTestdb(t *testing.T, option string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", option)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if err := createTables(context.TODO(), db); err != nil {
		return nil, err
	}

	return db, nil
}

func TestImplementsStatusStorer(t *testing.T) {
	t.Parallel()
	require.Implements(t, (*storers.StatusStorer)(nil), &Store{})
}
func TestCreateTables(t *testing.T) {
	t.Parallel()
	_, err := makeTestdb(t, ":memory:")
	require.NoError(t, err)
}

func TestCreateTablesIdempotent(t *testing.T) {
	t.Parallel()
	db, err := makeTestdb(t, ":memory:")
	require.NoError(t, err)
	cerr := createTables(context.TODO(), db)
	require.NoError(t, cerr)
}
