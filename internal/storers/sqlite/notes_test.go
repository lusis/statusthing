package sqlite

import (
	"context"
	"testing"

	"github.com/lusis/statusthing/internal/filters"
	"github.com/lusis/statusthing/internal/serrors"
	_ "github.com/lusis/statusthing/internal/storers/sqlite/driver" // sql driver
	"github.com/lusis/statusthing/internal/testutils"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestNoteLifeCycle(t *testing.T) {
	ctx := context.TODO()
	db, dberr := makeTestdb(t, ":memory:")
	require.NoError(t, dberr)
	require.NotNil(t, db)
	store, _ := New(db)
	ires, ierr := store.StoreItem(ctx, testutils.MakeItem(t.Name()))
	require.NoError(t, ierr)
	require.NotNil(t, ires)

	res, reserr := store.StoreNote(ctx, testutils.MakeNote(t.Name()), ires.GetId())
	require.NoError(t, reserr)
	require.NotNil(t, res)

	gres, gerr := store.GetNote(ctx, res.GetId())
	require.NoError(t, gerr)
	require.NotNil(t, gres)
	require.Equal(t, res.GetId(), gres.GetId())
	require.Equal(t, res.GetText(), gres.GetText())
	require.NotNil(t, gres.GetTimestamps)
	require.NotNil(t, gres.GetTimestamps().GetCreated())
	require.NotNil(t, gres.GetTimestamps().GetUpdated())

	fres, ferr := store.FindNotes(ctx, ires.GetId())
	require.NoError(t, ferr)
	require.Len(t, fres, 1)
	require.Equal(t, res.GetId(), fres[0].GetId())
	require.Equal(t, res.GetText(), fres[0].GetText())
	require.NotNil(t, fres[0].GetTimestamps)
	require.NotNil(t, fres[0].GetTimestamps().GetCreated())
	require.NotNil(t, fres[0].GetTimestamps().GetUpdated())

	uerr := store.UpdateNote(ctx, res.GetId(), filters.WithNoteText("new-text"))
	require.NoError(t, uerr)

	cres, cerr := store.GetNote(ctx, res.GetId())
	require.NoError(t, cerr)
	require.NotNil(t, cres)
	require.Equal(t, res.GetId(), cres.GetId())
	require.Equal(t, "new-text", cres.GetText())
	require.NotNil(t, gres.GetTimestamps)
	require.NotNil(t, gres.GetTimestamps().GetCreated())
	require.NotNil(t, gres.GetTimestamps().GetUpdated())

	derr := store.DeleteNote(ctx, res.GetId())
	require.NoError(t, derr)

	dcerr := store.DeleteNote(ctx, res.GetId())
	require.ErrorIs(t, dcerr, serrors.ErrNotFound)

	// make sure we didn't somehow delete the item entry with an unintended cascade
	cires, cierr := store.GetItem(ctx, ires.GetId())
	require.NoError(t, cierr)
	require.True(t, proto.Equal(cires, ires))
}
