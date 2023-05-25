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
	*DbTimestamps
}

// DbTimestamps represents a [statusthingv1.Timestamps] as a db record
type DbTimestamps struct {
	Created uint64  `db:"created"`
	Updated uint64  `db:"updated"`
	Deleted *uint64 `db:"deleted"`
}

// MakeDbTimestamps converts a [DbTimestamps] to a [statusthingv1.Timestamps]
func MakeDbTimestamps(timestamps *statusthingv1.Timestamps) (*DbTimestamps, error) {
	created := timestamps.GetCreated()
	updated := timestamps.GetUpdated()
	deleted := timestamps.GetDeleted()
	if timestamps == nil {
		return nil, serrors.NewError("timestamps", serrors.ErrMissingTimestamp)
	}
	if !created.IsValid() {
		return nil, serrors.NewError("created", serrors.ErrMissingTimestamp)
	}
	if !updated.IsValid() {
		return nil, serrors.NewError("updated", serrors.ErrMissingTimestamp)
	}
	dbt := &DbTimestamps{
		Created: storers.TsToUInt64(created),
		Updated: storers.TsToUInt64(updated),
	}
	if deleted.IsValid() {
		dbt.Deleted = storers.TsToUInt64Ptr(deleted)
	}
	return dbt, nil
}

// MakeDbCommon builds the minimum required common fields of all db records
func MakeDbCommon(id, name, desc string, timestamps *statusthingv1.Timestamps) (*DbCommon, error) {
	id = html.EscapeString(id)
	name = html.EscapeString(name)
	desc = html.EscapeString(desc)

	if !validation.ValidString(id) {
		return nil, serrors.NewError("id", serrors.ErrEmptyString)
	}
	if !validation.ValidString(name) {
		return nil, serrors.NewError("name", serrors.ErrEmptyString)
	}

	dbc := &DbCommon{
		ID:   id,
		Name: name,
	}
	if validation.ValidString(desc) {
		dbc.Description = storers.StringPtr(desc)
	}
	ts, tserr := MakeDbTimestamps(timestamps)
	if tserr != nil {
		return nil, tserr
	}
	dbc.DbTimestamps = ts
	return dbc, nil
}
