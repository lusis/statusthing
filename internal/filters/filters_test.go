package filters

import (
	"testing"

	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/errors"
	"github.com/lusis/statusthing/internal/testutils"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// simpleTestCase is meant for very simple validation of a set of FilterOption
type simpleTestCase struct {
	// the type of error, if any, that should be returned
	err error
	// the filter options to be provided
	opts []FilterOption
	// validation func is a func that will be run to allow test code to be simpler
	validationFunc func(*Filters)
}

func TestNew(t *testing.T) {
	t.Parallel()
	protonow := timestamppb.Now()
	happyTs := &statusthingv1.Timestamps{Created: protonow}
	cases := map[string]simpleTestCase{
		"itemID-happy-path": {
			opts:           []FilterOption{WithItemID(t.Name())},
			validationFunc: func(f *Filters) { require.Equal(t, t.Name(), f.ItemID()) },
		},
		"itemID-empty": {
			opts: []FilterOption{WithItemID("")},
			err:  errors.ErrEmptyString,
		},
		"itemID-already-set": {
			opts: []FilterOption{WithItemID(t.Name()), WithItemID(t.Name())},
			err:  errors.ErrAlreadySet,
		},
		"noteID-happy-path": {
			opts:           []FilterOption{WithNoteID(t.Name())},
			validationFunc: func(f *Filters) { require.Equal(t, t.Name(), f.NoteID()) },
		},
		"noteID-empty": {
			opts: []FilterOption{WithNoteID("")},
			err:  errors.ErrEmptyString,
		},
		"noteID-already-set": {
			opts: []FilterOption{WithNoteID(t.Name()), WithNoteID(t.Name())},
			err:  errors.ErrAlreadySet,
		},
		"statusID-status-already-set": {
			opts: []FilterOption{WithStatus(&statusthingv1.Status{}), WithStatusID(t.Name())},
			err:  errors.ErrAlreadySet,
		},
		"statusID-happy-path": {
			opts:           []FilterOption{WithStatusID(t.Name())},
			validationFunc: func(f *Filters) { require.Equal(t, t.Name(), f.StatusID()) },
		},
		"statusID-empty": {
			opts: []FilterOption{WithStatusID("")},
			err:  errors.ErrEmptyString,
		},
		"statusID-already-set": {
			opts: []FilterOption{WithStatusID(t.Name()), WithStatusID(t.Name())},
			err:  errors.ErrAlreadySet,
		},
		"status-statusid-conflict": {
			opts: []FilterOption{WithStatusID(t.Name()), WithStatus(testutils.MakeStatus(t.Name()))},
			err:  errors.ErrAlreadySet,
		},
		"status-happypath": {
			opts:           []FilterOption{WithStatus(&statusthingv1.Status{Id: t.Name()})},
			validationFunc: func(f *Filters) { require.Equal(t, t.Name(), f.Status().GetId()) },
		},
		"status-nil-status": {
			opts: []FilterOption{WithStatus(nil)},
			err:  errors.ErrNilVal,
		},
		"status-already-set": {
			opts: []FilterOption{WithStatus(&statusthingv1.Status{Id: t.Name()}), WithStatus(&statusthingv1.Status{Id: t.Name()})},
			err:  errors.ErrAlreadySet,
		},
		"color-happy-path": {
			opts:           []FilterOption{WithColor(t.Name())},
			validationFunc: func(f *Filters) { require.Equal(t, t.Name(), f.Color()) },
		},
		"color-empty": {
			opts: []FilterOption{WithColor("")},
			err:  errors.ErrEmptyString,
		},
		"color-already-set": {
			opts: []FilterOption{WithColor(t.Name()), WithColor(t.Name())},
			err:  errors.ErrAlreadySet,
		},
		"description-happy-path": {
			opts:           []FilterOption{WithDescription(t.Name())},
			validationFunc: func(f *Filters) { require.Equal(t, t.Name(), f.Description()) },
		},
		"description-empty": {
			opts: []FilterOption{WithDescription("")},
			err:  errors.ErrEmptyString,
		},
		"description-already-set": {
			opts: []FilterOption{WithDescription(t.Name()), WithDescription(t.Name())},
			err:  errors.ErrAlreadySet,
		},
		"statuskind-happy-path": {
			opts:           []FilterOption{WithStatusKind(statusthingv1.StatusKind_STATUS_KIND_AVAILABLE)},
			validationFunc: func(f *Filters) { require.Equal(t, statusthingv1.StatusKind_STATUS_KIND_AVAILABLE, f.StatusKind()) },
		},
		"statuskind-zero-val": {
			opts: []FilterOption{WithStatusKind(statusthingv1.StatusKind_STATUS_KIND_UNKNOWN)},
			err:  errors.ErrEmptyEnum,
		},
		"statuskind-already-set": {
			opts: []FilterOption{WithStatusKind(statusthingv1.StatusKind_STATUS_KIND_AVAILABLE), WithStatusKind(statusthingv1.StatusKind_STATUS_KIND_DOWN)},
			err:  errors.ErrAlreadySet,
		},
		"timestamps-happy-path": {
			opts:           []FilterOption{WithTimestamps(happyTs)},
			validationFunc: func(f *Filters) { require.Equal(t, happyTs, f.Timestamps()) },
		},
		"timestamps-zero-val": {
			opts: []FilterOption{WithTimestamps(nil)},
			err:  errors.ErrNilVal,
		},
		"timestamps-already-set": {
			opts: []FilterOption{WithTimestamps(&statusthingv1.Timestamps{Created: protonow}), WithTimestamps(&statusthingv1.Timestamps{Created: timestamppb.Now()})},
			err:  errors.ErrAlreadySet,
		},
		"statusids-happy-path": {
			opts:           []FilterOption{WithStatusIDs("1", "2")},
			validationFunc: func(f *Filters) { require.Equal(t, []string{"1", "2"}, f.StatusIDs()) },
		},
		"statusids-atleastone": {
			opts: []FilterOption{WithStatusIDs()},
			err:  errors.ErrAtLeastOne,
		},
		"statusids-already-set": {
			opts: []FilterOption{WithStatusIDs(t.Name()), WithStatusIDs(t.Name())},
			err:  errors.ErrAlreadySet,
		},
		"name-happy-path": {
			opts:           []FilterOption{WithName("1")},
			validationFunc: func(f *Filters) { require.Equal(t, "1", f.Name()) },
		},
		"name-emptystring": {
			opts: []FilterOption{WithName("")},
			err:  errors.ErrEmptyString,
		},
		"name-already-set": {
			opts: []FilterOption{WithName(t.Name()), WithName(t.Name())},
			err:  errors.ErrAlreadySet,
		},
		"note-text-happy-path": {
			opts:           []FilterOption{WithNoteText("1")},
			validationFunc: func(f *Filters) { require.Equal(t, "1", f.NoteText()) },
		},
		"note-text-emptystring": {
			opts: []FilterOption{WithNoteText("")},
			err:  errors.ErrEmptyString,
		},
		"note-text-already-set": {
			opts: []FilterOption{WithNoteText(t.Name()), WithNoteText(t.Name())},
			err:  errors.ErrAlreadySet,
		},
		"statuskinds-happy-path": {
			opts:           []FilterOption{WithStatusKinds(statusthingv1.StatusKind_STATUS_KIND_AVAILABLE)},
			validationFunc: func(f *Filters) { require.Equal(t, statusthingv1.StatusKind_STATUS_KIND_AVAILABLE, f.StatusKinds()[0]) },
		},
		"statuskinds-atleastone": {
			opts: []FilterOption{WithStatusKinds()},
			err:  errors.ErrAtLeastOne,
		},
		"statuskinds-enum-unknown": {
			opts: []FilterOption{WithStatusKinds(statusthingv1.StatusKind_STATUS_KIND_UNKNOWN)},
			err:  errors.ErrEmptyEnum,
		},
		"statuskinds-already-set": {
			opts: []FilterOption{WithStatusKinds(statusthingv1.StatusKind_STATUS_KIND_AVAILABLE), WithStatusKinds(statusthingv1.StatusKind_STATUS_KIND_AVAILABLE)},
			err:  errors.ErrAlreadySet,
		},
	}
	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			res, err := New(tc.opts...)
			if err != nil {
				require.ErrorIs(t, err, tc.err)
				require.Nil(t, res, "error should return nil Filters")
			} else {
				require.NoError(t, err, "error should be nil")
				require.NotNil(t, res, "should not be nil")
			}
			if tc.validationFunc != nil {
				tc.validationFunc(res)
			}
		})
	}
}
