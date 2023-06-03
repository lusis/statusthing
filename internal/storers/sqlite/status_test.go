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
)

func TestStatusLifecycle(t *testing.T) {
	ctx := context.TODO()
	db, dberr := makeTestdb(t, ":memory:")
	require.NoError(t, dberr)
	require.NotNil(t, db)
	store, _ := New(db)

	// Store
	res, reserr := store.StoreStatus(ctx, &statusthingv1.Status{
		Id:         t.Name() + "id",
		Name:       t.Name() + "name",
		Kind:       statusthingv1.StatusKind_STATUS_KIND_CREATED,
		Timestamps: testutils.MakeTimestamps(false),
	})
	require.NoError(t, reserr)
	require.NotNil(t, res)
	require.Equal(t, t.Name()+"id", res.GetId())
	require.Equal(t, t.Name()+"name", res.GetName())
	require.Equal(t, statusthingv1.StatusKind_STATUS_KIND_CREATED, res.GetKind())
	require.NotNil(t, res.GetTimestamps)
	require.NotNil(t, res.GetTimestamps().GetCreated())
	require.NotNil(t, res.GetTimestamps().GetUpdated())

	// Update
	uerr := store.UpdateStatus(ctx, res.GetId(),
		filters.WithColor("new-color"),
		filters.WithDescription("new-description"),
		filters.WithName("new-name"),
		filters.WithStatusKind(statusthingv1.StatusKind_STATUS_KIND_DECOMM),
	)
	require.NoError(t, uerr)

	// Get
	gres, gerr := store.GetStatus(ctx, res.GetId())
	require.NoError(t, gerr)
	require.NotNil(t, gres)
	require.Equal(t, res.GetId(), gres.GetId())
	require.Equal(t, "new-name", gres.GetName())
	require.Equal(t, "new-description", gres.GetDescription())
	require.Equal(t, "new-color", gres.GetColor())
	require.Equal(t, statusthingv1.StatusKind_STATUS_KIND_DECOMM, gres.GetKind())

	// Delete
	delerr := store.DeleteStatus(ctx, res.GetId())
	require.NoError(t, delerr)
	delagainerr := store.DeleteStatus(ctx, res.GetId())
	require.ErrorIs(t, delagainerr, serrors.ErrNotFound)
}

func TestFindStatus(t *testing.T) {
	ctx := context.TODO()
	type testCase struct {
		opts  []filters.FilterOption
		kinds []statusthingv1.StatusKind
		ids   []string
		err   error
		count int
	}
	testCases := map[string]testCase{
		"all-returns-empty": {
			count: 0,
		},
		"finds-ids": {
			count: 1,
			ids:   []string{"id-one"},
			opts:  []filters.FilterOption{filters.WithStatusIDs("id-one_status_id")}, // status_id is appended by helper
		},
		"not-finds-ids": {
			count: 0,
			ids:   []string{"id-one"},
			opts:  []filters.FilterOption{filters.WithStatusIDs("invalid")}, // status_id is appended by helper
		},
		"finds-kind": {
			count: 1,
			kinds: []statusthingv1.StatusKind{
				statusthingv1.StatusKind_STATUS_KIND_CREATED,
			},
			opts: []filters.FilterOption{
				filters.WithStatusKinds(statusthingv1.StatusKind_STATUS_KIND_CREATED),
			},
		},
		"ids-and-kinds": {
			count: 2,
			kinds: []statusthingv1.StatusKind{
				statusthingv1.StatusKind_STATUS_KIND_CREATED,
				statusthingv1.StatusKind_STATUS_KIND_DOWN,
				statusthingv1.StatusKind_STATUS_KIND_DECOMM,
			},
			ids: []string{"test-1", "test-2", "test-3"},
			opts: []filters.FilterOption{
				filters.WithStatusKinds(statusthingv1.StatusKind_STATUS_KIND_DECOMM),
				filters.WithStatusIDs("test-1_status_id"),
			},
		},
		"not-find-kind": {
			count: 0,
			kinds: []statusthingv1.StatusKind{
				statusthingv1.StatusKind_STATUS_KIND_CREATED,
			},
			opts: []filters.FilterOption{
				filters.WithStatusKinds(statusthingv1.StatusKind_STATUS_KIND_DECOMM),
			},
		},
		"multiple-same-kind": {
			count: 2,
			kinds: []statusthingv1.StatusKind{
				statusthingv1.StatusKind_STATUS_KIND_CREATED,
				statusthingv1.StatusKind_STATUS_KIND_CREATED,
				statusthingv1.StatusKind_STATUS_KIND_DOWN,
			},
			opts: []filters.FilterOption{
				filters.WithStatusKinds(statusthingv1.StatusKind_STATUS_KIND_CREATED),
			},
		},
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			db, dberr := makeTestdb(t, ":memory:")
			defer db.Close()
			require.NoError(t, dberr)
			require.NotNil(t, db)
			store, err := New(db)
			require.NoError(t, err)

			for i, kind := range tc.kinds {
				testStatus := testutils.MakeStatus(fmt.Sprintf("status-%d", i))
				testStatus.Kind = kind

				_, serr := store.StoreStatus(ctx, testStatus)
				require.NoError(t, serr)
			}
			for _, id := range tc.ids {
				testStatus := testutils.MakeStatus(id)
				_, serr := store.StoreStatus(ctx, testStatus)
				require.NoError(t, serr)
			}
			fres, ferr := store.FindStatus(ctx, tc.opts...)
			if tc.err != nil {
				require.ErrorIs(t, ferr, tc.err)
				require.Nil(t, fres)
			} else {
				require.NoError(t, ferr)
				require.Len(t, fres, tc.count)
			}
		})

	}
}
