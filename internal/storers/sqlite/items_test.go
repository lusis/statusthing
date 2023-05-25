package sqlite

import (
	"context"
	"testing"

	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/filters"
	"github.com/lusis/statusthing/internal/serrors"
	"github.com/lusis/statusthing/internal/testutils"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

type itemTestcase struct {
	dbthing *dbItem
	pbthing *statusthingv1.Item
	errtext string
	err     error
}

func TestItemLifecycle(t *testing.T) {
	ctx := context.TODO()
	db, dberr := makeTestdb(t, ":memory:")
	require.NoError(t, dberr)
	require.NotNil(t, db)
	store, _ := New(db)

	// Store
	res, reserr := store.StoreItem(ctx, &statusthingv1.Item{
		Id:         t.Name() + "id",
		Name:       t.Name() + "name",
		Timestamps: testutils.MakeTimestamps(false),
	})
	require.NoError(t, reserr)
	require.NotNil(t, res)
	require.Equal(t, t.Name()+"id", res.GetId())
	require.Equal(t, t.Name()+"name", res.GetName())

	require.NotNil(t, res.GetTimestamps)
	require.NotNil(t, res.GetTimestamps().GetCreated())
	require.NotNil(t, res.GetTimestamps().GetUpdated())

	// Get
	gres, gerr := store.GetItem(ctx, res.GetId())
	require.NoError(t, gerr)
	require.NotNil(t, gres)
	require.Equal(t, res.GetId(), gres.GetId())
	require.Equal(t, res.GetName(), gres.GetName())

	// to update we need a status
	sres, serr := store.StoreStatus(ctx, testutils.MakeStatus(t.Name()))
	require.NoError(t, serr)
	require.NotNil(t, sres)
	// Update
	uerr := store.UpdateItem(ctx, res.GetId(),
		filters.WithDescription("new-description"),
		filters.WithName("new-name"),
		filters.WithStatusID(sres.GetId()),
	)
	require.NoError(t, uerr)

	// Get Again
	checkres, checkerr := store.GetItem(ctx, res.GetId())
	require.NoError(t, checkerr)
	require.NotNil(t, checkres)
	require.Equal(t, res.GetId(), checkres.GetId())
	require.Equal(t, "new-name", checkres.GetName())
	require.Equal(t, "new-description", checkres.GetDescription())
	require.NotNil(t, checkres.GetStatus())
	require.True(t, proto.Equal(sres, checkres.GetStatus()))

	// Delete
	delerr := store.DeleteItem(ctx, res.GetId())
	require.NoError(t, delerr)
	delagainerr := store.DeleteItem(ctx, res.GetId())
	require.ErrorIs(t, delagainerr, serrors.ErrNotFound)
}
