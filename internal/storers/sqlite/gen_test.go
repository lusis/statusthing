package sqlite

import (
	"context"
	"fmt"
	"testing"

	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	_ "github.com/lusis/statusthing/internal/storers/sqlite/driver" // sql driver

	"github.com/lusis/statusthing/internal/filters"
	"github.com/lusis/statusthing/internal/storers"
	"github.com/lusis/statusthing/internal/testutils"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/require"
)

func TestGenData(t *testing.T) {
	t.Skip()
	ctx := context.TODO()
	db, dberr := makeTestdb(t, ":memory:")
	require.NoError(t, dberr)
	require.NotNil(t, db)
	store, _ := New(db)

	defaultIDs := []string{}
	// Create default statuses
	for _, d := range storers.DefaultStatuses {
		d.Timestamps = testutils.MakeTimestamps(false)
		res, err := store.StoreStatus(ctx, d)
		require.NoError(t, err)
		require.NotNil(t, res)
		defaultIDs = append(defaultIDs, res.GetId())
	}
	// items
	for i := 0; i <= 1000; i++ {
		// add some items with no status
		nires, nierr := store.StoreItem(ctx, &statusthingv1.Item{
			Id:         ksuid.New().String(),
			Name:       fmt.Sprintf("no-status-item-name-%d", i),
			Timestamps: testutils.MakeTimestamps(false),
		})
		require.NoError(t, nierr)
		require.NotNil(t, nires)

		for idi, id := range defaultIDs {
			sires, sireserr := store.StoreItem(ctx, &statusthingv1.Item{
				Id:         ksuid.New().String(),
				Name:       fmt.Sprintf("statusid-make-status-item-name-%d-%d", idi, i),
				Timestamps: testutils.MakeTimestamps(false),
			})
			require.NoError(t, sireserr)
			require.NotNil(t, sires)

			uerr := store.UpdateItem(ctx, sires.GetId(), filters.WithStatusID(id))
			require.NoError(t, uerr)
		}

	}
}
