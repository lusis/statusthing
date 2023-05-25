package internal

import statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"

// DbNote is a common representation of a [statusthingv1.Note] in a database
type DbNote struct {
	*DbCommon
	ItemID string
}

// DbNoteFromProto returns a [DbNote] from a [statusthingv1.Note]
func DbNoteFromProto(note *statusthingv1.Note) (*DbNote, error) {
	return nil, nil
}

// ToProto returns a [statusthingv1.Note] from a [DbNote]
func (n *DbNote) ToProto() (*statusthingv1.Note, error) {
	return nil, nil
}
