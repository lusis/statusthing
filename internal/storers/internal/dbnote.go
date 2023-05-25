package internal

import (
	"html"

	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/serrors"
	"github.com/lusis/statusthing/internal/storers"
	"github.com/lusis/statusthing/internal/validation"
)

// DbNote is a common representation of a [statusthingv1.Note] in a database
type DbNote struct {
	ID       string
	NoteText string
	ItemID   string
	*DbTimestamps
}

// DbNoteFromProto returns a [DbNote] from a [statusthingv1.Note]
func DbNoteFromProto(pbnote *statusthingv1.Note) (*DbNote, error) {
	if pbnote == nil {
		return nil, serrors.NewError("note", serrors.ErrNilVal)
	}
	id := html.EscapeString(pbnote.GetId())
	txt := html.EscapeString(pbnote.GetText())
	if !validation.ValidString(id) {
		return nil, serrors.NewError("id", serrors.ErrEmptyString)
	}
	if !validation.ValidString(txt) {
		return nil, serrors.NewError("text", serrors.ErrEmptyString)
	}
	ts, err := MakeDbTimestamps(pbnote.GetTimestamps())
	if err != nil {
		return nil, err
	}
	res := &DbNote{
		ID:           id,
		NoteText:     txt,
		DbTimestamps: ts,
	}
	// we set the item_id externally
	return res, nil
}

// ToProto returns a [statusthingv1.Note] from a [DbNote]
func (n *DbNote) ToProto() (*statusthingv1.Note, error) {
	res := &statusthingv1.Note{
		Timestamps: &statusthingv1.Timestamps{},
	}

	if !validation.ValidString(n.ID) {
		return nil, serrors.NewError("id", serrors.ErrInvalidData)
	}
	if !validation.ValidString(n.NoteText) {
		return nil, serrors.NewError("text", serrors.ErrInvalidData)
	}
	// id/name
	res.Id = html.UnescapeString(n.ID)
	res.Text = html.UnescapeString(n.NoteText)
	// timestamps
	pbcreated := storers.Int64ToTs(int64(n.Created))
	pbupdated := storers.Int64ToTs(int64(n.Updated))

	if pbcreated == nil {
		return nil, serrors.NewError("created", serrors.ErrInvalidData)
	}
	if pbupdated == nil {
		return nil, serrors.NewError("updated", serrors.ErrInvalidData)
	}
	res.Timestamps.Created = pbcreated
	res.Timestamps.Updated = pbupdated

	if n.Deleted != nil {
		res.Timestamps.Deleted = storers.Int64ToTs(int64(*n.Deleted))
	}
	return res, nil
}
