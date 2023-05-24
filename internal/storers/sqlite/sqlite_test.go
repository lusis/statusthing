package sqlite

import (
	"context"
	"database/sql"
	"io/ioutil"
	"os"
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

func makeTestdb(t *testing.T, option string) (*sql.DB, func(), error) {
	tempFilename := makeTempFilename(t)
	url := tempFilename + option

	cleanupFunc := func() {
		err := os.Remove(tempFilename)
		if err != nil {
			t.Error("temp file remove error:", err)
		}
	}

	db, err := sql.Open("sqlite", url)
	if err != nil {
		return nil, nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, cleanupFunc, err
	}
	if err := createTables(context.TODO(), db); err != nil {
		return nil, cleanupFunc, err
	}

	return db, cleanupFunc, nil
}

func TestImplementsStatusStorer(t *testing.T) {
	t.Parallel()
	require.Implements(t, (*storers.StatusStorer)(nil), &Store{})
}
func TestCreateTables(t *testing.T) {
	t.Parallel()
	_, cleanup, err := makeTestdb(t, "")
	defer cleanup()
	require.NoError(t, err)
}

func TestCreateTablesIdempotent(t *testing.T) {
	t.Parallel()
	db, cleanup, err := makeTestdb(t, "")
	defer cleanup()
	require.NoError(t, err)
	cerr := createTables(context.TODO(), db)
	require.NoError(t, cerr)
}
