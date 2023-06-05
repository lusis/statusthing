package internal

import (
	"html"

	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/serrors"
	"github.com/lusis/statusthing/internal/storers"
	"github.com/lusis/statusthing/internal/validation"

	"google.golang.org/protobuf/proto"
)

// DbStatus represents a common representation of a [statusthingv1.Status] in a db
type DbStatus struct {
	*DbCommon
	Color *string `db:"color"`
	Kind  *string `db:"kind"`
}

// DbStatusFromProto creates a [DbStatus] from a [statusthingv1.Status]
func DbStatusFromProto(pbstatus *statusthingv1.Status) (*DbStatus, error) {
	if pbstatus == nil {
		return nil, serrors.NewError("status", serrors.ErrNilVal)
	}
	dbCommon, err := MakeDbCommon(pbstatus.GetId(), pbstatus.GetName(), pbstatus.GetDescription(), pbstatus.GetTimestamps())
	if err != nil {
		return nil, err
	}

	color := html.EscapeString(pbstatus.GetColor())
	kind := pbstatus.GetKind()

	if kind == statusthingv1.StatusKind_STATUS_KIND_UNKNOWN {
		return nil, serrors.NewError("kind", serrors.ErrEmptyEnum)
	}

	dbs := &DbStatus{
		DbCommon: dbCommon,
		Kind:     storers.StringPtr(kind.String()),
	}
	if validation.ValidString(color) {
		dbs.Color = storers.StringPtr(color)
	}
	return dbs, nil
}

// ToProto converts a [DbStatus] to a [statusthingv1.Status]
func (s *DbStatus) ToProto() (proto.Message, error) {
	res := &statusthingv1.Status{
		Timestamps: &statusthingv1.Timestamps{},
	}
	if !validation.ValidString(s.ID) {
		return nil, serrors.NewError("id", serrors.ErrInvalidData)
	}
	if !validation.ValidString(s.Name) {
		return nil, serrors.NewError("name", serrors.ErrInvalidData)
	}
	if s.Kind == nil || s.Kind == storers.StringPtr(statusthingv1.StatusKind_STATUS_KIND_UNKNOWN.String()) {
		return nil, serrors.NewError("kind", serrors.ErrInvalidData)
	}

	// id/name
	res.Id = html.UnescapeString(s.ID)
	res.Name = html.UnescapeString(s.Name)

	// kind
	res.Kind = statusthingv1.StatusKind(statusthingv1.StatusKind_value[*s.Kind])

	//desc/color
	if s.Description != nil {
		res.Description = html.UnescapeString(*s.Description)
	}
	if s.Color != nil {
		res.Color = html.UnescapeString(*s.Color)
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
	return res, nil
}
