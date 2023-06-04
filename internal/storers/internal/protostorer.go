package internal

import (
	"context"

	"google.golang.org/protobuf/proto"

	"github.com/lusis/statusthing/internal/filters"
)

// ProtoStorer is something that can store protobuf messages
// Not sold on this yet - just needed to dump it
type ProtoStorer interface {
	Insert(ctx context.Context, pb proto.Message) (proto.Message, error)
	Select(ctx context.Context, idval string) (proto.Message, error)
	Update(ctx context.Context, idval string, opts ...filters.FilterOption) error
	Delete(ctx context.Context, idval string) error
	Where(ctx context.Context, opts ...filters.FilterOption) ([]proto.Message, error)
}
