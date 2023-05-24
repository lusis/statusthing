package unimplemented

import (
	"context"

	v1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/errors"
	"github.com/lusis/statusthing/internal/filters"
)

// ItemStore ...
type ItemStore struct{}

// StoreItem stores the provided [statusthingv1.Item]
func (is *ItemStore) StoreItem(ctx context.Context, item *v1.Item) (*v1.Item, error) { // nolint: revive
	return nil, errors.ErrNotImplemented
}

// GetItem gets a [statusthingv1.Item] by its id
func (is *ItemStore) GetItem(ctx context.Context, itemID string) (*v1.Item, error) { // nolint: revive
	return nil, errors.ErrNotImplemented
}

// FindItems returns all known [statusthingv1.Item] optionally filtered by the provided [filters.FilterOption]
func (is *ItemStore) FindItems(ctx context.Context, opts ...filters.FilterOption) ([]*v1.Item, error) { // nolint: revive
	return nil, errors.ErrNotImplemented
}

// UpdateItem updates the [statusthingv1.Item] by its id with the provided [filters.FilterOption]
func (is *ItemStore) UpdateItem(ctx context.Context, itemID string, opts ...filters.FilterOption) error { // nolint: revive
	return errors.ErrNotImplemented
}

// DeleteItem deletes the [statusthingv1.Item] by its id
func (is *ItemStore) DeleteItem(ctx context.Context, itemID string) error { // nolint: revive
	return errors.ErrNotImplemented
}
