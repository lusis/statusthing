package services

import (
	"context"
	"testing"

	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/errors"
	"github.com/lusis/statusthing/internal/filters"
	"github.com/lusis/statusthing/internal/storers/memdb"
	"github.com/stretchr/testify/require"
)

func TestAddStatus(t *testing.T) {
	t.Run("nil-store", func(t *testing.T) {
		sts := &StatusThingService{}
		_, err := sts.NewStatus(context.TODO(), "", 0)
		require.ErrorIs(t, err, errors.ErrStoreUnavailable)
	})
	t.Run("happy-path", func(t *testing.T) {
		sts, err := NewStatusThingService(&testStatusThingStore{
			customStatuses: []*statusthingv1.Status{},
		})
		require.NoError(t, err, "should not error")
		require.NotNil(t, sts, "should not be nil")
		s, err := sts.NewStatus(context.TODO(), t.Name(), statusthingv1.StatusKind_STATUS_KIND_AVAILABLE)
		require.NoError(t, err, "should not error")
		require.NotNil(t, s, "should not be nil")
	})

	t.Run("happy-path-custom-id", func(t *testing.T) {
		store := &testStatusThingStore{
			customStatuses: []*statusthingv1.Status{},
		}
		sts, err := NewStatusThingService(store)
		require.NoError(t, err, "should not error")
		require.NotNil(t, sts, "should not be nil")
		statusID := t.Name() + "_status_id"
		s, err := sts.NewStatus(context.TODO(), t.Name(), statusthingv1.StatusKind_STATUS_KIND_AVAILABLE, filters.WithStatusID(statusID))
		require.NoError(t, err, "should not error")
		require.NotNil(t, s, "should not be nil")
		// we use constructor so items will be added at instantiation
		lastElem := store.customStatuses[len(store.customStatuses)-1]
		require.Equal(t, statusID, lastElem.GetId(), "id should be set")
	})
	t.Run("happy-path-custom-description", func(t *testing.T) {
		store := &testStatusThingStore{
			customStatuses: []*statusthingv1.Status{},
		}
		sts, err := NewStatusThingService(store)
		require.NoError(t, err, "should not error")
		require.NotNil(t, sts, "should not be nil")
		s, err := sts.NewStatus(context.TODO(), t.Name(), statusthingv1.StatusKind_STATUS_KIND_AVAILABLE, filters.WithDescription(t.Name()))
		require.NoError(t, err, "should not error")
		require.NotNil(t, s, "should not be nil")
		// we use constructor so items will be added at instantiation
		lastElem := store.customStatuses[len(store.customStatuses)-1]
		require.Equal(t, t.Name(), lastElem.GetDescription(), "description should be set")
	})
	t.Run("happy-path-color", func(t *testing.T) {
		store := &testStatusThingStore{
			customStatuses: []*statusthingv1.Status{},
		}
		sts, err := NewStatusThingService(store)
		require.NoError(t, err, "should not error")
		require.NotNil(t, sts, "should not be nil")
		s, err := sts.NewStatus(context.TODO(), t.Name(), statusthingv1.StatusKind_STATUS_KIND_AVAILABLE, filters.WithColor(t.Name()))
		require.NoError(t, err, "should not error")
		require.NotNil(t, s, "should not be nil")
		// we use constructor so items will be added at instantiation
		lastElem := store.customStatuses[len(store.customStatuses)-1]
		require.Equal(t, t.Name(), lastElem.GetColor(), "color should be set")
	})
	t.Run("missing-status-name", func(t *testing.T) {
		sts, err := NewStatusThingService(&testStatusThingStore{
			customStatuses: []*statusthingv1.Status{},
		})
		require.NoError(t, err, "should not error")
		require.NotNil(t, sts, "should not be nil")
		s, err := sts.NewStatus(context.TODO(), "", statusthingv1.StatusKind_STATUS_KIND_AVAILABLE)
		require.ErrorIs(t, err, errors.ErrEmptyString)
		require.Nil(t, s, "should be nil")
	})
	t.Run("zero-value-enum", func(t *testing.T) {
		sts, err := NewStatusThingService(&testStatusThingStore{
			customStatuses: []*statusthingv1.Status{},
		})
		require.NoError(t, err, "should not error")
		require.NotNil(t, sts, "should not be nil")
		s, err := sts.NewStatus(context.TODO(), t.Name(), 0)
		require.ErrorIs(t, err, errors.ErrEmptyEnum)
		require.Nil(t, s, "should be nil")
	})
}
func TestRemoveStatus(t *testing.T) {
	t.Run("nil-store", func(t *testing.T) {
		sts := &StatusThingService{}
		err := sts.DeleteStatus(context.TODO(), t.Name())
		require.ErrorIs(t, err, errors.ErrStoreUnavailable)
	})
	t.Run("empty-id", func(t *testing.T) {
		sts := &StatusThingService{store: &testStatusThingStore{}}
		err := sts.DeleteStatus(context.TODO(), "")
		require.ErrorIs(t, err, errors.ErrEmptyString)
	})
	t.Run("happy-path", func(t *testing.T) {
		sts, err := NewStatusThingService(&testStatusThingStore{
			customStatuses: []*statusthingv1.Status{},
		})
		require.NoError(t, err, "should not error")
		require.NotNil(t, sts, "should not be nil")
		err = sts.DeleteStatus(context.TODO(), t.Name())
		require.NoError(t, err, "should not error")
	})
}

func TestEditStatus(t *testing.T) {
	ctx := context.TODO()
	t.Run("nil-store", func(t *testing.T) {
		sts := &StatusThingService{}
		err := sts.EditStatus(context.TODO(), t.Name())
		require.ErrorIs(t, err, errors.ErrStoreUnavailable)
	})
	t.Run("empty-id", func(t *testing.T) {
		sts := &StatusThingService{store: &testStatusThingStore{}}
		err := sts.EditStatus(context.TODO(), "", filters.WithColor("red"))
		require.ErrorIs(t, err, errors.ErrEmptyString)
	})
	t.Run("no-filters", func(t *testing.T) {
		sts := &StatusThingService{store: &testStatusThingStore{}}
		err := sts.EditStatus(context.TODO(), t.Name())
		require.ErrorIs(t, err, errors.ErrAtLeastOne)
	})
	t.Run("happy-path", func(t *testing.T) {
		mem, _ := memdb.New()
		sts, _ := NewStatusThingService(mem, WithDefaults())

		for _, s := range sts.GetCreatedDefaults() {
			rerr := sts.EditStatus(ctx, s.GetId(), filters.WithDescription(s.GetKind().String()))
			require.NoError(t, rerr)

			gres, gerr := sts.GetStatus(ctx, s.GetId())
			require.NoError(t, gerr)
			require.Equal(t, s.GetKind().String(), gres.GetDescription())
		}
	})
}
