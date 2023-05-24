package memdb

import (
	"context"
	"testing"

	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/errors"
	"github.com/lusis/statusthing/internal/filters"
	"github.com/lusis/statusthing/internal/testutils"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestDbNoteFromProto(t *testing.T) {
	pbNow := timestamppb.Now()
	intNow := tsToInt(pbNow)
	testcases := map[string]struct {
		protoNote      *statusthingv1.Note
		itemID         string
		dbNote         *dbNote
		err            error
		validationFunc func(*testing.T, *dbNote, *statusthingv1.Note)
	}{
		"nil-note": {
			err: errors.ErrNilVal,
		},
		"missing-item-id": {
			protoNote: &statusthingv1.Note{},
			err:       errors.ErrEmptyString,
		},
		"missing-note-id": {
			protoNote: &statusthingv1.Note{},
			itemID:    t.Name(),
			err:       errors.ErrEmptyString,
		},
		"missing-note-text": {
			protoNote: &statusthingv1.Note{Id: t.Name()},
			itemID:    t.Name(),
			err:       errors.ErrEmptyString,
		},
		"missing-timestamps": {
			protoNote: &statusthingv1.Note{Id: t.Name(), Text: t.Name()},
			itemID:    t.Name(),
			err:       errors.ErrNilVal,
		},
		"invalid-created": {
			protoNote: &statusthingv1.Note{Id: t.Name(), Text: t.Name(), Timestamps: &statusthingv1.Timestamps{Updated: pbNow}},
			itemID:    t.Name(),
			err:       errors.ErrInvalidData,
		},
		"invalid-updated": {
			protoNote: &statusthingv1.Note{Id: t.Name(), Text: t.Name(), Timestamps: &statusthingv1.Timestamps{Created: pbNow}},
			itemID:    t.Name(),
			err:       errors.ErrInvalidData,
		},
		"happy-path": {
			protoNote: &statusthingv1.Note{Id: t.Name(), Text: t.Name(), Timestamps: &statusthingv1.Timestamps{Created: pbNow, Updated: pbNow, Deleted: pbNow}},
			itemID:    "my-item-id",
			dbNote: &dbNote{
				ID:       t.Name(),
				NoteData: t.Name(),
				ItemID:   "my-item-id",
				Created:  intNow,
				Updated:  intNow,
				Deleted:  intNow,
			},
			validationFunc: func(t *testing.T, dn *dbNote, n *statusthingv1.Note) {
				require.Equal(t, dn.ID, n.GetId())
				require.Equal(t, dn.NoteData, n.GetText())
				require.NotZero(t, dn.Created)
				require.Equal(t, intNow, dn.Created)
				require.NotZero(t, dn.Updated)
				require.Equal(t, intNow, dn.Updated)
				require.NotZero(t, dn.Deleted)
				require.Equal(t, intNow, dn.Deleted)
				require.Equal(t, "my-item-id", dn.ItemID)
			},
		},
	}
	t.Parallel()
	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			nres, nerr := noteFromProto(tc.protoNote, tc.itemID)
			if tc.err != nil {
				require.ErrorIs(t, nerr, tc.err)
				require.Nil(t, nres)
			} else {
				require.NoError(t, nerr)
				require.NotNil(t, nres)
			}
			if tc.validationFunc != nil {
				tc.validationFunc(t, nres, tc.protoNote)
			}
		})
	}
}

func TestDbNoteToProto(t *testing.T) {
	pbNow := timestamppb.Now()
	intNow := tsToInt(pbNow)
	t.Parallel()
	testcases := map[string]struct {
		dbNote         *dbNote
		protoNote      *statusthingv1.Note
		err            error
		validationFunc func(*testing.T, *dbNote, *statusthingv1.Note)
	}{
		"missing-id": {
			dbNote: &dbNote{},
			err:    errors.ErrInvalidData,
		},
		"missing-note-data": {
			dbNote: &dbNote{ID: t.Name()},
			err:    errors.ErrInvalidData,
		},
		"missing-created": {
			dbNote: &dbNote{ID: t.Name(), NoteData: t.Name(), Updated: intNow},
			err:    errors.ErrInvalidData,
		},
		"missing-updated": {
			dbNote: &dbNote{ID: t.Name(), NoteData: t.Name(), Created: intNow},
			err:    errors.ErrInvalidData,
		},
		"happy-path": {
			dbNote: &dbNote{ID: t.Name(), NoteData: t.Name(), Updated: intNow, Created: intNow, Deleted: intNow},
			validationFunc: func(t *testing.T, dn *dbNote, n *statusthingv1.Note) {
				require.Equal(t, dn.ID, n.GetId())
				require.Equal(t, dn.NoteData, n.GetText())
				require.Equal(t, pbNow, n.GetTimestamps().GetCreated())
				require.Equal(t, pbNow, n.GetTimestamps().GetDeleted())
				require.Equal(t, pbNow, n.GetTimestamps().GetUpdated())
			},
		},
	}
	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			nres, nerr := tc.dbNote.toProto()
			if tc.err != nil {
				require.ErrorIs(t, nerr, tc.err)
				require.Nil(t, nres)
			} else {
				require.NoError(t, nerr)
				require.NotNil(t, nres)
				if tc.validationFunc != nil {
					tc.validationFunc(t, tc.dbNote, nres)
				}
			}
		})
	}
}
func TestNotesLifeCycle(t *testing.T) {
	store, _ := New()
	note := testutils.MakeNote(t.Name())

	// we have to have an item in the db to add the note
	item := testutils.MakeItem(t.Name())

	ires, ierr := store.StoreItem(context.TODO(), item)
	require.NoError(t, ierr)
	require.NotNil(t, ires)

	res, err := store.StoreNote(context.TODO(), note, item.GetId())
	require.NoError(t, err)
	require.NotNil(t, res)
	require.True(t, proto.Equal(item, ires), "original note and stored note should be the same")

	check, checkErr := store.GetNote(context.TODO(), res.GetId())
	require.NoError(t, checkErr)
	require.NotNil(t, check)
	require.True(t, proto.Equal(res, check))

	checkAll, checkAllErr := store.FindNotes(context.TODO(), ires.GetId())
	require.NoError(t, checkAllErr)
	require.Len(t, checkAll, 1)
	first := checkAll[0]
	require.True(t, proto.Equal(res, first), "stored note and found note should be the same")

	updateErr := store.UpdateNote(context.TODO(), first.GetId(), filters.WithNoteText("new_text"))
	require.NoError(t, updateErr)

	checkAgain, againErr := store.FindNotes(context.TODO(), ires.GetId())
	require.NoError(t, againErr)
	require.Len(t, checkAgain, 1)
	firstAgain := checkAgain[0]
	require.Equal(t, note.GetId(), firstAgain.GetId())
	require.Equal(t, "new_text", firstAgain.GetText())
	require.Greater(t, firstAgain.GetTimestamps().GetUpdated().AsTime(), note.GetTimestamps().GetUpdated().AsTime(), "updated timestamp should be updated")

	delErr := store.DeleteNote(context.TODO(), firstAgain.GetId())
	require.NoError(t, delErr)

	delCheck, delCheckErr := store.FindNotes(context.TODO(), ires.GetId())
	require.NoError(t, delCheckErr)
	require.Len(t, delCheck, 0)
}

func TestStoreNoteErrors(t *testing.T) {
	t.Parallel()
	type testcase struct {
		note   *statusthingv1.Note
		itemID string
		err    error
	}
	testcases := map[string]testcase{
		"nil-note": {
			note:   nil,
			itemID: "foo",
			err:    errors.ErrNilVal,
		},
		"missing-item-id": {
			note: &statusthingv1.Note{},
			err:  errors.ErrEmptyString,
		},
	}
	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			mem, _ := New()
			res, err := mem.StoreNote(context.TODO(), tc.note, tc.itemID)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				require.Nil(t, res)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.note.GetId(), res.GetId())
				require.Equal(t, tc.note.GetText(), res.GetText())
				require.True(t, testutils.TimestampsEqual(tc.note.GetTimestamps(), res.GetTimestamps()))
			}
		})
	}
}
