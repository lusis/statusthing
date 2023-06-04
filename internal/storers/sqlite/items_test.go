package sqlite

import (
	"context"
	"fmt"
	"testing"

	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/filters"
	"github.com/lusis/statusthing/internal/serrors"
	_ "github.com/lusis/statusthing/internal/storers/sqlite/driver" // sql driver
	"github.com/lusis/statusthing/internal/testutils"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

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

func TestFindItems(t *testing.T) {
	ctx := context.TODO()
	type testcase struct {
		items     []*statusthingv1.Item
		opts      []filters.FilterOption
		kinds     []statusthingv1.StatusKind
		statusids []string
		err       error
		count     int
	}
	testCases := map[string]testcase{
		"all-items": {
			items: []*statusthingv1.Item{
				testutils.MakeItem("item-1"),
				testutils.MakeItem("item-2"),
				func() *statusthingv1.Item {
					item := testutils.MakeItem("item-3")
					status := testutils.MakeStatus("status-3")
					item.Status = status
					return item
				}(),
			},
			count: 3,
		},
		"one-with-id": {
			items: []*statusthingv1.Item{
				testutils.MakeItem("item-1"),
				testutils.MakeItem("item-2"),
				func() *statusthingv1.Item {
					item := testutils.MakeItem("item-3")
					status := testutils.MakeStatus("status-3")
					// override for determinism
					status.Id = "status-3"
					item.Status = status
					return item
				}(),
			},
			opts:  []filters.FilterOption{filters.WithStatusIDs("status-3")},
			count: 1,
		},
		"one-with-kind": {
			items: []*statusthingv1.Item{
				testutils.MakeItem("one-with-kind-item-1"),
				testutils.MakeItem("one-with-kind-item-2"),
				func() *statusthingv1.Item {
					item := testutils.MakeItem("one-with-kind-item-3")
					status := testutils.MakeStatus("one-with-kind-status-3")
					// override for determinism
					status.Kind = statusthingv1.StatusKind_STATUS_KIND_DECOMM
					item.Status = status
					return item
				}(),
			},
			kinds: []statusthingv1.StatusKind{statusthingv1.StatusKind_STATUS_KIND_UNAVAILABLE, statusthingv1.StatusKind_STATUS_KIND_OFFLINE},
			opts:  []filters.FilterOption{filters.WithStatusKinds(statusthingv1.StatusKind_STATUS_KIND_DECOMM)},
			count: 1,
		},
		"one-with-kind-and-id": {
			items: []*statusthingv1.Item{
				testutils.MakeItem("one-with-kind-item-1"),
				testutils.MakeItem("one-with-kind-item-2"),
				func() *statusthingv1.Item {
					item := testutils.MakeItem("one-with-kind-item-3")
					status := testutils.MakeStatus("one-with-kind-status-3")
					// override for determinism
					status.Kind = statusthingv1.StatusKind_STATUS_KIND_DECOMM
					item.Status = status
					return item
				}(),
				func() *statusthingv1.Item {
					item := testutils.MakeItem("one-with-kind-item-4")
					status := testutils.MakeStatus("one-with-kind-status-4")
					// override for determinism
					status.Id = "wanted_id"
					status.Kind = statusthingv1.StatusKind_STATUS_KIND_WARNING
					item.Status = status
					return item
				}(),
			},
			kinds: []statusthingv1.StatusKind{statusthingv1.StatusKind_STATUS_KIND_UNAVAILABLE, statusthingv1.StatusKind_STATUS_KIND_OFFLINE},
			opts:  []filters.FilterOption{filters.WithStatusKinds(statusthingv1.StatusKind_STATUS_KIND_DECOMM), filters.WithStatusIDs("wanted_id")},
			count: 2,
		},
	}
	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			db, dberr := makeTestdb(t, ":memory:")
			require.NoError(t, dberr)
			require.NotNil(t, db)
			store, _ := New(db)

			for i, kind := range tc.kinds {
				testStatus := testutils.MakeStatus(fmt.Sprintf("adhoc-statuskind-%d", i))
				// overwrite the generated kind to align
				testStatus.Kind = kind

				_, serr := store.StoreStatus(ctx, testStatus)
				require.NoError(t, serr)
			}
			for _, id := range tc.statusids {
				testStatus := testutils.MakeStatus("adhoc-statusid-")
				// overwrite the generated id to align
				testStatus.Id = id
				_, serr := store.StoreStatus(ctx, testStatus)
				require.NoError(t, serr)
			}
			for _, item := range tc.items {
				ires, ierr := store.StoreItem(ctx, item)
				require.NoError(t, ierr)
				require.NotNil(t, ires)
			}
			// get all items for debugging
			// all, allerr := store.FindItems(ctx)
			// require.NoError(t, allerr)
			// for _, i := range all {
			// 	t.Logf("all: %+v\n", i)
			// }
			fres, ferr := store.FindItems(ctx, tc.opts...)
			require.NoError(t, ferr)
			require.NotNil(t, fres)

			require.Len(t, fres, tc.count)
		})
	}
}

func TestUpdateItem(t *testing.T) {
	ctx := context.TODO()
	t.Run("invalid-status-id", func(t *testing.T) {
		db, dberr := makeTestdb(t, ":memory:")
		require.NoError(t, dberr)
		require.NotNil(t, db)
		store, _ := New(db)
		item, itemerr := store.StoreItem(ctx, testutils.MakeItem(t.Name()))
		require.NoError(t, itemerr)
		require.NotNil(t, item)

		uerr := store.UpdateItem(ctx, item.GetId(), filters.WithStatusID("invalid"))
		require.Error(t, uerr, serrors.ErrNotFound)
		require.Contains(t, uerr.Error(), "item")

		citem, cerr := store.GetItem(ctx, item.GetId())
		require.NoError(t, cerr)
		require.NotNil(t, citem)
	})

}
