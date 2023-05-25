package internal

import (
	"testing"

	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/serrors"
	"github.com/lusis/statusthing/internal/storers"
	"github.com/lusis/statusthing/internal/testutils"
	"github.com/lusis/statusthing/internal/validation"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type noteTestCase struct {
	dbnote  *DbNote
	pbnote  *statusthingv1.Note
	errtext string
	err     error
}

func TestNoteFromProto(t *testing.T) {
	t.Parallel()

	t.Run("nil-check", func(t *testing.T) {
		s, serr := DbNoteFromProto(nil)
		require.ErrorIs(t, serr, serrors.ErrNilVal)
		require.Nil(t, s)
	})

	// t.Name inside this map refers to the name of the parent test not the iteration.
	// It's intentional since we're not worried about conflict here
	testcases := map[string]noteTestCase{
		"happy-path": {},
		"missing-id": {
			pbnote:  &statusthingv1.Note{Id: ""},
			err:     serrors.ErrEmptyString,
			errtext: "id",
		},
		"missing-txt": {
			pbnote:  &statusthingv1.Note{Id: t.Name()},
			err:     serrors.ErrEmptyString,
			errtext: "text",
		},
		"missing-timestamps": {
			pbnote:  &statusthingv1.Note{Id: t.Name(), Text: t.Name()},
			err:     serrors.ErrMissingTimestamp,
			errtext: "timestamps",
		},
		"missing-created": {
			pbnote: &statusthingv1.Note{Id: t.Name(), Text: t.Name(), Timestamps: &statusthingv1.Timestamps{
				Updated: timestamppb.Now(),
			}},
			err:     serrors.ErrMissingTimestamp,
			errtext: "created",
		},
		"missing-updated": {
			pbnote: &statusthingv1.Note{Id: t.Name(), Text: t.Name(), Timestamps: &statusthingv1.Timestamps{
				Created: timestamppb.Now(),
			}},
			err:     serrors.ErrMissingTimestamp,
			errtext: "updated",
		},
	}

	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			pb := tc.pbnote
			if pb == nil {
				pb = testutils.MakeNote(t.Name())
				pb.Timestamps = testutils.MakeTimestamps(true)
			}
			s, serr := DbNoteFromProto(pb)
			if tc.err != nil {
				require.ErrorIs(t, serr, tc.err)
				require.Nil(t, s)
				if validation.ValidString(tc.errtext) {
					require.ErrorContains(t, serr, tc.errtext)
				}
			} else {
				require.NoError(t, serr)
				require.NotNil(t, s)
				require.Equal(t, pb.GetId(), s.ID)
				require.Equal(t, pb.GetText(), s.NoteText)
				require.NotZero(t, s.Created)
				require.NotZero(t, s.Updated)
				if tc.pbnote.GetTimestamps().GetDeleted().IsValid() {
					require.NotNil(t, s.Deleted)
				}
			}
		})
	}
}

func TestNoteToProto(t *testing.T) {
	t.Parallel()
	// t.Name inside this map refers to the name of the parent test not the iteration.
	// It's intentional since we're not worried about conflict here
	testcases := map[string]noteTestCase{
		"happy-path": {},
		"missing-id": {
			dbnote:  &DbNote{ID: ""},
			err:     serrors.ErrInvalidData,
			errtext: "id",
		},
		"missing-text": {
			dbnote:  &DbNote{ID: t.Name()},
			err:     serrors.ErrInvalidData,
			errtext: "text",
		},
		"missing-created": {
			dbnote: &DbNote{
				ID:           t.Name(),
				NoteText:     t.Name(),
				DbTimestamps: &DbTimestamps{Updated: uint64(storers.TsToInt64(timestamppb.Now()))},
			},
			err:     serrors.ErrInvalidData,
			errtext: "created",
		},
		"missing-updated": {
			dbnote: &DbNote{
				ID:           t.Name(),
				NoteText:     t.Name(),
				DbTimestamps: &DbTimestamps{Created: uint64(storers.TsToInt64(timestamppb.Now()))},
			},
			err:     serrors.ErrInvalidData,
			errtext: "updated",
		},
	}

	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			pb := tc.dbnote
			if pb == nil {
				pb = &DbNote{
					ID:       t.Name(),
					NoteText: t.Name(),
					DbTimestamps: &DbTimestamps{
						Created: storers.TsToUInt64(timestamppb.Now()),
						Updated: storers.TsToUInt64(timestamppb.Now()),
						Deleted: storers.TsToUInt64Ptr(timestamppb.Now())},
					ItemID: t.Name(),
				}
			}
			s, serr := pb.ToProto()
			if tc.err != nil {
				require.ErrorIs(t, serr, tc.err)
				require.Nil(t, s)
				if validation.ValidString(tc.errtext) {
					require.ErrorContains(t, serr, tc.errtext)
				}
			} else {
				require.NoError(t, serr)
				require.NotNil(t, s)
				require.Equal(t, pb.ID, s.GetId())
				require.Equal(t, pb.NoteText, s.GetText())
				require.True(t, s.GetTimestamps().GetCreated().IsValid())
				require.True(t, s.GetTimestamps().GetUpdated().IsValid())
				if pb.Deleted != nil {
					require.True(t, s.GetTimestamps().GetDeleted().IsValid())
				}
			}
		})
	}
}
