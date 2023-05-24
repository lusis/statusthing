package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/bufbuild/connect-go"

	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	v1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/filters"
	serrors "github.com/lusis/statusthing/internal/serrors"
	"github.com/lusis/statusthing/internal/services"
)

// APIHandler is something that can handle api requests
type APIHandler struct {
	sts *services.StatusThingService
}

// NewAPIHandler creates a new api handler
func NewAPIHandler(svc *services.StatusThingService) (*APIHandler, error) {
	if svc == nil {
		return nil, fmt.Errorf("svc: %w", serrors.ErrNilVal)
	}

	handler := &APIHandler{
		sts: svc,
	}
	return handler, nil
}

// GetItem gets an Item by its Id
func (api *APIHandler) GetItem(ctx context.Context, req *connect.Request[v1.GetItemRequest]) (*connect.Response[v1.GetItemResponse], error) {
	itemID := req.Msg.GetItemId()
	res, err := api.sts.GetItem(ctx, itemID)
	if err != nil {
		return nil, handleError(err)
	}
	return connect.NewResponse(&statusthingv1.GetItemResponse{Item: res}), nil
}

// ListItems gets all known Item
func (api *APIHandler) ListItems(ctx context.Context, req *connect.Request[v1.ListItemsRequest]) (*connect.Response[v1.ListItemsResponse], error) {
	opts := []filters.FilterOption{}

	if req.Msg.GetStatusIds() != nil {
		opts = append(opts, filters.WithStatusIDs(req.Msg.GetStatusIds()...))
	}
	if req.Msg.GetKinds() != nil {
		opts = append(opts, filters.WithStatusKinds(req.Msg.GetKinds()...))
	}
	res, err := api.sts.FindItems(ctx, opts...)
	if err != nil {
		return nil, handleError(err)
	}
	return connect.NewResponse(&v1.ListItemsResponse{
		Items: res,
	}), nil
}

// AddItem adds a new Item
func (api *APIHandler) AddItem(ctx context.Context, req *connect.Request[v1.AddItemRequest]) (*connect.Response[v1.AddItemResponse], error) {
	msg := req.Msg
	name := msg.GetName()
	desc := msg.GetDescription()
	initialStatus := msg.GetInitialStatus()
	initialNote := msg.GetInitialNoteText()
	initialStatusID := msg.GetInitialStatusId()

	opts := []filters.FilterOption{}
	if strings.TrimSpace(desc) != "" {
		opts = append(opts, filters.WithDescription(desc))
	}
	if strings.TrimSpace(initialStatusID) != "" {
		opts = append(opts, filters.WithStatusID(initialStatusID))
	} else if initialStatus != nil {
		opts = append(opts, filters.WithStatus(initialStatus))
	}
	if strings.TrimSpace(initialNote) != "" {
		opts = append(opts, filters.WithNoteText(initialNote))
	}
	res, err := api.sts.NewItem(ctx, name, opts...)
	if err != nil {
		return nil, handleError(err)
	}
	return connect.NewResponse(&v1.AddItemResponse{
		Item: res,
	}), nil
}

// UpdateItem updates an existing Item
func (api *APIHandler) UpdateItem(ctx context.Context, req *connect.Request[v1.UpdateItemRequest]) (*connect.Response[v1.UpdateItemResponse], error) {
	msg := req.Msg
	itemID := msg.GetItemId()
	name := msg.GetName()
	desc := msg.GetDescription()
	statusID := msg.GetStatusId()

	opts := []filters.FilterOption{}

	if name != "" {
		opts = append(opts, filters.WithName(name))
	}
	if desc != "" {
		opts = append(opts, filters.WithDescription(desc))
	}
	if statusID != "" {
		opts = append(opts, filters.WithStatusID(statusID))
	}

	if err := api.sts.EditItem(ctx, itemID, opts...); err != nil {
		return nil, handleError(err)
	}
	return &connect.Response[v1.UpdateItemResponse]{}, nil
}

// DeleteItem deletes an exisiting Item
func (api *APIHandler) DeleteItem(ctx context.Context, req *connect.Request[v1.DeleteItemRequest]) (*connect.Response[v1.DeleteItemResponse], error) {
	itemID := req.Msg.GetItemId()
	if err := api.sts.DeleteItem(ctx, itemID); err != nil {
		return nil, handleError(err)
	}
	return &connect.Response[v1.DeleteItemResponse]{}, nil
}

// GetNote gets a Note by its Id
func (api *APIHandler) GetNote(ctx context.Context, req *connect.Request[v1.GetNoteRequest]) (*connect.Response[v1.GetNoteResponse], error) {
	noteID := req.Msg.GetNoteId()
	res, err := api.sts.GetNote(ctx, noteID)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&v1.GetNoteResponse{Note: res}), nil
}

// ListNotes gets all Note for an Item
func (api *APIHandler) ListNotes(ctx context.Context, req *connect.Request[v1.ListNotesRequest]) (*connect.Response[v1.ListNotesResponse], error) {
	itemID := req.Msg.GetItemId()
	res, err := api.sts.AllNotes(ctx, itemID)
	if err != nil {
		return nil, handleError(err)
	}
	return connect.NewResponse(&v1.ListNotesResponse{Notes: res}), nil
}

// AddNote adds a Note to an Item
func (api *APIHandler) AddNote(ctx context.Context, req *connect.Request[v1.AddNoteRequest]) (*connect.Response[v1.AddNoteResponse], error) {
	itemID := req.Msg.GetItemId()
	noteText := req.Msg.GetNoteText()
	res, err := api.sts.NewNote(ctx, itemID, noteText)
	if err != nil {
		return nil, handleError(err)
	}
	return connect.NewResponse(&v1.AddNoteResponse{Note: res}), nil
}

// UpdateNote edits an existing Note
func (api *APIHandler) UpdateNote(ctx context.Context, req *connect.Request[v1.UpdateNoteRequest]) (*connect.Response[v1.UpdateNoteResponse], error) {
	noteID := req.Msg.GetNoteId()
	noteText := req.Msg.GetNoteText()
	if err := api.sts.EditNote(ctx, noteID, noteText); err != nil {
		return nil, handleError(err)
	}
	return &connect.Response[v1.UpdateNoteResponse]{}, nil
}

// DeleteNote deletes a Note from an Item
func (api *APIHandler) DeleteNote(ctx context.Context, req *connect.Request[v1.DeleteNoteRequest]) (*connect.Response[v1.DeleteNoteResponse], error) {
	noteID := req.Msg.GetNoteId()
	if err := api.sts.DeleteNote(ctx, noteID); err != nil {
		return nil, handleError(err)
	}
	return &connect.Response[v1.DeleteNoteResponse]{}, nil
}

// GetStatus gets a Status by its Id
func (api *APIHandler) GetStatus(ctx context.Context, req *connect.Request[v1.GetStatusRequest]) (*connect.Response[v1.GetStatusResponse], error) {
	statusID := req.Msg.GetStatusId()
	res, err := api.sts.GetStatus(ctx, statusID)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&v1.GetStatusResponse{Status: res}), nil
}

// ListStatus gets all known Statuses
func (api *APIHandler) ListStatus(ctx context.Context, req *connect.Request[v1.ListStatusRequest]) (*connect.Response[v1.ListStatusResponse], error) {
	opts := []filters.FilterOption{}
	if req.Msg.GetKinds() != nil {
		opts = append(opts, filters.WithStatusKinds(req.Msg.Kinds...))
	}
	res, err := api.sts.AllStatuses(ctx, opts...)
	if err != nil {
		return nil, handleError(err)
	}
	return connect.NewResponse(&v1.ListStatusResponse{Statuses: res}), nil
}

// AddStatus adds a new Status
func (api *APIHandler) AddStatus(ctx context.Context, req *connect.Request[v1.AddStatusRequest]) (*connect.Response[v1.AddStatusResponse], error) {
	name := req.Msg.GetName()
	kind := req.Msg.GetKind()
	color := req.Msg.GetColor()
	desc := req.Msg.GetDescription()

	opts := []filters.FilterOption{}
	if strings.TrimSpace(color) != "" {
		opts = append(opts, filters.WithColor(color))
	}
	if strings.TrimSpace(desc) != "" {
		opts = append(opts, filters.WithDescription(desc))
	}
	res, err := api.sts.NewStatus(ctx, name, kind, opts...)
	if err != nil {
		return nil, handleError(err)
	}
	return connect.NewResponse(&v1.AddStatusResponse{Status: res}), nil
}

// UpdateStatus updates an existing Status
func (api *APIHandler) UpdateStatus(ctx context.Context, req *connect.Request[v1.UpdateStatusRequest]) (*connect.Response[v1.UpdateStatusResponse], error) {
	id := req.Msg.GetStatusId()
	opts := []filters.FilterOption{}
	if name := req.Msg.GetName(); name != "" {
		opts = append(opts, filters.WithName(name))
	}
	if kind := req.Msg.GetKind(); kind != statusthingv1.StatusKind_STATUS_KIND_UNKNOWN {
		opts = append(opts, filters.WithStatusKind(kind))
	}
	if description := req.Msg.GetDescription(); description != "" {
		opts = append(opts, filters.WithDescription(description))
	}
	if color := req.Msg.GetColor(); color != "" {
		opts = append(opts, filters.WithColor(color))
	}

	if err := api.sts.EditStatus(ctx, id, opts...); err != nil {
		return nil, handleError(err)
	}
	return &connect.Response[v1.UpdateStatusResponse]{}, nil
}

// DeleteStatus deletes a Status
func (api *APIHandler) DeleteStatus(ctx context.Context, req *connect.Request[v1.DeleteStatusRequest]) (*connect.Response[v1.DeleteStatusResponse], error) {
	statusID := req.Msg.GetStatusId()
	if err := api.sts.DeleteStatus(ctx, statusID); err != nil {
		return nil, handleError(err)
	}
	return &connect.Response[v1.DeleteStatusResponse]{}, nil
}

func handleError(err error) *connect.Error {
	if errors.Is(err, serrors.ErrEmptyString) {
		return connect.NewError(connect.CodeInvalidArgument, err)
	}

	if errors.Is(err, serrors.ErrNotImplemented) {
		return connect.NewError(connect.CodeUnimplemented, err)
	}

	if errors.Is(err, serrors.ErrNotFound) {
		return connect.NewError(connect.CodeInvalidArgument, err)
	}
	if errors.Is(err, serrors.ErrStoreUnavailable) {
		return connect.NewError(connect.CodeFailedPrecondition, err)
	}
	// fallthrough
	return connect.NewError(connect.CodeInternal, fmt.Errorf("unexpected error"))
}
