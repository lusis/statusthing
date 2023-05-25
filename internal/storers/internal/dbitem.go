package internal

import (
	"html"

	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/serrors"
	"github.com/lusis/statusthing/internal/storers"
	"github.com/lusis/statusthing/internal/validation"
)

// DbItem represents a common representation of an [statusthingv1.Item] in a db
type DbItem struct {
	*DbCommon
	StatusID *string `db:"status_id"`
}

// DbItemFromProto creates a [DbItem] from a [statusthingv1.Item]
func DbItemFromProto(pbitem *statusthingv1.Item) (*DbItem, error) {
	if pbitem == nil {
		return nil, serrors.NewError("status", serrors.ErrNilVal)
	}
	id := html.EscapeString(pbitem.GetId())
	name := html.EscapeString(pbitem.GetName())
	desc := html.EscapeString(pbitem.GetDescription())
	statusID := pbitem.GetStatus().GetId()
	created := pbitem.GetTimestamps().GetCreated()
	updated := pbitem.GetTimestamps().GetUpdated()
	deleted := pbitem.GetTimestamps().GetDeleted()

	if !validation.ValidString(id) {
		return nil, serrors.NewError("id", serrors.ErrEmptyString)
	}
	if !validation.ValidString(name) {
		return nil, serrors.NewError("name", serrors.ErrEmptyString)
	}
	if pbitem.GetTimestamps() == nil {
		return nil, serrors.NewError("timestamps", serrors.ErrMissingTimestamp)
	}

	if !created.IsValid() {
		return nil, serrors.NewError("created", serrors.ErrMissingTimestamp)
	}
	if !updated.IsValid() {
		return nil, serrors.NewError("updated", serrors.ErrMissingTimestamp)
	}
	dbs := &DbItem{
		DbCommon: &DbCommon{
			ID:      id,
			Name:    name,
			Created: storers.TsToUInt64(created),
			Updated: storers.TsToUInt64(updated),
		},
	}

	if validation.ValidString(statusID) {
		dbs.StatusID = storers.StringPtr(statusID)
	}

	if validation.ValidString(desc) {
		dbs.Description = storers.StringPtr(desc)
	}

	if deleted.IsValid() {
		dbs.Deleted = storers.TsToUInt64Ptr(deleted)
	}
	return dbs, nil
}

// ToProto converts a [DbItem] to a [statusthingv1.Item]
func (s *DbItem) ToProto() (*statusthingv1.Item, error) {
	res := &statusthingv1.Item{
		Timestamps: &statusthingv1.Timestamps{},
	}
	if !validation.ValidString(s.ID) {
		return nil, serrors.NewError("id", serrors.ErrInvalidData)
	}
	if !validation.ValidString(s.Name) {
		return nil, serrors.NewError("name", serrors.ErrInvalidData)
	}

	// id/name
	res.Id = html.UnescapeString(s.ID)
	res.Name = html.UnescapeString(s.Name)

	//desc
	if s.Description != nil {
		res.Description = html.UnescapeString(*s.Description)
	}

	// timestamps
	pbcreated := storers.Int64ToTs(int64(s.Created))
	pbupdated := storers.Int64ToTs(int64(s.Updated))

	if pbcreated == nil {
		return nil, serrors.NewError("created", serrors.ErrInvalidData)
	}
	if pbupdated == nil {
		return nil, serrors.NewError("updated", serrors.ErrInvalidData)
	}
	res.Timestamps.Created = pbcreated
	res.Timestamps.Updated = pbupdated

	if s.Deleted != nil {
		res.Timestamps.Deleted = storers.Int64ToTs(int64(*s.Deleted))
	}
	// status will be populated outside of here for now
	//
	return res, nil
}
