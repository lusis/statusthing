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

type itemTestCase struct {
	dbitem  *DbItem
	pbitem  *statusthingv1.Item
	errtext string
	err     error
}

func TestItemFromProto(t *testing.T) {
	t.Parallel()

	t.Run("nil-check", func(t *testing.T) {
		s, serr := DbItemFromProto(nil)
		require.ErrorIs(t, serr, serrors.ErrNilVal)
		require.Nil(t, s)
	})

	// t.Name inside this map refers to the name of the parent test not the iteration.
	// It's intentional since we're not worried about conflict here
	testcases := map[string]itemTestCase{
		"happy-path": {},
		"missing-id": {
			pbitem:  &statusthingv1.Item{Id: ""},
			err:     serrors.ErrEmptyString,
			errtext: "id",
		},
		"missing-name": {
			pbitem:  &statusthingv1.Item{Id: t.Name()},
			err:     serrors.ErrEmptyString,
			errtext: "name",
		},
		"missing-timestamps": {
			pbitem:  &statusthingv1.Item{Id: t.Name(), Name: t.Name()},
			err:     serrors.ErrMissingTimestamp,
			errtext: "timestamps",
		},
		"missing-created": {
			pbitem: &statusthingv1.Item{Id: t.Name(), Name: t.Name(), Timestamps: &statusthingv1.Timestamps{
				Updated: timestamppb.Now(),
			}},
			err:     serrors.ErrMissingTimestamp,
			errtext: "created",
		},
		"missing-updated": {
			pbitem: &statusthingv1.Item{Id: t.Name(), Name: t.Name(), Timestamps: &statusthingv1.Timestamps{
				Created: timestamppb.Now(),
			}},
			err:     serrors.ErrMissingTimestamp,
			errtext: "updated",
		},
	}

	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			pb := tc.pbitem
			if pb == nil {
				pb = testutils.MakeItem(t.Name())
				pb.Description = t.Name()
				pb.Status = testutils.MakeStatus(t.Name())
				pb.Timestamps = testutils.MakeTimestamps(true)
			}
			s, serr := DbItemFromProto(pb)
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
				require.Equal(t, pb.GetName(), s.Name)
				require.Equal(t, pb.GetDescription(), *s.Description)
				require.NotZero(t, s.Created)
				require.NotZero(t, s.Updated)
				if pb.GetStatus() != nil {
					require.Equal(t, pb.GetStatus().GetId(), *s.StatusID)
				}
				if tc.pbitem.GetTimestamps().GetDeleted().IsValid() {
					require.NotNil(t, s.Deleted)
				}
			}
		})
	}
}

func TestItemToProto(t *testing.T) {
	t.Parallel()
	// t.Name inside this map refers to the name of the parent test not the iteration.
	// It's intentional since we're not worried about conflict here
	testcases := map[string]itemTestCase{
		"happy-path": {},
		"missing-id": {
			dbitem:  &DbItem{DbCommon: &DbCommon{ID: ""}},
			err:     serrors.ErrInvalidData,
			errtext: "id",
		},
		"missing-name": {
			dbitem:  &DbItem{DbCommon: &DbCommon{ID: t.Name()}},
			err:     serrors.ErrInvalidData,
			errtext: "name",
		},
		"missing-created": {
			dbitem: &DbItem{
				DbCommon: &DbCommon{
					ID:           t.Name(),
					Name:         t.Name(),
					DbTimestamps: &DbTimestamps{Updated: uint64(storers.TsToInt64(timestamppb.Now()))},
				},
			},
			err:     serrors.ErrInvalidData,
			errtext: "created",
		},
		"missing-updated": {
			dbitem: &DbItem{
				DbCommon: &DbCommon{
					ID:           t.Name(),
					Name:         t.Name(),
					DbTimestamps: &DbTimestamps{Created: uint64(storers.TsToInt64(timestamppb.Now()))},
				},
			},
			err:     serrors.ErrInvalidData,
			errtext: "updated",
		},
	}

	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			dbitem := tc.dbitem
			if dbitem == nil {
				dbitem = &DbItem{
					DbCommon: &DbCommon{
						ID:          t.Name(),
						Name:        t.Name(),
						Description: storers.StringPtr(t.Name()),
						DbTimestamps: &DbTimestamps{
							Created: storers.TsToUInt64(timestamppb.Now()),
							Updated: storers.TsToUInt64(timestamppb.Now()),
							Deleted: storers.TsToUInt64Ptr(timestamppb.Now())},
					},
					StatusID: storers.StringPtr(t.Name()),
				}
			}
			s, serr := dbitem.ToProto()
			if tc.err != nil {
				require.ErrorIs(t, serr, tc.err)
				require.Nil(t, s)
				if validation.ValidString(tc.errtext) {
					require.ErrorContains(t, serr, tc.errtext)
				}
			} else {
				require.NoError(t, serr)
				require.NotNil(t, s)
				require.Equal(t, dbitem.ID, s.GetId())
				require.Equal(t, dbitem.Name, s.GetName())
				require.Equal(t, *dbitem.Description, s.GetDescription())
				require.True(t, s.GetTimestamps().GetCreated().IsValid())
				require.True(t, s.GetTimestamps().GetUpdated().IsValid())
				if validation.ValidString(*dbitem.StatusID) {
					require.Equal(t, *dbitem.StatusID, s.GetStatus().GetId())
				}
				if dbitem.Deleted != nil {
					require.True(t, s.GetTimestamps().GetDeleted().IsValid())
				}
			}
		})
	}
}
