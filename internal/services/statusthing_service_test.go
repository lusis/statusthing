package services

import (
	"context"
	"fmt"
	"testing"

	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/errors"
	"github.com/lusis/statusthing/internal/filters"
	"github.com/lusis/statusthing/internal/storers/memdb"
	"github.com/lusis/statusthing/internal/storers/unimplemented"
	"github.com/stretchr/testify/require"
)

var testErr = fmt.Errorf("snarf")

func TestNew(t *testing.T) {
	t.Run("nil-store", func(t *testing.T) {
		sts, err := NewStatusThingService(nil)
		require.ErrorIs(t, err, errors.ErrStoreUnavailable)
		require.Nil(t, sts, "should be nil")
	})

	var allStatusErr = fmt.Errorf("all-error")
	t.Run("all-status-error", func(t *testing.T) {
		sts, err := NewStatusThingService(&testStatusThingStore{
			allErr: allStatusErr,
		}, WithDefaults())
		require.ErrorIs(t, err, allStatusErr)
		require.Nil(t, sts, "should be nil")
	})
	var cStoreErr = fmt.Errorf("cstore-error")
	t.Run("c-store-error", func(t *testing.T) {
		sts, err := NewStatusThingService(&testStatusThingStore{
			cstoreErr: cStoreErr,
		})
		require.NoError(t, err, "should not error inserting individual items")
		require.NotNil(t, sts, "should not be nil")
	})
	t.Run("happy-path", func(t *testing.T) {

		sts, err := NewStatusThingService(&testStatusThingStore{
			customStatuses: []*statusthingv1.Status{},
		})
		require.NoError(t, err, "should not error")
		require.NotNil(t, sts, "should not be nil")
	})
	t.Run("with-default-statuses", func(t *testing.T) {
		mdb, err := memdb.New()
		require.NoError(t, err)
		require.NotNil(t, mdb)
		sts, err := NewStatusThingService(mdb, WithDefaults())
		require.NoError(t, err)
		require.NotNil(t, sts)
		all, allerr := sts.AllStatuses(context.TODO())
		require.NoError(t, allerr)
		require.NotNil(t, all)
		require.Len(t, all, 3)
	})
}

type testStatusThingStore struct {
	*unimplemented.StatusThingStore
	customStatuses  []*statusthingv1.Status
	things          []*statusthingv1.Item
	notes           map[string][]*statusthingv1.Note
	lastUpdatedNote *statusthingv1.Note
	// the following are errors returned by each function if set
	// AllStatuses
	allErr error
	// StoreStatus
	cstoreErr error
	// Store
	addErr error
	// GetStatus
	getStatusErr error
	// Update
	updateErr error
	// Delete
	deleteErr error
	// DeleteStatus
	cdelErr error
	// StoreNote
	sNoteErr error
	// UpdateNote
	updateNoteErr error
	// DeleteNote
	delNoteErr error
}

func (ts *testStatusThingStore) UpdateItem(_ context.Context, _ string, _ ...filters.FilterOption) error {
	return ts.updateErr
}

func (ts *testStatusThingStore) DeleteItem(_ context.Context, _ string) error {
	return ts.deleteErr
}
func (ts *testStatusThingStore) FindStatus(_ context.Context, _ ...filters.FilterOption) ([]*statusthingv1.Status, error) {
	if ts.allErr != nil {
		return nil, ts.allErr
	}
	return ts.customStatuses, nil
}

// this is not a fancy in-memory store. We're only going to return the first item
func (ts *testStatusThingStore) GetStatus(_ context.Context, _ string) (*statusthingv1.Status, error) {
	if ts.getStatusErr != nil {
		return nil, ts.getStatusErr
	}
	if len(ts.customStatuses) == 0 {
		return nil, errors.ErrNotFound
	}
	return ts.customStatuses[0], nil
}

func (ts *testStatusThingStore) StoreStatus(_ context.Context, s *statusthingv1.Status) (*statusthingv1.Status, error) {
	if ts.cstoreErr != nil {
		return nil, ts.cstoreErr
	}
	ts.customStatuses = append(ts.customStatuses, s)
	return s, nil
}

func (ts *testStatusThingStore) DeleteStatus(_ context.Context, _ string) error {
	return ts.cdelErr
}

func (ts *testStatusThingStore) StoreItem(_ context.Context, item *statusthingv1.Item) (*statusthingv1.Item, error) {
	if ts.addErr != nil {
		return nil, ts.addErr
	}
	ts.things = append(ts.things, item)
	return item, nil
}

func (ts *testStatusThingStore) StoreNote(_ context.Context, note *statusthingv1.Note, thingID string) (*statusthingv1.Note, error) {
	if ts.sNoteErr != nil {
		return nil, ts.sNoteErr
	}
	if len(ts.notes) == 0 {
		ts.notes = map[string][]*statusthingv1.Note{}
	}
	if len(ts.notes[thingID]) == 0 {
		ts.notes[thingID] = make([]*statusthingv1.Note, 0)
	}
	ts.notes[thingID] = append(ts.notes[thingID], note)
	return note, nil
}

func (ts *testStatusThingStore) UpdateNote(_ context.Context, noteID string, opts ...filters.FilterOption) error {
	if ts.updateNoteErr != nil {
		return ts.updateNoteErr
	}
	// we need to parse here for some checks
	f, err := filters.New(opts...)
	if err != nil {
		return err
	}
	ts.lastUpdatedNote = &statusthingv1.Note{
		Id:         noteID,
		Timestamps: f.Timestamps(),
		Text:       f.NoteText(),
	}
	return nil
}

func (ts *testStatusThingStore) DeleteNote(_ context.Context, _ string) error {
	return ts.delNoteErr
}
