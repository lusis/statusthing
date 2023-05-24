package services

import (
	"context"
	"testing"
	"time"

	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/filters"
	"github.com/lusis/statusthing/internal/serrors"
	"github.com/lusis/statusthing/internal/storers/memdb"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestAddNote(t *testing.T) {
	t.Run("nil-store", func(t *testing.T) {
		sts := &StatusThingService{}
		_, err := sts.NewNote(context.TODO(), "", "")
		require.ErrorIs(t, err, serrors.ErrStoreUnavailable)
	})
	t.Run("happy-path", func(t *testing.T) {
		sts, err := NewStatusThingService(&testStatusThingStore{
			customStatuses: []*statusthingv1.Status{},
		})
		require.NoError(t, err, "should not error")
		require.NotNil(t, sts, "should not be nil")
		res, err := sts.NewNote(context.TODO(), t.Name(), t.Name())
		require.NoError(t, err, "should not error")
		require.NotNil(t, res)
	})
	t.Run("missing-item-id", func(t *testing.T) {
		sts, err := NewStatusThingService(&testStatusThingStore{
			customStatuses: []*statusthingv1.Status{},
		})
		require.NoError(t, err, "should not error")
		require.NotNil(t, sts, "should not be nil")
		res, err := sts.NewNote(context.TODO(), "", t.Name())
		require.ErrorIs(t, err, serrors.ErrEmptyString)
		require.Nil(t, res)
	})
	t.Run("missing-note-text", func(t *testing.T) {
		sts, err := NewStatusThingService(&testStatusThingStore{
			customStatuses: []*statusthingv1.Status{},
		})
		require.NoError(t, err, "should not error")
		require.NotNil(t, sts, "should not be nil")
		res, err := sts.NewNote(context.TODO(), t.Name(), "")
		require.ErrorIs(t, err, serrors.ErrEmptyString)
		require.Nil(t, res)
	})
}

func TestEditNote(t *testing.T) {
	t.Run("nil-store", func(t *testing.T) {
		sts := &StatusThingService{}
		err := sts.EditNote(context.TODO(), "", "")
		require.ErrorIs(t, err, serrors.ErrStoreUnavailable)
	})
	t.Run("happy-path", func(t *testing.T) {
		sts, err := NewStatusThingService(&testStatusThingStore{
			customStatuses: []*statusthingv1.Status{},
		})
		require.NoError(t, err, "should not error")
		require.NotNil(t, sts, "should not be nil")
		err = sts.EditNote(context.TODO(), t.Name(), t.Name())
		require.NoError(t, err, "should not error")
	})
	t.Run("happy-path-custom-timestamps", func(t *testing.T) {
		ts := makeTsNow()
		deleted := time.Now().Add(-24 * time.Hour)
		ts.Deleted = timestamppb.New(deleted)
		testStore := &testStatusThingStore{
			customStatuses: []*statusthingv1.Status{},
		}
		sts, err := NewStatusThingService(testStore)
		require.NoError(t, err, "should not error")
		require.NotNil(t, sts, "should not be nil")
		err = sts.EditNote(context.TODO(), t.Name(), t.Name(), filters.WithTimestamps(ts))
		require.NoError(t, err, "should not error")
		require.NotNil(t, testStore.lastUpdatedNote, "last note should be updated")
		require.Equal(t, ts, testStore.lastUpdatedNote.GetTimestamps(), "timestamps should be equal")
	})
	t.Run("missing-thing-id", func(t *testing.T) {
		sts, err := NewStatusThingService(&testStatusThingStore{
			customStatuses: []*statusthingv1.Status{},
		})
		require.NoError(t, err, "should not error")
		require.NotNil(t, sts, "should not be nil")
		err = sts.EditNote(context.TODO(), "", t.Name())
		require.ErrorIs(t, err, serrors.ErrEmptyString)
	})
	t.Run("missing-note-text", func(t *testing.T) {
		sts, err := NewStatusThingService(&testStatusThingStore{
			customStatuses: []*statusthingv1.Status{},
		})
		require.NoError(t, err, "should not error")
		require.NotNil(t, sts, "should not be nil")
		err = sts.EditNote(context.TODO(), t.Name(), "")
		require.ErrorIs(t, err, serrors.ErrEmptyString)
	})
}
func TestRemoveNote(t *testing.T) {
	t.Run("nil-store", func(t *testing.T) {
		sts := &StatusThingService{}
		err := sts.DeleteNote(context.TODO(), "")
		require.ErrorIs(t, err, serrors.ErrStoreUnavailable)
	})
	t.Run("happy-path", func(t *testing.T) {
		sts, err := NewStatusThingService(&testStatusThingStore{
			customStatuses: []*statusthingv1.Status{},
		})
		require.NoError(t, err, "should not error")
		require.NotNil(t, sts, "should not be nil")
		err = sts.DeleteNote(context.TODO(), t.Name())
		require.NoError(t, err, "should not error")
	})
	t.Run("missing-note-id", func(t *testing.T) {
		sts, err := NewStatusThingService(&testStatusThingStore{
			customStatuses: []*statusthingv1.Status{},
		})
		require.NoError(t, err, "should not error")
		require.NotNil(t, sts, "should not be nil")
		err = sts.DeleteNote(context.TODO(), "")
		require.ErrorIs(t, err, serrors.ErrEmptyString)
	})
}

func TestAllNotes(t *testing.T) {
	ctx := context.TODO()
	t.Run("nil-store", func(t *testing.T) {
		sts := &StatusThingService{}
		res, err := sts.AllNotes(ctx, t.Name())
		require.ErrorIs(t, err, serrors.ErrStoreUnavailable)
		require.Nil(t, res)
	})
	t.Run("missing-itemid", func(t *testing.T) {
		mem, _ := memdb.New()
		sts := &StatusThingService{store: mem}
		res, err := sts.AllNotes(ctx, "")
		require.ErrorIs(t, err, serrors.ErrEmptyString)
		require.Nil(t, res)
	})
}

func TestGetNote(t *testing.T) {
	ctx := context.TODO()
	t.Run("nil-store", func(t *testing.T) {
		sts := &StatusThingService{}
		res, err := sts.GetNote(ctx, t.Name())
		require.ErrorIs(t, err, serrors.ErrStoreUnavailable)
		require.Nil(t, res)
	})
	t.Run("missing-itemid", func(t *testing.T) {
		mem, _ := memdb.New()
		sts := &StatusThingService{store: mem}
		res, err := sts.GetNote(ctx, "")
		require.ErrorIs(t, err, serrors.ErrEmptyString)
		require.Nil(t, res)
	})
	t.Run("happy-path", func(t *testing.T) {
		mem, _ := memdb.New()
		sts := &StatusThingService{store: mem}
		ires, ierr := sts.NewItem(ctx, t.Name(), filters.WithNoteText(t.Name()))
		require.NoError(t, ierr)
		require.NotNil(t, ires)
		require.Len(t, ires.GetNotes(), 1)
		note := ires.GetNotes()[0]
		gres, gerr := sts.GetNote(ctx, note.GetId())
		require.NoError(t, gerr)
		require.NotNil(t, gres)
		require.Equal(t, note.GetId(), gres.GetId())
		require.Equal(t, note.GetText(), gres.GetText())
	})
}
