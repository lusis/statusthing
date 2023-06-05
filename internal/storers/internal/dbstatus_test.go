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

type statusTestcase struct {
	dbstatus *DbStatus
	pbstatus *statusthingv1.Status
	errtext  string
	err      error
}

func TestStatusFromProto(t *testing.T) {
	t.Parallel()

	t.Run("nil-check", func(t *testing.T) {
		s, serr := DbStatusFromProto(nil)
		require.ErrorIs(t, serr, serrors.ErrNilVal)
		require.Nil(t, s)
	})

	// t.Name inside this map refers to the name of the parent test not the iteration.
	// It's intentional since we're not worried about conflict here
	testcases := map[string]statusTestcase{
		"happy-path": {},
		"missing-id": {
			pbstatus: &statusthingv1.Status{Id: ""},
			err:      serrors.ErrEmptyString,
			errtext:  "id",
		},
		"missing-name": {
			pbstatus: &statusthingv1.Status{Id: t.Name()},
			err:      serrors.ErrEmptyString,
			errtext:  "name",
		},
		"missing-kind": {
			pbstatus: &statusthingv1.Status{Id: t.Name(), Name: t.Name(), Timestamps: testutils.MakeTimestamps(false)},
			err:      serrors.ErrEmptyEnum,
			errtext:  "kind",
		},
		"missing-timestamps": {
			pbstatus: &statusthingv1.Status{Id: t.Name(), Name: t.Name(), Kind: statusthingv1.StatusKind_STATUS_KIND_UNAVAILABLE},
			err:      serrors.ErrMissingTimestamp,
			errtext:  "timestamps",
		},
		"missing-created": {
			pbstatus: &statusthingv1.Status{Id: t.Name(), Name: t.Name(), Kind: statusthingv1.StatusKind_STATUS_KIND_UNAVAILABLE, Timestamps: &statusthingv1.Timestamps{
				Updated: timestamppb.Now(),
			}},
			err:     serrors.ErrMissingTimestamp,
			errtext: "created",
		},
		"missing-updated": {
			pbstatus: &statusthingv1.Status{Id: t.Name(), Name: t.Name(), Kind: statusthingv1.StatusKind_STATUS_KIND_UNAVAILABLE, Timestamps: &statusthingv1.Timestamps{
				Created: timestamppb.Now(),
			}},
			err:     serrors.ErrMissingTimestamp,
			errtext: "updated",
		},
	}

	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			pbstatus := tc.pbstatus
			if pbstatus == nil {
				pbstatus = testutils.MakeStatus(t.Name())
				pbstatus.Color = t.Name()
				pbstatus.Description = t.Name()
				pbstatus.Timestamps = testutils.MakeTimestamps(true)
			}
			s, serr := DbStatusFromProto(pbstatus)
			if tc.err != nil {
				require.ErrorIs(t, serr, tc.err)
				require.Nil(t, s)
				if validation.ValidString(tc.errtext) {
					require.ErrorContains(t, serr, tc.errtext)
				}
			} else {
				require.NoError(t, serr)
				require.NotNil(t, s)
				require.Equal(t, pbstatus.GetId(), s.ID)
				require.Equal(t, pbstatus.GetName(), s.Name)
				require.Equal(t, pbstatus.GetKind().String(), *s.Kind)
				require.Equal(t, pbstatus.GetDescription(), *s.Description)
				require.Equal(t, pbstatus.GetColor(), *s.Color)
				require.NotZero(t, s.Created)
				require.NotZero(t, s.Updated)
				if tc.pbstatus.GetTimestamps().GetDeleted().IsValid() {
					require.NotNil(t, s.Deleted)
				}
			}
		})
	}
}

func TestStatusToProto(t *testing.T) {
	t.Parallel()
	// t.Name inside this map refers to the name of the parent test not the iteration.
	// It's intentional since we're not worried about conflict here
	testcases := map[string]statusTestcase{
		"happy-path": {},
		"missing-id": {
			dbstatus: &DbStatus{DbCommon: &DbCommon{ID: ""}},
			err:      serrors.ErrInvalidData,
			errtext:  "id",
		},
		"missing-name": {
			dbstatus: &DbStatus{DbCommon: &DbCommon{ID: t.Name()}},
			err:      serrors.ErrInvalidData,
			errtext:  "name",
		},
		"missing-kind": {
			dbstatus: &DbStatus{DbCommon: &DbCommon{ID: t.Name(), Name: t.Name()}},
			err:      serrors.ErrInvalidData,
			errtext:  "kind",
		},
		"missing-created": {
			dbstatus: &DbStatus{
				DbCommon: &DbCommon{
					ID:           t.Name(),
					Name:         t.Name(),
					DbTimestamps: &DbTimestamps{Updated: uint64(storers.TsToInt64(timestamppb.Now()))},
				},
				Kind: storers.StringPtr(statusthingv1.StatusKind_STATUS_KIND_UNAVAILABLE.String()),
			},
			err:     serrors.ErrInvalidData,
			errtext: "created",
		},
		"missing-updated": {
			dbstatus: &DbStatus{
				DbCommon: &DbCommon{
					ID:           t.Name(),
					Name:         t.Name(),
					DbTimestamps: &DbTimestamps{Created: uint64(storers.TsToInt64(timestamppb.Now()))},
				}, Kind: storers.StringPtr(statusthingv1.StatusKind_STATUS_KIND_UNAVAILABLE.String()),
			},
			err:     serrors.ErrInvalidData,
			errtext: "updated",
		},
	}

	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			dbstatus := tc.dbstatus
			if dbstatus == nil {
				dbstatus = &DbStatus{
					DbCommon: &DbCommon{
						ID:          t.Name(),
						Name:        t.Name(),
						Description: storers.StringPtr(t.Name()),
						DbTimestamps: &DbTimestamps{
							Created: storers.TsToUInt64(timestamppb.Now()),
							Updated: storers.TsToUInt64(timestamppb.Now()),
							Deleted: storers.TsToUInt64Ptr(timestamppb.Now())},
					},
					Kind:  storers.StringPtr(statusthingv1.StatusKind_STATUS_KIND_AVAILABLE.String()),
					Color: storers.StringPtr(t.Name()),
				}
			}
			pb, serr := dbstatus.ToProto()

			if tc.err != nil {
				require.ErrorIs(t, serr, tc.err)
				require.Nil(t, pb)
				if validation.ValidString(tc.errtext) {
					require.ErrorContains(t, serr, tc.errtext)
				}
			} else {
				s, ok := pb.(*statusthingv1.Status)
				require.True(t, ok)
				require.NoError(t, serr)
				require.NotNil(t, s)
				require.Equal(t, dbstatus.ID, s.GetId())
				require.Equal(t, dbstatus.Name, s.GetName())
				require.Equal(t, *dbstatus.Kind, s.GetKind().String())
				require.Equal(t, *dbstatus.Description, s.GetDescription())
				require.Equal(t, *dbstatus.Color, s.GetColor())
				require.True(t, s.GetTimestamps().GetCreated().IsValid())
				require.True(t, s.GetTimestamps().GetUpdated().IsValid())
				if dbstatus.Deleted != nil {
					require.True(t, s.GetTimestamps().GetDeleted().IsValid())
				}
			}
		})
	}
}
