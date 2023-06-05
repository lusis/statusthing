package sqlite

import (
	"context"
	"database/sql"
	"io/ioutil"
	"os"
	"testing"

	"github.com/lusis/statusthing/internal/storers"
	_ "github.com/lusis/statusthing/internal/storers/sqlite/driver" // sql driver
	"github.com/lusis/statusthing/migrations"

	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slog"
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
	db, err := migrations.MigrateDatabase(context.TODO(), "sqlite3", option, slog.NewTextHandler(os.Stdout, nil))
	if err != nil {
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
	t.Skip("fix")
	t.Parallel()
	db, err := makeTestdb(t, ":memory:")
	require.NoError(t, err)
	require.NotNil(t, db)

	db, err = migrations.MigrateDatabase(context.TODO(), "sqlite3", ":memory:", slog.NewTextHandler(os.Stdout, nil))

	require.NoError(t, err)
}
