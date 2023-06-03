package internal

import (
	"context"

	"google.golang.org/protobuf/proto"

	"github.com/lusis/statusthing/internal/filters"
)

// protoStorer is something that can store protobuf messages
// Not sold on this yet - just needed to dump it
type protoStorer interface {
	Store(ctx context.Context, pb proto.Message) (proto.Message, error)
	Get(ctx context.Context, idval string) (proto.Message, error)
	Find(ctx context.Context, opts ...filters.FilterOption) ([]proto.Message, error)
	Update(ctx context.Context, idval string, opts ...filters.FilterOption) error
	Delete(ctx context.Context, idval string) error
}
