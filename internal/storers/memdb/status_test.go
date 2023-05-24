package memdb

import (
	"context"
	"testing"

	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/filters"
	"github.com/lusis/statusthing/internal/serrors"
	"github.com/lusis/statusthing/internal/testutils"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/stretchr/testify/require"
)

func TestUpdateStatusRequirements(t *testing.T) {
	mem, _ := New()
	err := mem.UpdateStatus(context.TODO(), "", filters.WithColor("red"))
	require.ErrorIs(t, err, serrors.ErrEmptyString)
	err = mem.UpdateStatus(context.TODO(), "foo")
	require.ErrorIs(t, err, serrors.ErrAtLeastOne)
}
func TestStatusLifecycle(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		err    error
		status *statusthingv1.Status
	}{
		"minimum": {
			err:    nil,
			status: testutils.MakeStatus(t.Name()),
		},
		"invalid-kind": {
			err:    serrors.ErrEmptyString,
			status: &statusthingv1.Status{Name: t.Name(), Kind: statusthingv1.StatusKind_STATUS_KIND_UNKNOWN},
		},
		"all-fields": {
			err: nil,
			status: func() *statusthingv1.Status {
				status := testutils.MakeStatus(t.Name())
				status.Timestamps = testutils.MakeTimestamps(true)
				status.Description = t.Name() + "_status_description"
				status.Color = t.Name() + "_status_color"
				status.Kind = statusthingv1.StatusKind_STATUS_KIND_DOWN
				return status
			}(),
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			store, _ := New()
			res, err := store.StoreStatus(context.TODO(), tc.status)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				require.Nil(t, res)
			} else {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, tc.status.GetId(), res.GetId())
				require.Equal(t, tc.status.GetName(), res.GetName())
				require.Equal(t, tc.status.GetKind(), res.GetKind())
				require.Equal(t, tc.status.GetColor(), res.GetColor())
				require.Equal(t, tc.status.GetDescription(), res.GetDescription())
				require.True(t, testutils.TimestampsEqual(tc.status.Timestamps, res.Timestamps))

				all, err := store.FindStatus(context.TODO())
				require.NoError(t, err)
				require.Len(t, all, 1)
				first := all[0]
				require.Equal(t, tc.status.GetId(), first.GetId())
				require.Equal(t, tc.status.GetName(), first.GetName())
				require.Equal(t, tc.status.GetKind(), first.GetKind())
				require.Equal(t, tc.status.GetColor(), first.GetColor())
				require.Equal(t, tc.status.GetDescription(), first.GetDescription())
				require.True(t, testutils.TimestampsEqual(tc.status.Timestamps, first.Timestamps))

				updateErr := store.UpdateStatus(
					context.TODO(),
					first.GetId(),
					filters.WithColor("new_color"),
					filters.WithName("new_name"),
					filters.WithDescription("new_description"),
					filters.WithStatusKind(statusthingv1.StatusKind_STATUS_KIND_INVESTIGATING),
				)
				require.NoError(t, updateErr)
				checkUp, checkUpErr := store.GetStatus(context.TODO(), first.GetId())
				require.NoError(t, checkUpErr)
				require.Equal(t, tc.status.GetId(), checkUp.GetId())
				require.Equal(t, "new_name", checkUp.GetName())
				require.Equal(t, statusthingv1.StatusKind_STATUS_KIND_INVESTIGATING, checkUp.GetKind())
				require.Equal(t, "new_color", checkUp.GetColor())
				require.Equal(t, "new_description", checkUp.GetDescription())
				require.True(t, testutils.TimestampsEqual(tc.status.Timestamps, checkUp.Timestamps))

				delErr := store.DeleteStatus(context.TODO(), res.Id)
				require.NoError(t, delErr)

				nfErr := store.DeleteStatus(context.TODO(), res.Id)
				require.ErrorIs(t, nfErr, serrors.ErrNotFound)
			}
		})
	}
}

func TestStatusKindFiltering(t *testing.T) {
	t.Parallel()
	store, _ := New()
	statuses := []*statusthingv1.Status{
		{Id: "available", Name: "available-status", Kind: statusthingv1.StatusKind_STATUS_KIND_AVAILABLE},
		{Id: "unavailable", Name: "unavailable-status", Kind: statusthingv1.StatusKind_STATUS_KIND_UNAVAILABLE},
		{Id: "down", Name: "down-status", Kind: statusthingv1.StatusKind_STATUS_KIND_DOWN},
		{Id: "up", Name: "up-status", Kind: statusthingv1.StatusKind_STATUS_KIND_UP},
		{Id: "warning", Name: "warning-status", Kind: statusthingv1.StatusKind_STATUS_KIND_WARNING},
		{Id: "investigating", Name: "investigating-status", Kind: statusthingv1.StatusKind_STATUS_KIND_INVESTIGATING},
		{Id: "observing", Name: "observing-status", Kind: statusthingv1.StatusKind_STATUS_KIND_OBSERVING},
		{Id: "created", Name: "created-status", Kind: statusthingv1.StatusKind_STATUS_KIND_CREATED},
	}
	for _, status := range statuses {
		t.Run(status.GetId(), func(t *testing.T) {
			status.Timestamps = makeTestTsNow()
			res, err := store.StoreStatus(context.TODO(), status)
			require.NoError(t, err)
			require.NotNil(t, res)
			require.Equal(t, status.GetId(), res.GetId())
			require.Equal(t, status.GetName(), res.GetName())
			require.Equal(t, status.GetKind(), res.GetKind())

			all, allerr := store.FindStatus(context.TODO(), filters.WithStatusKinds(status.GetKind()))
			require.NoError(t, allerr)
			require.Len(t, all, 1)
			require.Equal(t, status.GetId(), all[0].GetId())
			require.Equal(t, status.GetName(), all[0].GetName())
			require.Equal(t, status.GetKind(), all[0].GetKind())
		})

	}
}
func TestDbStatusFromProto(t *testing.T) {
	type testcase struct {
		protostatus *statusthingv1.Status
		dbstatus    *dbStatus
		err         error
	}
	testcases := map[string]testcase{
		"nil": {err: serrors.ErrNilVal},
		"happy-path": {
			protostatus: &statusthingv1.Status{
				Id:          t.Name() + "id",
				Name:        t.Name(),
				Kind:        statusthingv1.StatusKind_STATUS_KIND_AVAILABLE,
				Description: t.Name() + "desc",
				Color:       t.Name() + "color",
				Timestamps:  makeTestTsNow(),
			},
			dbstatus: &dbStatus{
				ID:          t.Name() + "id",
				Name:        t.Name(),
				Description: t.Name() + "desc",
				Color:       t.Name() + "color",
				Kind:        "STATUS_KIND_AVAILABLE",
			},
		},
		"invalid-kind": {
			protostatus: &statusthingv1.Status{
				Id:          t.Name() + "id",
				Name:        t.Name(),
				Description: t.Name() + "desc",
				Color:       t.Name() + "color",
				Timestamps:  makeTestTsNow(),
			},
			err: serrors.ErrEmptyEnum,
		},
		"invalid-id": {
			protostatus: &statusthingv1.Status{
				Name:        t.Name(),
				Kind:        statusthingv1.StatusKind_STATUS_KIND_AVAILABLE,
				Description: t.Name() + "desc",
				Color:       t.Name() + "color",
				Timestamps:  makeTestTsNow(),
			},
			err: serrors.ErrEmptyString,
		},
		"invalid-name": {
			protostatus: &statusthingv1.Status{
				Id:          t.Name(),
				Kind:        statusthingv1.StatusKind_STATUS_KIND_AVAILABLE,
				Description: t.Name() + "desc",
				Color:       t.Name() + "color",
				Timestamps:  makeTestTsNow(),
			},
			err: serrors.ErrEmptyString,
		},
		"nil-timestamps": {
			protostatus: &statusthingv1.Status{
				Id:          t.Name(),
				Name:        t.Name(),
				Kind:        statusthingv1.StatusKind_STATUS_KIND_AVAILABLE,
				Description: t.Name() + "desc",
				Color:       t.Name() + "color",
			},
			err: serrors.ErrNilVal,
		},
		"missing-created": {
			protostatus: &statusthingv1.Status{
				Id:          t.Name(),
				Name:        t.Name(),
				Kind:        statusthingv1.StatusKind_STATUS_KIND_AVAILABLE,
				Description: t.Name() + "desc",
				Color:       t.Name() + "color",
				Timestamps: &statusthingv1.Timestamps{
					Updated: timestamppb.Now(),
				},
			},
			err: serrors.ErrInvalidData,
		},
		"missing-updated": {
			protostatus: &statusthingv1.Status{
				Id:          t.Name(),
				Name:        t.Name(),
				Kind:        statusthingv1.StatusKind_STATUS_KIND_AVAILABLE,
				Description: t.Name() + "desc",
				Color:       t.Name() + "color",
				Timestamps: &statusthingv1.Timestamps{
					Created: timestamppb.Now(),
				},
			},
			err: serrors.ErrInvalidData,
		},
	}

	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			res, err := statusFromProto(tc.protostatus)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				require.Nil(t, res)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.dbstatus.Name, tc.protostatus.GetName())
				require.Equal(t, tc.dbstatus.Description, tc.protostatus.GetDescription())
				require.Equal(t, tc.dbstatus.Color, tc.protostatus.GetColor())
				require.Equal(t, tc.dbstatus.Kind, tc.protostatus.GetKind().String())
			}
		})
	}
}

func TestProtoToDbStatus(t *testing.T) {
	type testcase struct {
		protostatus *statusthingv1.Status
		dbstatus    *dbStatus
		err         error
	}
	timestamps := makeTestTsNow()
	testcases := map[string]testcase{
		"happy-path": {
			protostatus: &statusthingv1.Status{
				Id:          t.Name() + "id",
				Name:        t.Name(),
				Kind:        statusthingv1.StatusKind_STATUS_KIND_AVAILABLE,
				Description: t.Name() + "desc",
				Color:       t.Name() + "color",
			},
			dbstatus: &dbStatus{
				ID:          t.Name() + "id",
				Name:        t.Name(),
				Description: t.Name() + "desc",
				Color:       t.Name() + "color",
				Kind:        "STATUS_KIND_AVAILABLE",
				Created:     tsToInt(timestamps.GetCreated()),
				Updated:     tsToInt(timestamps.GetUpdated()),
			},
		},
		"missing-id": {
			dbstatus: &dbStatus{
				Name:        t.Name(),
				Description: t.Name() + "desc",
				Color:       t.Name() + "color",
				Kind:        "STATUS_KIND_AVAILABLE",
				Created:     tsToInt(timestamps.GetCreated()),
				Updated:     tsToInt(timestamps.GetUpdated()),
			},
			err: serrors.ErrInvalidData,
		},
		"missing-name": {
			dbstatus: &dbStatus{
				ID:          t.Name(),
				Description: t.Name() + "desc",
				Color:       t.Name() + "color",
				Kind:        "STATUS_KIND_AVAILABLE",
				Created:     tsToInt(timestamps.GetCreated()),
				Updated:     tsToInt(timestamps.GetUpdated()),
			},
			err: serrors.ErrInvalidData,
		},
		"unknown-kind": {
			dbstatus: &dbStatus{
				ID:          t.Name(),
				Name:        t.Name(),
				Description: t.Name() + "desc",
				Color:       t.Name() + "color",
				Kind:        "STATUS_KIND_UNKNOWN",
				Created:     tsToInt(timestamps.GetCreated()),
				Updated:     tsToInt(timestamps.GetUpdated()),
			},
			err: serrors.ErrInvalidData,
		},
		"missing-kind": {
			dbstatus: &dbStatus{
				ID:          t.Name(),
				Name:        t.Name(),
				Description: t.Name() + "desc",
				Color:       t.Name() + "color",
				Created:     tsToInt(timestamps.GetCreated()),
				Updated:     tsToInt(timestamps.GetUpdated()),
			},
			err: serrors.ErrInvalidData,
		},
		"invalid-created": {
			dbstatus: &dbStatus{
				ID:          t.Name(),
				Name:        t.Name(),
				Description: t.Name() + "desc",
				Color:       t.Name() + "color",
				Kind:        statusthingv1.StatusKind_STATUS_KIND_AVAILABLE.String(),
				Created:     0,
				Updated:     tsToInt(timestamps.GetUpdated()),
			},
			err: serrors.ErrInvalidData,
		},
		"invalid-updated": {
			dbstatus: &dbStatus{
				ID:          t.Name(),
				Name:        t.Name(),
				Description: t.Name() + "desc",
				Color:       t.Name() + "color",
				Kind:        statusthingv1.StatusKind_STATUS_KIND_AVAILABLE.String(),
				Created:     tsToInt(timestamps.GetCreated()),
				Updated:     0,
			},
			err: serrors.ErrInvalidData,
		},
	}

	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			res, err := tc.dbstatus.toProto()
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				require.Nil(t, res)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.dbstatus.Name, tc.protostatus.GetName())
				require.Equal(t, tc.dbstatus.Description, tc.protostatus.GetDescription())
				require.Equal(t, tc.dbstatus.Color, tc.protostatus.GetColor())
				require.Equal(t, tc.dbstatus.Kind, tc.protostatus.GetKind().String())
			}
		})
	}
}
