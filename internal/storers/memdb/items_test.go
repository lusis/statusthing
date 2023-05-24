package memdb

import (
	"context"
	"testing"
	"time"

	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/errors"
	"github.com/lusis/statusthing/internal/filters"
	"github.com/lusis/statusthing/internal/testutils"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestProtoFromMemItem(t *testing.T) {
	testcases := map[string]struct {
		item           *dbItem
		expected       *statusthingv1.Item
		err            error
		validationFunc func(*testing.T, *dbItem, *statusthingv1.Item)
	}{
		"missing-id": {
			item: &dbItem{},
			err:  errors.ErrInvalidData,
		},
		"missing-name": {
			item: &dbItem{ID: "missing-name"},
			err:  errors.ErrInvalidData,
		},
		"missing-updated": {
			item: &dbItem{ID: "missing-updated", Name: "missing-updated", Created: int(time.Now().UTC().UnixNano())},
			err:  errors.ErrInvalidData,
		},
		"missing-created": {
			item: &dbItem{ID: "missing-updated", Name: "missing-updated", Updated: int(time.Now().UTC().UnixNano())},
			err:  errors.ErrInvalidData,
		},
		"with-deleted": {
			item: &dbItem{ID: "missing-updated", Name: "missing-updated", Created: int(time.Now().UTC().UnixNano()), Updated: int(time.Now().UTC().UnixNano()), Deleted: int(time.Now().UTC().UnixNano())},
			validationFunc: func(t *testing.T, di *dbItem, i *statusthingv1.Item) {
				require.NotEqual(t, 0, di.Deleted, "should have a deleted timestamp")
			},
		},
		"with-description": {
			item: &dbItem{ID: "missing-updated", Name: "missing-updated", Created: int(time.Now().UTC().UnixNano()), Updated: int(time.Now().UTC().UnixNano()), Description: "mydescription"},
			validationFunc: func(t *testing.T, di *dbItem, i *statusthingv1.Item) {
				require.Equal(t, "mydescription", i.GetDescription(), "should have a description")
			},
		},
		"with-status": {
			item: &dbItem{ID: "missing-updated", Name: "missing-updated", Created: int(time.Now().UTC().UnixNano()), Updated: int(time.Now().UTC().UnixNano()), StatusID: "statusid"},
			validationFunc: func(t *testing.T, di *dbItem, i *statusthingv1.Item) {
				require.Nil(t, i.GetStatus())
			},
		},
	}
	t.Parallel()
	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			res, err := tc.item.toProto()
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				require.Nil(t, res)
			} else {
				require.NoError(t, err)
				require.NotNil(t, res)
			}
			if tc.validationFunc != nil {
				tc.validationFunc(t, tc.item, res)
			}
		})
	}
}

func TestDbItemFromProto(t *testing.T) {
	testcases := map[string]struct {
		item           *statusthingv1.Item
		expected       *dbItem
		err            error
		validationFunc func(*testing.T, *dbItem, *statusthingv1.Item)
	}{
		"nil-item": {
			item: nil,
			err:  errors.ErrNilVal,
		},
		"missing-id": {
			item: &statusthingv1.Item{},
			err:  errors.ErrEmptyString,
		},
		"missing-name": {
			item: &statusthingv1.Item{Id: "missing-name"},
			err:  errors.ErrEmptyString,
		},
		"missing-timestamps": {
			item: &statusthingv1.Item{Id: "missing-timestamps", Name: "missing-timestamps"},
			err:  errors.ErrNilVal,
		},
		"missing-updated": {
			item: &statusthingv1.Item{Id: "missing-updated", Name: "missing-updated", Timestamps: &statusthingv1.Timestamps{
				Created: timestamppb.Now(),
			}},
			err: errors.ErrInvalidData,
		},
		"missing-created": {
			item: &statusthingv1.Item{Id: "missing-created", Name: "missing-created", Timestamps: &statusthingv1.Timestamps{
				Updated: timestamppb.Now(),
			}},
			err: errors.ErrInvalidData,
		},
		"with-deleted": {
			item: &statusthingv1.Item{Id: "with-deleted", Name: "with-deleted", Timestamps: &statusthingv1.Timestamps{
				Created: timestamppb.Now(),
				Updated: timestamppb.Now(),
				Deleted: timestamppb.Now(),
			}},
			validationFunc: func(t *testing.T, di *dbItem, i *statusthingv1.Item) {
				require.NotEqual(t, 0, di.Deleted, "should have a deleted timestamp")
			},
		},
		"with-description": {
			item: &statusthingv1.Item{Id: "with-description", Name: "with-description", Timestamps: &statusthingv1.Timestamps{
				Created: timestamppb.Now(),
				Updated: timestamppb.Now(),
			}, Description: "mydesc"},
			validationFunc: func(t *testing.T, di *dbItem, i *statusthingv1.Item) {
				require.Equal(t, "mydesc", di.Description, "should have a description")
			},
		},
		"with-status": {
			item: &statusthingv1.Item{Id: "with-status", Name: "with-status", Timestamps: &statusthingv1.Timestamps{
				Created: timestamppb.Now(),
				Updated: timestamppb.Now(),
			}, Status: &statusthingv1.Status{Id: "myid"}},
			validationFunc: func(t *testing.T, di *dbItem, i *statusthingv1.Item) {
				require.Equal(t, "myid", di.StatusID)
			},
		},
	}
	t.Parallel()
	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			res, err := memItemFromProto((tc.item))
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				require.Nil(t, res)
			} else {
				require.NoError(t, err)
				require.NotNil(t, res)
			}
			if tc.validationFunc != nil {
				tc.validationFunc(t, res, tc.item)
			}
		})
	}
}
func TestItemsLifecycleWithOptions(t *testing.T) {
	// TODO: missed supporting adding a status at item creation time
}

func TestAllItems(t *testing.T) {
	type testitem struct {
		item *statusthingv1.Item
		opts []filters.FilterOption
	}
	type testcase struct {
		opts        []filters.FilterOption
		items       []*testitem
		expectedLen int
		err         error
	}
	testcases := map[string]testcase{
		"no-options": {
			items: []*testitem{
				{item: testutils.MakeItem(t.Name())},
			},
			expectedLen: 1,
		},
		"statuskinds-not-matched": {
			items: []*testitem{
				{
					item: &statusthingv1.Item{
						Id:         t.Name(),
						Name:       t.Name(),
						Timestamps: makeTestTsNow(),
						Status: &statusthingv1.Status{
							Id:         t.Name(),
							Name:       t.Name(),
							Kind:       statusthingv1.StatusKind_STATUS_KIND_DOWN,
							Timestamps: makeTestTsNow(),
						},
					},
				}},
			expectedLen: 0,
			opts:        []filters.FilterOption{filters.WithStatusKinds(statusthingv1.StatusKind_STATUS_KIND_AVAILABLE)},
		},
		"statuskinds-matched": {
			items: []*testitem{
				{
					item: &statusthingv1.Item{
						Id:         t.Name(),
						Name:       t.Name(),
						Timestamps: makeTestTsNow(),
						Status: &statusthingv1.Status{
							Id:         t.Name(),
							Name:       t.Name(),
							Kind:       statusthingv1.StatusKind_STATUS_KIND_DOWN,
							Timestamps: makeTestTsNow(),
						},
					},
				}},
			expectedLen: 1,
			opts:        []filters.FilterOption{filters.WithStatusKinds(statusthingv1.StatusKind_STATUS_KIND_DOWN)},
		},
		"statusids-matched": {
			items: []*testitem{
				{
					item: &statusthingv1.Item{
						Id:         t.Name(),
						Name:       t.Name(),
						Timestamps: makeTestTsNow(),
						Status: &statusthingv1.Status{
							Id:         "my-status-id",
							Name:       t.Name(),
							Kind:       statusthingv1.StatusKind_STATUS_KIND_DOWN,
							Timestamps: makeTestTsNow(),
						},
					},
				}},
			expectedLen: 1,
			opts:        []filters.FilterOption{filters.WithStatusIDs("my-status-id")},
		},
		"statusids-not-matched": {
			items: []*testitem{
				{
					item: &statusthingv1.Item{
						Id:         t.Name(),
						Name:       t.Name(),
						Timestamps: makeTestTsNow(),
						Status: &statusthingv1.Status{
							Id:         "my-status-id",
							Name:       t.Name(),
							Kind:       statusthingv1.StatusKind_STATUS_KIND_DOWN,
							Timestamps: makeTestTsNow(),
						},
					},
				}},
			expectedLen: 0,
			opts:        []filters.FilterOption{filters.WithStatusIDs("invalid-status-id")},
		},
	}

	for n, tc := range testcases {
		ctx := context.TODO()
		t.Run(n, func(t *testing.T) {
			store, _ := New()
			require.NotNil(t, store)

			for _, is := range tc.items {
				ires, ierr := store.StoreItem(ctx, is.item)
				require.NoError(t, ierr)
				require.NotNil(t, ires)
				if tc.err != nil {
					require.ErrorIs(t, ierr, tc.err)
					require.Nil(t, ires)
				} else {
					require.NoError(t, ierr)
					require.NotNil(t, ires)

					allres, allerr := store.FindItems(ctx, tc.opts...)
					require.NoError(t, allerr)
					require.Len(t, allres, tc.expectedLen)
				}
			}
		})
	}
}
func TestItemsLifecycle(t *testing.T) {
	ctx := context.TODO()
	store, _ := New()
	item := testutils.MakeItem(t.Name())
	status := testutils.MakeStatus(t.Name())

	res, err := store.StoreItem(context.TODO(), item)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.True(t, proto.Equal(res, item), "item and stored item should be the same")

	sres, serr := store.StoreStatus(ctx, status)
	require.NoError(t, serr, "status should be stored")
	require.NotNil(t, sres, "returned status should not be nil")

	fres, err := store.FindItems(ctx)
	require.NoError(t, err)
	require.Len(t, fres, 1, "should find one item")
	require.True(t, proto.Equal(res, fres[0]), "item and stored item should be the same")

	uerr := store.UpdateItem(
		ctx,
		res.GetId(),
		filters.WithName("new_name"),
		filters.WithDescription("new_description"),
		filters.WithStatus(status),
	)
	require.NoError(t, uerr)

	checkRes, checkErr := store.GetItem(ctx, res.GetId())
	require.NoError(t, checkErr)
	require.Greater(t, checkRes.GetTimestamps().GetUpdated().AsTime(), res.GetTimestamps().GetUpdated().AsTime(), "updated should be later than create")
	require.Equal(t, res.GetId(), checkRes.GetId(), "id should be same")
	require.Equal(t, "new_name", checkRes.GetName(), "name should be updated")
	require.Equal(t, "new_description", checkRes.GetDescription(), "description should be updated")
	require.NotNil(t, checkRes.GetStatus())
	require.True(t, proto.Equal(sres, checkRes.GetStatus()), "status should match previously inserted status")

	delErr := store.DeleteItem(ctx, res.GetId())
	require.NoError(t, delErr, "item should delete without error")

	checkDelRes, checkDelErr := store.GetItem(ctx, res.GetId())
	require.ErrorIs(t, checkDelErr, errors.ErrNotFound)
	require.Nil(t, checkDelRes)
}

func TestGetItem(t *testing.T) {
	ctx := context.TODO()
	t.Run("happy-path", func(t *testing.T) {})
	t.Run("invalid-status-id", func(t *testing.T) {
		store, _ := New()
		// we have to manually add an entry to the db so we can set an invalid statusid
		testitem := testutils.MakeItem(t.Name())
		txn := store.db.Txn(true)
		mItem, merr := memItemFromProto(testitem)
		require.NoError(t, merr)
		require.NotNil(t, mItem)

		mItem.StatusID = t.Name()
		ierr := txn.Insert(itemTableName, mItem)
		require.NoError(t, ierr)
		txn.Commit()

		ures, uerr := store.GetItem(ctx, mItem.ID)
		require.ErrorIs(t, uerr, errors.ErrInvalidData)
		require.Nil(t, ures)
	})
}

func TestUpdateItem(t *testing.T) {
	ctx := context.TODO()
	t.Run("no-opts", func(t *testing.T) {
		store, _ := New()
		rerr := store.UpdateItem(ctx, t.Name())
		require.ErrorIs(t, rerr, errors.ErrAtLeastOne)
	})
	t.Run("empty-id", func(t *testing.T) {
		store, _ := New()
		rerr := store.UpdateItem(ctx, "", filters.WithColor(t.Name()))
		require.ErrorIs(t, rerr, errors.ErrEmptyString)
	})
}
