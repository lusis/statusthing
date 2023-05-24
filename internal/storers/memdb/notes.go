package memdb

import (
	"context"
	"fmt"
	"strings"

	hcmemdb "github.com/hashicorp/go-memdb"
	"google.golang.org/protobuf/types/known/timestamppb"

	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	v1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/filters"
	"github.com/lusis/statusthing/internal/serrors"
)

const (
	notesTableName    = "notes"
	notesPrimaryIndex = "id"
)

type dbNote struct {
	ID       string
	ItemID   string
	NoteData string
	Created  int
	Updated  int
	Deleted  int
}

func (n *dbNote) toProto() (*v1.Note, error) {
	if n.ID == "" {
		return nil, serrors.ErrInvalidData
	}
	if n.NoteData == "" {
		return nil, serrors.ErrInvalidData
	}
	note := &v1.Note{
		Id:         n.ID,
		Text:       n.NoteData,
		Timestamps: &v1.Timestamps{},
	}
	// timestamps
	created, updated, deleted := intToTs(n.Created), intToTs(n.Updated), intToTs(n.Deleted)
	if !created.IsValid() {
		return nil, fmt.Errorf("created: %w", serrors.ErrInvalidData)
	}
	note.Timestamps.Created = created
	if !updated.IsValid() {
		return nil, fmt.Errorf("updated: %w", serrors.ErrInvalidData)
	}
	note.Timestamps.Updated = updated

	if deleted.IsValid() {
		note.Timestamps.Deleted = deleted
	}
	return note, nil
}

func noteFromProto(pn *statusthingv1.Note, itemID string) (*dbNote, error) {
	if pn == nil {
		return nil, fmt.Errorf("note: %w", serrors.ErrNilVal)
	}
	if itemID == "" {
		return nil, serrors.ErrEmptyString
	}
	if pn.GetId() == "" {
		return nil, fmt.Errorf("id: %w", serrors.ErrEmptyString)
	}
	if pn.GetText() == "" {
		return nil, fmt.Errorf("text: %w", serrors.ErrEmptyString)
	}
	if pn.GetTimestamps() == nil {
		return nil, fmt.Errorf("timestamps: %w", serrors.ErrNilVal)
	}
	if !pn.GetTimestamps().GetCreated().IsValid() {
		return nil, serrors.ErrInvalidData
	}
	if !pn.GetTimestamps().GetUpdated().IsValid() {
		return nil, serrors.ErrInvalidData
	}
	n := &dbNote{
		ID:       pn.GetId(),
		ItemID:   itemID,
		NoteData: pn.GetText(),
		Created:  tsToInt(pn.GetTimestamps().GetCreated()),
		Updated:  tsToInt(pn.GetTimestamps().GetUpdated()),
	}
	if pn.GetTimestamps().GetDeleted().IsValid() {
		n.Deleted = tsToInt(pn.GetTimestamps().GetDeleted())
	}
	return n, nil
}

var notesSchema = &hcmemdb.TableSchema{
	Name: notesTableName,
	Indexes: map[string]*hcmemdb.IndexSchema{
		"id": {
			Name:    notesPrimaryIndex,
			Unique:  true,
			Indexer: &hcmemdb.StringFieldIndex{Field: "ID"},
		},
		"item_id": {
			Name:    "item_id",
			Indexer: &hcmemdb.StringFieldIndex{Field: "ItemID"},
		},
		"note_data": {
			Name:    "note_data",
			Indexer: &hcmemdb.StringFieldIndex{Field: "NoteData"},
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

// GetNote gets a [statusthingv1.Note] by its unique id
func (s *StatusThingStore) GetNote(ctx context.Context, noteID string) (*v1.Note, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()
	res, err := firstWithTxn(ctx, txn, notesTableName, notesPrimaryIndex, noteID)
	if err != nil {
		return nil, err
	}

	note, ok := res.(*dbNote)
	if !ok {
		return nil, serrors.ErrInvalidData
	}

	pbnote, err := note.toProto()
	if err != nil {
		return nil, err
	}
	return pbnote, nil
}

// StoreNote stores the provided [statusthingv1.Note] associated with the provided [statusthingv1.Item] by its id
func (s *StatusThingStore) StoreNote(ctx context.Context, note *v1.Note, itemID string) (*v1.Note, error) {
	if note == nil {
		return nil, serrors.ErrNilVal
	}
	if strings.TrimSpace(itemID) == "" {
		return nil, serrors.ErrEmptyString
	}
	if _, err := s.GetItem(ctx, itemID); err != nil {
		return nil, err
	}
	txn := s.db.Txn(true)
	memNote, err := noteFromProto(note, itemID)
	if err != nil {
		return nil, err
	}
	if err := txn.Insert(notesTableName, memNote); err != nil {
		txn.Abort()
		return nil, err
	}
	txn.Commit()
	return note, nil
}

// FindNotes gets all known [statusthingv1.Note]
func (s *StatusThingStore) FindNotes(ctx context.Context, itemID string, opts ...filters.FilterOption) ([]*v1.Note, error) {
	// we don't use any options right now for notes
	_, err := filters.New(opts...)
	if err != nil {
		return nil, err
	}

	if itemID == "" {
		return nil, serrors.ErrEmptyString
	}

	if _, err := s.GetItem(ctx, itemID); err != nil {
		return nil, err
	}

	txn := s.db.Txn(false)
	defer txn.Abort()
	return getNotesWithTxn(ctx, txn, itemID, opts...)
}

// UpdateNote updates the [statusthingv1.Note] with the provided [filters.FilterOption]
// supported filters:
// - [filters.WithNoteText]: for changing the note text
// - [filters.WithTimestamps]: for setting custom timestamps
func (s *StatusThingStore) UpdateNote(ctx context.Context, noteID string, opts ...filters.FilterOption) error {
	if len(opts) == 0 {
		return serrors.ErrAtLeastOne
	}
	if strings.TrimSpace(noteID) == "" {
		return serrors.ErrEmptyString
	}
	f, err := filters.New(opts...)
	if err != nil {
		return err
	}
	// we need the raw db record here so we can get the item its associated with for not-found checks
	txn := s.db.Txn(false)
	existing, err := firstWithTxn(ctx, txn, notesTableName, notesPrimaryIndex, noteID)
	if err != nil {
		txn.Abort()
		return err
	}
	// done with the txn. We need to clear it so we can call storenote
	txn.Abort()

	note, ok := existing.(*dbNote)
	if !ok {
		return serrors.ErrInvalidData
	}
	pnote, err := note.toProto()
	if err != nil {
		return err
	}
	if f.Timestamps() != nil {
		pnote.Timestamps = f.Timestamps()
	} else {
		pnote.GetTimestamps().Updated = timestamppb.Now()
	}
	if f.NoteText() != "" {
		pnote.Text = f.NoteText()
	}
	if _, err := s.StoreNote(ctx, pnote, note.ItemID); err != nil {
		return err
	}
	return nil
}

// DeleteNote deletes a [statusthingv1.Note] by its id
func (s *StatusThingStore) DeleteNote(ctx context.Context, noteID string) error {
	txn := s.db.Txn(true)
	if err := deleteWithTxn(ctx, txn, notesTableName, &dbNote{ID: noteID}); err != nil {
		txn.Abort()
		return err
	}
	txn.Commit()
	return nil
}

func getNotesWithTxn(ctx context.Context, txn *hcmemdb.Txn, itemID string, opts ...filters.FilterOption) ([]*v1.Note, error) {
	if txn == nil {
		return nil, serrors.ErrNilVal
	}
	// we don't use any options right now for notes
	_, err := filters.New(opts...)
	if err != nil {
		return nil, err
	}

	if itemID == "" {
		return nil, serrors.ErrEmptyString
	}

	// we get by item id here
	res, err := getWithTxn(ctx, txn, notesTableName, "item_id", itemID)
	if err != nil {
		return nil, err
	}
	notes := []*v1.Note{}
	for obj := res.Next(); obj != nil; obj = res.Next() {
		note, ok := obj.(*dbNote)
		if !ok {
			return nil, serrors.ErrInvalidData
		}

		pbnote, err := note.toProto()
		if err != nil {
			return nil, err
		}
		notes = append(notes, pbnote)
	}
	return notes, nil
}
