package unimplemented

import (
	"context"

	v1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/errors"
	"github.com/lusis/statusthing/internal/filters"
)

// StatusStore stores [statusthingv1.Status]
type StatusStore struct{}

// StoreStatus stores the provided [statusthingv1.Status]
func (ss *StatusStore) StoreStatus(ctx context.Context, status *v1.Status) (*v1.Status, error) { // nolint: revive
	return nil, errors.ErrNotImplemented
}

// GetStatus gets a [statusthingv1.CustomStatus] by its unique id
func (ss *StatusStore) GetStatus(ctx context.Context, statusID string) (*v1.Status, error) { // nolint: revive
	return nil, errors.ErrNotImplemented
}

// FindStatus returns all know [statusthingv1.Status] optionally filtered by the provided [filters.FilterOption]
func (ss *StatusStore) FindStatus(ctx context.Context, opts ...filters.FilterOption) ([]*v1.Status, error) { // nolint: revive
	return nil, errors.ErrNotImplemented
}

// UpdateStatus updates the [statusthingv1.Status] by id with the provided [filters.FilterOption]
func (ss *StatusStore) UpdateStatus(ctx context.Context, statusID string, opts ...filters.FilterOption) error { // nolint:revive
	return errors.ErrNotImplemented
}

// DeleteStatus deletes a [statusthingv1.Status] by its id
func (css *StatusStore) DeleteStatus(ctx context.Context, statusID string) error { // nolint:revive
	return errors.ErrNotImplemented
}
