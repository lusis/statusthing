package services

import (
	"context"
	"testing"

	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/errors"
	"github.com/lusis/statusthing/internal/filters"
	"github.com/lusis/statusthing/internal/storers/memdb"
	"github.com/lusis/statusthing/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestItems(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	testCases := map[string]struct {
		item           *statusthingv1.Item
		status         *statusthingv1.Status
		notelen        int
		itemStatusID   string
		opts           []filters.FilterOption
		err            error
		validationFunc func(*StatusThingService) bool
	}{
		"happy-path": {
			item: &statusthingv1.Item{Name: "my-service", Timestamps: makeTsNow()},
		},
		"missing-item-id": {
			err:  errors.ErrEmptyString,
			item: &statusthingv1.Item{Name: "my-service", Timestamps: makeTsNow()},
			opts: []filters.FilterOption{filters.WithItemID("")},
		},
		"with-alternate-status-missing-id": {
			item: &statusthingv1.Item{Name: "my-service", Timestamps: makeTsNow()},
			opts: []filters.FilterOption{
				filters.WithItemID("my-service"),
				filters.WithStatus(&statusthingv1.Status{Name: "my-status-name", Kind: statusthingv1.StatusKind_STATUS_KIND_AVAILABLE}),
			},
			validationFunc: func(sts *StatusThingService) bool {
				res, err := sts.GetItem(ctx, "my-service")
				if err != nil {
					t.Log(err)
					return false
				}
				if res == nil {
					return false
				}
				if res.GetStatus().GetName() != "my-status-name" {
					return false
				}
				return true
			},
		},
		"with-alternate-status": {
			itemStatusID: t.Name() + "_status_id",
			opts: []filters.FilterOption{
				filters.WithStatus(testutils.MakeStatus(t.Name())),
			},
			validationFunc: func(sts *StatusThingService) bool {
				res, err := sts.GetStatus(ctx, t.Name()+"_status_id")
				if err != nil {
					t.Log(err)
					return false
				}
				if res == nil {
					return false
				}
				return true
			},
		},
		"with-valid-status-id": {
			status: testutils.MakeStatus(t.Name()),
			opts:   []filters.FilterOption{filters.WithStatusID(t.Name() + "_status_id")},
		},
		"invalid-status-id": {
			opts: []filters.FilterOption{
				filters.WithStatusID("invalid"),
			},
			err: errors.ErrNotFound,
		},
		"with-initial-note": {
			// we need a custom item id to get the notes
			opts:    []filters.FilterOption{filters.WithNoteText("this is my text"), filters.WithItemID("my-item-id")},
			notelen: 1,
			validationFunc: func(sts *StatusThingService) bool {
				i, err := sts.AllNotes(ctx, "my-item-id")
				if err != nil {
					return false
				}
				if i[0].GetText() == "this is my text" {
					return true
				}
				return false
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// Set up our deps
			store, err := memdb.New()
			require.NoError(t, err)
			require.NotNil(t, store)
			sts, err := NewStatusThingService(store)
			require.NoError(t, err)
			require.NotNil(t, sts)

			item := tc.item
			if item == nil {
				item = testutils.MakeItem(t.Name())
			}

			if tc.status != nil {
				sres, serr := sts.store.StoreStatus(ctx, tc.status)
				require.NoError(t, serr)
				require.NotNil(t, sres)
			}

			ires, ierr := sts.NewItem(ctx, item.GetName(), tc.opts...)
			if tc.err != nil {
				require.ErrorIs(t, ierr, tc.err)
				require.Nil(t, ires)
			} else {
				require.NoError(t, ierr)
				require.NotNil(t, ires)
				require.Equal(t, item.GetName(), ires.GetName(), "name should match")
				require.Equal(t, item.GetDescription(), ires.GetDescription(), "description should match")
				require.Len(t, ires.GetNotes(), tc.notelen)
				if tc.itemStatusID != "" {
					require.NotNil(t, ires.GetStatus())
					require.Equal(t, tc.itemStatusID, ires.GetStatus().GetId())
				}
				// since we control timestamps, our own timestamps should never be used
				require.Falsef(t, testutils.TimestampsEqual(item.GetTimestamps(), ires.GetTimestamps()), "timestamps should not match. [expected] %+v | [actual] %+v\n", item.GetTimestamps(), ires.GetTimestamps())
			}
			if tc.validationFunc != nil {
				require.True(t, tc.validationFunc(sts), "validation func should pass")
			}
		})
	}
}

func TestEditItem(t *testing.T) {
	mem, _ := memdb.New()
	sts := &StatusThingService{store: mem}
	ctx := context.TODO()

	res, err := sts.NewItem(ctx, t.Name())
	require.NoError(t, err)
	require.NotNil(t, res)

	eerr := sts.EditItem(ctx, res.GetId(), filters.WithDescription("desc"))
	require.NoError(t, eerr)

	gres, gerr := sts.GetItem(ctx, res.GetId())
	require.NoError(t, gerr)
	require.NotNil(t, gres)
	require.Equal(t, "desc", gres.GetDescription())
}

func TestAddItem(t *testing.T) {
	t.Run("nil-store", func(t *testing.T) {
		sts := &StatusThingService{}
		_, err := sts.NewItem(context.TODO(), "")
		require.ErrorIs(t, err, errors.ErrStoreUnavailable)
	})
	t.Run("happy-path", func(t *testing.T) {
		sts, err := NewStatusThingService(&testStatusThingStore{
			customStatuses: []*statusthingv1.Status{},
		})
		require.NoError(t, err, "should not error")
		require.NotNil(t, sts, "should not be nil")
		res, err := sts.NewItem(context.TODO(), t.Name())
		require.NoError(t, err, "should not error")
		require.NotNil(t, res, "should not be nil")
		require.NotEmpty(t, res.GetId(), "id should be populated")
	})
	t.Run("happy-path-with-custom-id", func(t *testing.T) {
		sts, err := NewStatusThingService(&testStatusThingStore{
			customStatuses: []*statusthingv1.Status{},
		})
		require.NoError(t, err, "should not error")
		require.NotNil(t, sts, "should not be nil")
		res, err := sts.NewItem(context.TODO(), t.Name(), filters.WithItemID(t.Name()))
		require.NoError(t, err, "should not error")
		require.NotNil(t, res, "should not be nil")
		require.Equal(t, t.Name(), res.GetId())
	})
	t.Run("happy-path-with-custom-status", func(t *testing.T) {
		sts, err := NewStatusThingService(&testStatusThingStore{
			customStatuses: []*statusthingv1.Status{
				{Id: t.Name() + "_status_id"},
			},
		})
		require.NoError(t, err, "should not error")
		require.NotNil(t, sts, "should not be nil")
		res, err := sts.NewItem(context.TODO(), t.Name(), filters.WithStatusID(t.Name()+"_status_id"))
		require.NoError(t, err, "should not error")
		require.NotNil(t, res, "should not be nil")
		require.Equal(t, t.Name()+"_status_id", res.GetStatus().GetId())
	})
	t.Run("happy-path-with-custom-description", func(t *testing.T) {
		sts, err := NewStatusThingService(&testStatusThingStore{
			customStatuses: []*statusthingv1.Status{},
		})
		require.NoError(t, err, "should not error")
		require.NotNil(t, sts, "should not be nil")
		res, err := sts.NewItem(context.TODO(), t.Name(), filters.WithDescription(t.Name()))
		require.NoError(t, err, "should not error")
		require.NotNil(t, res, "should not be nil")
		require.Equal(t, t.Name(), res.GetDescription())
	})
}

func TestRemoveItem(t *testing.T) {
	t.Run("nil-store", func(t *testing.T) {
		sts := &StatusThingService{}
		err := sts.DeleteItem(context.TODO(), "")
		require.ErrorIs(t, err, errors.ErrStoreUnavailable)
	})
	t.Run("happy-path", func(t *testing.T) {
		sts, err := NewStatusThingService(&testStatusThingStore{
			customStatuses: []*statusthingv1.Status{},
		})
		require.NoError(t, err, "should not error")
		require.NotNil(t, sts, "should not be nil")
		err = sts.DeleteItem(context.TODO(), t.Name())
		require.NoError(t, err, "should not error")
	})
	t.Run("missing-thing-id", func(t *testing.T) {
		sts, err := NewStatusThingService(&testStatusThingStore{
			customStatuses: []*statusthingv1.Status{},
		})
		require.NoError(t, err, "should not error")
		require.NotNil(t, sts, "should not be nil")
		err = sts.DeleteItem(context.TODO(), "")
		require.ErrorIs(t, err, errors.ErrEmptyString)
	})
}

func TestAllItems(t *testing.T) {
	ctx := context.TODO()
	t.Run("no-opts", func(t *testing.T) {
		store, _ := memdb.New()
		svc, _ := NewStatusThingService(store, WithDefaults())
		require.NotNil(t, svc)

		for _, status := range svc.GetCreatedDefaults() {
			ires, ierr := svc.NewItem(ctx, status.GetName(), filters.WithStatus(status))
			require.NoError(t, ierr)
			require.NotNil(t, ires)

			idres, iderr := svc.FindItems(ctx, filters.WithStatusIDs(status.GetId()))
			require.NoError(t, iderr)
			require.Len(t, idres, 1)
			item := idres[0]
			require.Equal(t, ires.GetId(), item.GetId())
		}
	})
	t.Run("status-ids", func(t *testing.T) {

	})
	t.Run("status-kinds", func(t *testing.T) {

	})
}
