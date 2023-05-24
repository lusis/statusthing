// nolint: revive
// TODO: remove nolint after finishing implementations
package services

import (
	"context"
	"fmt"
	"sync"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"golang.org/x/exp/slog"

	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/errors"
	"github.com/lusis/statusthing/internal/storers"
)

// TODO: add interface after API is stabilized

// StatusThingService is a domain service for [statusthingv1.StatusThing] and related items
type StatusThingService struct {
	l               *sync.RWMutex
	store           storers.StatusThingStorer
	loadDefaults    bool
	defaultStatuses []*statusthingv1.Status
}

// NewStatusThingService returns a new [StatusThingService]
func NewStatusThingService(store storers.StatusThingStorer, opts ...ServiceOption) (*StatusThingService, error) {
	if store == nil {
		return nil, fmt.Errorf("nil store provided: %w", errors.ErrStoreUnavailable)
	}
	svc := &StatusThingService{
		l:               &sync.RWMutex{},
		store:           store,
		defaultStatuses: make([]*statusthingv1.Status, 0),
	}
	for _, opt := range opts {
		svc.l.Lock()
		if err := opt(svc); err != nil {
			svc.l.Unlock()
			return nil, err
		}
		svc.l.Unlock()
	}
	if svc.loadDefaults {
		if err := svc.CreateDefaultStatuses(); err != nil {
			return nil, err
		}
	}

	return svc, nil
}

// GetCreatedDefaults gets the default status that were created if requested
func (sts *StatusThingService) GetCreatedDefaults() []*statusthingv1.Status {
	sts.l.RLock()
	defer sts.l.RUnlock()
	return sts.defaultStatuses
}

func (sts *StatusThingService) CreateDefaultStatuses() error {
	existing, err := sts.store.FindStatus(context.TODO())
	if err != nil {
		return fmt.Errorf("unable to check existing statuses: %w", err)
	}
	if len(existing) == 0 {
		for _, s := range defaultStatuses {
			cp := proto.Clone(s).(*statusthingv1.Status)
			cp.Timestamps = makeTsNow()
			sts.l.Lock()
			created, err := sts.store.StoreStatus(context.TODO(), cp)
			if err != nil {
				slog.Warn("unable to add default status", "error", err, "status.name", cp.GetName())
			}
			sts.defaultStatuses = append(sts.defaultStatuses, created)
			sts.l.Unlock()
		}
	}
	return nil
}
func makeTsNow() *statusthingv1.Timestamps {
	now := timestamppb.Now()
	return &statusthingv1.Timestamps{
		Created: now,
		Updated: now,
	}
}
