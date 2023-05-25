package storers

import (
	"context"

	v1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/filters"
	"github.com/segmentio/ksuid"
)

// TODO: consider migrating this to generic if any non-type specific store functions are needed

// StatusThingStorer stores [sttypes.StatusThing]
// This is the only store we work directly with
// The embedding is to allow more flexible storage options
type StatusThingStorer interface {
	NoteStorer
	StatusStorer
	ItemStorer
}

// ItemStorer storers [statusthingv1.Item]
type ItemStorer interface {
	// StoreItem stores the provided [statusthingv1.Item]
	StoreItem(ctx context.Context, item *v1.Item) (*v1.Item, error)
	// GetItem gets a [statusthingv1.Item] by its id
	GetItem(ctx context.Context, itemID string) (*v1.Item, error)
	// FindItems returns all known [statusthingv1.Item] optionally filtered by the provided [filters.FilterOption]
	FindItems(ctx context.Context, opts ...filters.FilterOption) ([]*v1.Item, error)
	// UpdateItem updates the [statusthingv1.Item] by its id with the provided [filters.FilterOption]
	UpdateItem(ctx context.Context, itemID string, opts ...filters.FilterOption) error
	// DeleteItem deletes the [statusthingv1.Item] by its id
	DeleteItem(ctx context.Context, itemID string) error
}

// NoteStorer stores [statusthingv1.Note]
type NoteStorer interface {
	// StoreNote stores the provided [statusthingv1.Note] associated with the provided [statusthingv1.StatusThing] by its id
	StoreNote(ctx context.Context, note *v1.Note, statusThingID string) (*v1.Note, error)
	// GetNote gets a [statusthingv1.Note] by its id
	GetNote(ctx context.Context, noteID string) (*v1.Note, error)
	// FindNotes gets all known [statusthingv1.Note] for the provided item id
	FindNotes(ctx context.Context, itemID string, opts ...filters.FilterOption) ([]*v1.Note, error)
	// UpdateNote updates the [statusthingv1.Note] with the provided [filters.FilterOption]
	UpdateNote(ctx context.Context, noteID string, opts ...filters.FilterOption) error
	// DeleteNote deletes a [statusthingv1.Note] by its id
	DeleteNote(ctx context.Context, noteID string) error
}

// StatusStorer stores [statusthingv1.Status]
type StatusStorer interface {
	// StoreStatus stores the provided [statusthingv1.Status]
	StoreStatus(ctx context.Context, status *v1.Status) (*v1.Status, error)
	// GetStatus gets a [statusthingv1.Status] by its unique id
	GetStatus(ctx context.Context, statusID string) (*v1.Status, error)
	// FindStatus returns all know [statusthingv1.Status] optionally filtered by the provided [filters.FilterOption]
	FindStatus(ctx context.Context, opts ...filters.FilterOption) ([]*v1.Status, error)
	// UpdateStatus updates the [statusthingv1.Status] by id with the provided [filters.FilterOption]
	UpdateStatus(ctx context.Context, statusID string, opts ...filters.FilterOption) error
	// DeleteStatus deletes a [statusthingv1.Status] by its id
	DeleteStatus(ctx context.Context, statusID string) error
}

const (
	green  = "#5DFC0A"
	red    = "#FF0000"
	yellow = "#EEEB8D"
)

// DefaultStatuses are safe default values that can be populated
var DefaultStatuses = []*v1.Status{
	// up status
	{Id: func() string { return ksuid.New().String() }(), Name: "UP", Kind: v1.StatusKind_STATUS_KIND_UP, Color: green},
	{Id: func() string { return ksuid.New().String() }(), Name: "DOWN", Kind: v1.StatusKind_STATUS_KIND_DOWN, Color: red},
	{Id: func() string { return ksuid.New().String() }(), Name: "WARNING", Kind: v1.StatusKind_STATUS_KIND_WARNING, Color: yellow},
	{Id: func() string { return ksuid.New().String() }(), Name: "CREATED", Kind: v1.StatusKind_STATUS_KIND_CREATED, Color: green},
	{Id: func() string { return ksuid.New().String() }(), Name: "OFFLINE", Kind: v1.StatusKind_STATUS_KIND_OFFLINE, Color: red},
	{Id: func() string { return ksuid.New().String() }(), Name: "ONLINE", Kind: v1.StatusKind_STATUS_KIND_ONLINE, Color: green},
	{Id: func() string { return ksuid.New().String() }(), Name: "OBSERVING", Kind: v1.StatusKind_STATUS_KIND_OBSERVING, Color: yellow},
	{Id: func() string { return ksuid.New().String() }(), Name: "INVESTIGATING", Kind: v1.StatusKind_STATUS_KIND_INVESTIGATING, Color: yellow},
}
