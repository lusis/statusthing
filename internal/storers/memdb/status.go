package memdb

import (
	"context"
	"fmt"
	"strings"

	hcmemdb "github.com/hashicorp/go-memdb"

	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	v1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/filters"
	"github.com/lusis/statusthing/internal/serrors"
)

const (
	statusTableName    = "statuses"
	statusPrimaryIndex = "id"
)

var statusSchema = &hcmemdb.TableSchema{
	Name: statusTableName,
	Indexes: map[string]*hcmemdb.IndexSchema{
		"id": {
			Name:    statusPrimaryIndex,
			Unique:  true,
			Indexer: &hcmemdb.StringFieldIndex{Field: "ID"},
		},
		"name": {
			Name:    "name",
			Unique:  true,
			Indexer: &hcmemdb.StringFieldIndex{Field: "Name"},
		},
		"description": {
			Name:         "description",
			AllowMissing: true,
			Indexer:      &hcmemdb.StringFieldIndex{Field: "Description"},
		},
		"color": {
			Name:         "color",
			AllowMissing: true,
			Indexer:      &hcmemdb.StringFieldIndex{Field: "Color"},
		},
		"kind": {
			Name:    "kind",
			Indexer: &hcmemdb.StringFieldIndex{Field: "Kind"},
		},
		"created": {
			Name:    "created",
			Indexer: &hcmemdb.IntFieldIndex{Field: "Created"},
		},
		"updated": {
			Name:    "updated",
			Indexer: &hcmemdb.IntFieldIndex{Field: "Updated"},
		},
		"deleted": {
			Name:         "deleted",
			AllowMissing: true,
			Indexer:      &hcmemdb.IntFieldIndex{Field: "Deleted"},
		},
	},
}

type dbStatus struct {
	ID          string
	Name        string
	Description string
	Color       string
	Kind        string
	Created     int
	Updated     int
	Deleted     int
}

func (ms *dbStatus) toProto() (*statusthingv1.Status, error) {
	if ms.ID == "" {
		return nil, fmt.Errorf("id: %w", serrors.ErrInvalidData)
	}
	if ms.Name == "" {
		return nil, fmt.Errorf("name: %w", serrors.ErrInvalidData)
	}

	if ms.Kind == "" || ms.Kind == statusthingv1.StatusKind_STATUS_KIND_UNKNOWN.String() {
		return nil, fmt.Errorf("kind: %w", serrors.ErrInvalidData)
	}

	created, updated, deleted := intToTs(ms.Created), intToTs(ms.Updated), intToTs(ms.Deleted)
	if !created.IsValid() {
		return nil, fmt.Errorf("created: %w", serrors.ErrInvalidData)
	}
	if !updated.IsValid() {
		return nil, fmt.Errorf("updated: %w", serrors.ErrInvalidData)
	}
	status := &statusthingv1.Status{
		Id:          ms.ID,
		Name:        ms.Name,
		Description: ms.Description,
		Color:       ms.Color,
		Kind:        v1.StatusKind(statusthingv1.StatusKind_value[ms.Kind]),
		Timestamps: &v1.Timestamps{
			Created: created,
			Updated: updated,
		},
	}

	if deleted.IsValid() {
		status.Timestamps.Deleted = deleted
	}
	return status, nil
}
func statusFromProto(ps *statusthingv1.Status) (*dbStatus, error) {
	if ps == nil {
		return nil, serrors.ErrNilVal
	}

	if ps.GetId() == "" {
		return nil, serrors.ErrEmptyString
	}
	if ps.GetName() == "" {
		return nil, serrors.ErrEmptyString
	}
	if ps.GetTimestamps() == nil {
		return nil, fmt.Errorf("timestamps: %w", serrors.ErrNilVal)
	}
	if !ps.GetTimestamps().GetCreated().IsValid() {
		return nil, fmt.Errorf("created: %w", serrors.ErrInvalidData)
	}
	if !ps.GetTimestamps().GetUpdated().IsValid() {
		return nil, fmt.Errorf("updated: %w", serrors.ErrInvalidData)
	}
	if ps.GetKind() == statusthingv1.StatusKind_STATUS_KIND_UNKNOWN {
		return nil, fmt.Errorf("kind: %w", serrors.ErrEmptyEnum)
	}
	ms := &dbStatus{
		ID:      ps.GetId(),
		Name:    ps.GetName(),
		Created: tsToInt(ps.GetTimestamps().GetCreated()),
		Updated: tsToInt(ps.GetTimestamps().GetUpdated()),
		Kind:    ps.GetKind().String(),
	}
	if ps.GetTimestamps().GetDeleted().IsValid() {
		ms.Deleted = tsToInt(ps.GetTimestamps().GetDeleted())
	}
	if ps.GetDescription() != "" {
		ms.Description = ps.GetDescription()
	}
	if ps.GetColor() != "" {
		ms.Color = ps.GetColor()
	}
	return ms, nil
}

// StoreStatus stores the provided [statusthingv1.Status]
func (s *StatusThingStore) StoreStatus(ctx context.Context, status *v1.Status) (*v1.Status, error) {
	txn := s.db.Txn(true)
	ms, err := statusFromProto(status)
	if err != nil {
		return nil, err
	}
	if err := txn.Insert(statusTableName, ms); err != nil {
		txn.Abort()
		return nil, err
	}
	txn.Commit()
	return s.GetStatus(ctx, status.GetId())
}

// GetStatus gets a [statusthingv1.Status] by its unique id
func (s *StatusThingStore) GetStatus(ctx context.Context, statusID string) (*v1.Status, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()
	return getStatusWithTxn(ctx, txn, statusID)
}

func getStatusWithTxn(ctx context.Context, txn *hcmemdb.Txn, statusID string) (*v1.Status, error) {
	res, err := firstWithTxn(ctx, txn, statusTableName, statusPrimaryIndex, statusID)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, serrors.ErrNotFound
	}
	status, ok := res.(*dbStatus)
	if !ok {
		return nil, serrors.ErrInvalidData
	}
	return status.toProto()
}

// FindStatus returns all know [statusthingv1.Status] optionally filtered by the provided [filters.FilterOption]
func (s *StatusThingStore) FindStatus(ctx context.Context, opts ...filters.FilterOption) ([]*v1.Status, error) {
	f, err := filters.New(opts...)
	if err != nil {
		return nil, err
	}

	txn := s.db.Txn(false)
	defer txn.Abort()
	res, err := getWithTxn(ctx, txn, statusTableName, statusPrimaryIndex)
	if err != nil {
		return nil, err
	}
	statuses := []*v1.Status{}
	for obj := res.Next(); obj != nil; obj = res.Next() {
		statusEntry, ok := obj.(*dbStatus)
		if !ok {
			return nil, serrors.ErrInvalidData
		}
		// with memdb we filter as we process each results
		if len(f.StatusKinds()) != 0 {
			hasKind := false
			for _, k := range f.StatusKinds() {
				if k.String() == statusEntry.Kind {
					hasKind = true
				}
			}
			if !hasKind {
				continue
			}
		}
		pb, err := statusEntry.toProto()
		if err != nil {
			return nil, err
		}

		statuses = append(statuses, pb)
	}
	return statuses, nil
}

// UpdateStatus updates the [statusthingv1.Status] by id with the provided [filters.FilterOption]
// supported filters:
// - [filters.WithColor]: to change the color
// - [filters.WithDescription]: to change the description
// - [filters.WithStatusKind]: to change the kind of status
// - [filters.WithName]: to change the name
func (s *StatusThingStore) UpdateStatus(ctx context.Context, statusID string, opts ...filters.FilterOption) error {
	if len(opts) == 0 {
		return serrors.ErrAtLeastOne
	}
	if strings.TrimSpace(statusID) == "" {
		return serrors.ErrEmptyString
	}
	f, err := filters.New(opts...)
	if err != nil {
		return err
	}
	existing, err := s.GetStatus(ctx, statusID)
	if err != nil {
		return err
	}
	if f.Color() != "" {
		existing.Color = f.Color()
	}
	if f.Description() != "" {
		existing.Description = f.Description()
	}
	if f.StatusKind() != statusthingv1.StatusKind_STATUS_KIND_UNKNOWN {
		existing.Kind = f.StatusKind()
	}
	if f.Name() != "" {
		existing.Name = f.Name()
	}
	// updates in memdb are just inserts to the same id so we can just call our own store func here
	if _, err := s.StoreStatus(ctx, existing); err != nil {
		return err
	}

	return nil
}

// DeleteStatus deletes a [statusthingv1.Status] by its id
func (s *StatusThingStore) DeleteStatus(ctx context.Context, statusID string) error {
	txn := s.db.Txn(true)
	if err := deleteWithTxn(ctx, txn, statusTableName, &dbStatus{ID: statusID}); err != nil {
		txn.Abort()
		return err
	}
	txn.Commit()
	return nil
}
