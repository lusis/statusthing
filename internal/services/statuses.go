package services

import (
	"context"
	"fmt"
	"strings"

	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/filters"
	"github.com/lusis/statusthing/internal/serrors"

	"github.com/segmentio/ksuid"
)

// NewStatus adds a new [statusthingv1.CustomStatus] with the provided name and [statusthingv1.StatusKind]
// supported filters:
// - [filters.WithColor]
// - [filters.WithDescription]
func (sts *StatusThingService) NewStatus(ctx context.Context, statusName string, statusKind statusthingv1.StatusKind, opts ...filters.FilterOption) (*statusthingv1.Status, error) {
	if sts.store == nil {
		return nil, fmt.Errorf("store was nil: %w", serrors.ErrStoreUnavailable)
	}
	if strings.TrimSpace(statusName) == "" {
		return nil, fmt.Errorf("statusName: %w", serrors.ErrEmptyString)
	}
	if statusKind == statusthingv1.StatusKind_STATUS_KIND_UNKNOWN {
		return nil, fmt.Errorf("statusKind: %w", serrors.ErrEmptyEnum)
	}
	f, err := filters.New(opts...)
	if err != nil {
		return nil, err
	}

	id := ksuid.New().String()
	if strings.TrimSpace(f.StatusID()) != "" {
		id = f.StatusID()
	}

	newStatus := &statusthingv1.Status{
		Id:         id,
		Name:       statusName,
		Timestamps: makeTsNow(),
		Kind:       statusKind,
	}
	if strings.TrimSpace(f.Description()) != "" {
		newStatus.Description = f.Description()
	}

	if strings.TrimSpace(f.Color()) != "" {
		newStatus.Color = f.Color()
	}
	return sts.store.StoreStatus(ctx, newStatus)
}

// EditStatus updates a [statusthingv1.tatus] by its id
func (sts *StatusThingService) EditStatus(ctx context.Context, statusID string, opts ...filters.FilterOption) error {
	if sts.store == nil {
		return fmt.Errorf("store was nil: %w", serrors.ErrStoreUnavailable)
	}
	if len(opts) == 0 {
		return fmt.Errorf("opts: %w", serrors.ErrAtLeastOne)
	}
	if strings.TrimSpace(statusID) == "" {
		return fmt.Errorf("statusName: %w", serrors.ErrEmptyString)
	}

	return sts.store.UpdateStatus(ctx, statusID, opts...)
}

// AllStatuses gets all [statusthingv1.Status]
func (sts *StatusThingService) AllStatuses(ctx context.Context, opts ...filters.FilterOption) ([]*statusthingv1.Status, error) {
	return sts.store.FindStatus(ctx, opts...)
}

// DeleteStatus removes a [statusthingv1.Status] by its unique id
func (sts *StatusThingService) DeleteStatus(ctx context.Context, statusID string) error {
	if statusID == "" {
		return fmt.Errorf("statusID: %w", serrors.ErrEmptyString)
	}
	if sts.store == nil {
		return fmt.Errorf("store was nil: %w", serrors.ErrStoreUnavailable)
	}
	return sts.store.DeleteStatus(ctx, statusID)
}

// GetStatus gets a [statusthingv1.Status] by id
func (sts *StatusThingService) GetStatus(ctx context.Context, statusID string) (*statusthingv1.Status, error) {
	if sts.store == nil {
		return nil, fmt.Errorf("store was nil: %w", serrors.ErrStoreUnavailable)
	}
	return sts.store.GetStatus(ctx, statusID)
}
