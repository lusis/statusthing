package memdb

import (
	"testing"

	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/storers"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func makeTestTsNow() *statusthingv1.Timestamps {
	now := timestamppb.Now()
	return &statusthingv1.Timestamps{
		Created: now,
		Updated: now,
	}
}
func TestImplements(t *testing.T) {
	require.Implements(t, (*storers.StatusThingStorer)(nil), new(StatusThingStore), "unimplemented status thing store should sastify interface")
}

func TestNew(t *testing.T) {
	store, err := New()
	require.NoError(t, err)
	require.NotNil(t, store)
}
