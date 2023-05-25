// Package internal contains internal storer code
package internal

import (
	"html"

	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/serrors"
	"github.com/lusis/statusthing/internal/storers"
	"github.com/lusis/statusthing/internal/validation"
)

// DbCommon is a struct that maps to common fields we use
type DbCommon struct {
	ID          string  `db:"id" goqu:"skipupdate"`
	Name        string  `db:"name"`
	Description *string `db:"description"`
	Created     uint64  `db:"created"`
	Updated     uint64  `db:"updated"`
	Deleted     *uint64 `db:"deleted"`
}

// MakeDbCommon builds the minimum required common fields of all db records
func MakeDbCommon(id, name, desc string, timestamps *statusthingv1.Timestamps) (*DbCommon, error) {
	id = html.EscapeString(id)
	name = html.EscapeString(name)
	desc = html.EscapeString(desc)
	created := timestamps.GetCreated()
	updated := timestamps.GetUpdated()
	deleted := timestamps.GetDeleted()

	if !validation.ValidString(id) {
		return nil, serrors.NewError("id", serrors.ErrEmptyString)
	}
	if !validation.ValidString(name) {
		return nil, serrors.NewError("name", serrors.ErrEmptyString)
	}
	if timestamps == nil {
		return nil, serrors.NewError("timestamps", serrors.ErrMissingTimestamp)
	}
	if !created.IsValid() {
		return nil, serrors.NewError("created", serrors.ErrMissingTimestamp)
	}
	if !updated.IsValid() {
		return nil, serrors.NewError("updated", serrors.ErrMissingTimestamp)
	}
	dbc := &DbCommon{
		ID:      id,
		Name:    name,
		Created: storers.TsToUInt64(created),
		Updated: storers.TsToUInt64(updated),
	}
	if validation.ValidString(desc) {
		dbc.Description = storers.StringPtr(desc)
	}
	if deleted.IsValid() {
		dbc.Deleted = storers.TsToUInt64Ptr(deleted)
	}
	return dbc, nil
}
