// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: statusthing/v1/services.proto

package statusthingv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	ItemsService_GetItem_FullMethodName    = "/statusthing.v1.ItemsService/GetItem"
	ItemsService_ListItems_FullMethodName  = "/statusthing.v1.ItemsService/ListItems"
	ItemsService_AddItem_FullMethodName    = "/statusthing.v1.ItemsService/AddItem"
	ItemsService_UpdateItem_FullMethodName = "/statusthing.v1.ItemsService/UpdateItem"
	ItemsService_DeleteItem_FullMethodName = "/statusthing.v1.ItemsService/DeleteItem"
)

// ItemsServiceClient is the client API for ItemsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ItemsServiceClient interface {
	// GetItem gets an Item by its Id
	GetItem(ctx context.Context, in *GetItemRequest, opts ...grpc.CallOption) (*GetItemResponse, error)
	// ListItems gets all known Items
	ListItems(ctx context.Context, in *ListItemsRequest, opts ...grpc.CallOption) (*ListItemsResponse, error)
	// AddItem adds a new Item
	AddItem(ctx context.Context, in *AddItemRequest, opts ...grpc.CallOption) (*AddItemResponse, error)
	// UpdateItem updates an existing Item
	UpdateItem(ctx context.Context, in *UpdateItemRequest, opts ...grpc.CallOption) (*UpdateItemResponse, error)
	// DeleteItem deletes an exisiting Item
	DeleteItem(ctx context.Context, in *DeleteItemRequest, opts ...grpc.CallOption) (*DeleteItemResponse, error)
}

type itemsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewItemsServiceClient(cc grpc.ClientConnInterface) ItemsServiceClient {
	return &itemsServiceClient{cc}
}

func (c *itemsServiceClient) GetItem(ctx context.Context, in *GetItemRequest, opts ...grpc.CallOption) (*GetItemResponse, error) {
	out := new(GetItemResponse)
	err := c.cc.Invoke(ctx, ItemsService_GetItem_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *itemsServiceClient) ListItems(ctx context.Context, in *ListItemsRequest, opts ...grpc.CallOption) (*ListItemsResponse, error) {
	out := new(ListItemsResponse)
	err := c.cc.Invoke(ctx, ItemsService_ListItems_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *itemsServiceClient) AddItem(ctx context.Context, in *AddItemRequest, opts ...grpc.CallOption) (*AddItemResponse, error) {
	out := new(AddItemResponse)
	err := c.cc.Invoke(ctx, ItemsService_AddItem_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *itemsServiceClient) UpdateItem(ctx context.Context, in *UpdateItemRequest, opts ...grpc.CallOption) (*UpdateItemResponse, error) {
	out := new(UpdateItemResponse)
	err := c.cc.Invoke(ctx, ItemsService_UpdateItem_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *itemsServiceClient) DeleteItem(ctx context.Context, in *DeleteItemRequest, opts ...grpc.CallOption) (*DeleteItemResponse, error) {
	out := new(DeleteItemResponse)
	err := c.cc.Invoke(ctx, ItemsService_DeleteItem_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ItemsServiceServer is the server API for ItemsService service.
// All implementations must embed UnimplementedItemsServiceServer
// for forward compatibility
type ItemsServiceServer interface {
	// GetItem gets an Item by its Id
	GetItem(context.Context, *GetItemRequest) (*GetItemResponse, error)
	// ListItems gets all known Items
	ListItems(context.Context, *ListItemsRequest) (*ListItemsResponse, error)
	// AddItem adds a new Item
	AddItem(context.Context, *AddItemRequest) (*AddItemResponse, error)
	// UpdateItem updates an existing Item
	UpdateItem(context.Context, *UpdateItemRequest) (*UpdateItemResponse, error)
	// DeleteItem deletes an exisiting Item
	DeleteItem(context.Context, *DeleteItemRequest) (*DeleteItemResponse, error)
	mustEmbedUnimplementedItemsServiceServer()
}

// UnimplementedItemsServiceServer must be embedded to have forward compatible implementations.
type UnimplementedItemsServiceServer struct {
}

func (UnimplementedItemsServiceServer) GetItem(context.Context, *GetItemRequest) (*GetItemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetItem not implemented")
}
func (UnimplementedItemsServiceServer) ListItems(context.Context, *ListItemsRequest) (*ListItemsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListItems not implemented")
}
func (UnimplementedItemsServiceServer) AddItem(context.Context, *AddItemRequest) (*AddItemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddItem not implemented")
}
func (UnimplementedItemsServiceServer) UpdateItem(context.Context, *UpdateItemRequest) (*UpdateItemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateItem not implemented")
}
func (UnimplementedItemsServiceServer) DeleteItem(context.Context, *DeleteItemRequest) (*DeleteItemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteItem not implemented")
}
func (UnimplementedItemsServiceServer) mustEmbedUnimplementedItemsServiceServer() {}

// UnsafeItemsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ItemsServiceServer will
// result in compilation errors.
type UnsafeItemsServiceServer interface {
	mustEmbedUnimplementedItemsServiceServer()
}

func RegisterItemsServiceServer(s grpc.ServiceRegistrar, srv ItemsServiceServer) {
	s.RegisterService(&ItemsService_ServiceDesc, srv)
}

func _ItemsService_GetItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemsServiceServer).GetItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ItemsService_GetItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemsServiceServer).GetItem(ctx, req.(*GetItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ItemsService_ListItems_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListItemsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemsServiceServer).ListItems(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ItemsService_ListItems_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemsServiceServer).ListItems(ctx, req.(*ListItemsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ItemsService_AddItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemsServiceServer).AddItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ItemsService_AddItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemsServiceServer).AddItem(ctx, req.(*AddItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ItemsService_UpdateItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemsServiceServer).UpdateItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ItemsService_UpdateItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemsServiceServer).UpdateItem(ctx, req.(*UpdateItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ItemsService_DeleteItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemsServiceServer).DeleteItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ItemsService_DeleteItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemsServiceServer).DeleteItem(ctx, req.(*DeleteItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ItemsService_ServiceDesc is the grpc.ServiceDesc for ItemsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ItemsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "statusthing.v1.ItemsService",
	HandlerType: (*ItemsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetItem",
			Handler:    _ItemsService_GetItem_Handler,
		},
		{
			MethodName: "ListItems",
			Handler:    _ItemsService_ListItems_Handler,
		},
		{
			MethodName: "AddItem",
			Handler:    _ItemsService_AddItem_Handler,
		},
		{
			MethodName: "UpdateItem",
			Handler:    _ItemsService_UpdateItem_Handler,
		},
		{
			MethodName: "DeleteItem",
			Handler:    _ItemsService_DeleteItem_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "statusthing/v1/services.proto",
}

const (
	StatusService_GetStatus_FullMethodName    = "/statusthing.v1.StatusService/GetStatus"
	StatusService_ListStatus_FullMethodName   = "/statusthing.v1.StatusService/ListStatus"
	StatusService_AddStatus_FullMethodName    = "/statusthing.v1.StatusService/AddStatus"
	StatusService_UpdateStatus_FullMethodName = "/statusthing.v1.StatusService/UpdateStatus"
	StatusService_DeleteStatus_FullMethodName = "/statusthing.v1.StatusService/DeleteStatus"
)

// StatusServiceClient is the client API for StatusService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StatusServiceClient interface {
	// GetStatus gets a Status by its Id
	GetStatus(ctx context.Context, in *GetStatusRequest, opts ...grpc.CallOption) (*GetStatusResponse, error)
	// ListStatus gets all known Status
	ListStatus(ctx context.Context, in *ListStatusRequest, opts ...grpc.CallOption) (*ListStatusResponse, error)
	// AddStatus adds a new status
	AddStatus(ctx context.Context, in *AddStatusRequest, opts ...grpc.CallOption) (*AddStatusResponse, error)
	// UpdateStatus updates an existing status
	UpdateStatus(ctx context.Context, in *UpdateStatusRequest, opts ...grpc.CallOption) (*UpdateStatusResponse, error)
	// DeleteStatus deletes a Status
	DeleteStatus(ctx context.Context, in *DeleteStatusRequest, opts ...grpc.CallOption) (*DeleteStatusResponse, error)
}

type statusServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStatusServiceClient(cc grpc.ClientConnInterface) StatusServiceClient {
	return &statusServiceClient{cc}
}

func (c *statusServiceClient) GetStatus(ctx context.Context, in *GetStatusRequest, opts ...grpc.CallOption) (*GetStatusResponse, error) {
	out := new(GetStatusResponse)
	err := c.cc.Invoke(ctx, StatusService_GetStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statusServiceClient) ListStatus(ctx context.Context, in *ListStatusRequest, opts ...grpc.CallOption) (*ListStatusResponse, error) {
	out := new(ListStatusResponse)
	err := c.cc.Invoke(ctx, StatusService_ListStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statusServiceClient) AddStatus(ctx context.Context, in *AddStatusRequest, opts ...grpc.CallOption) (*AddStatusResponse, error) {
	out := new(AddStatusResponse)
	err := c.cc.Invoke(ctx, StatusService_AddStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statusServiceClient) UpdateStatus(ctx context.Context, in *UpdateStatusRequest, opts ...grpc.CallOption) (*UpdateStatusResponse, error) {
	out := new(UpdateStatusResponse)
	err := c.cc.Invoke(ctx, StatusService_UpdateStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statusServiceClient) DeleteStatus(ctx context.Context, in *DeleteStatusRequest, opts ...grpc.CallOption) (*DeleteStatusResponse, error) {
	out := new(DeleteStatusResponse)
	err := c.cc.Invoke(ctx, StatusService_DeleteStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StatusServiceServer is the server API for StatusService service.
// All implementations must embed UnimplementedStatusServiceServer
// for forward compatibility
type StatusServiceServer interface {
	// GetStatus gets a Status by its Id
	GetStatus(context.Context, *GetStatusRequest) (*GetStatusResponse, error)
	// ListStatus gets all known Status
	ListStatus(context.Context, *ListStatusRequest) (*ListStatusResponse, error)
	// AddStatus adds a new status
	AddStatus(context.Context, *AddStatusRequest) (*AddStatusResponse, error)
	// UpdateStatus updates an existing status
	UpdateStatus(context.Context, *UpdateStatusRequest) (*UpdateStatusResponse, error)
	// DeleteStatus deletes a Status
	DeleteStatus(context.Context, *DeleteStatusRequest) (*DeleteStatusResponse, error)
	mustEmbedUnimplementedStatusServiceServer()
}

// UnimplementedStatusServiceServer must be embedded to have forward compatible implementations.
type UnimplementedStatusServiceServer struct {
}

func (UnimplementedStatusServiceServer) GetStatus(context.Context, *GetStatusRequest) (*GetStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStatus not implemented")
}
func (UnimplementedStatusServiceServer) ListStatus(context.Context, *ListStatusRequest) (*ListStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListStatus not implemented")
}
func (UnimplementedStatusServiceServer) AddStatus(context.Context, *AddStatusRequest) (*AddStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddStatus not implemented")
}
func (UnimplementedStatusServiceServer) UpdateStatus(context.Context, *UpdateStatusRequest) (*UpdateStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateStatus not implemented")
}
func (UnimplementedStatusServiceServer) DeleteStatus(context.Context, *DeleteStatusRequest) (*DeleteStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteStatus not implemented")
}
func (UnimplementedStatusServiceServer) mustEmbedUnimplementedStatusServiceServer() {}

// UnsafeStatusServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StatusServiceServer will
// result in compilation errors.
type UnsafeStatusServiceServer interface {
	mustEmbedUnimplementedStatusServiceServer()
}

func RegisterStatusServiceServer(s grpc.ServiceRegistrar, srv StatusServiceServer) {
	s.RegisterService(&StatusService_ServiceDesc, srv)
}

func _StatusService_GetStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatusServiceServer).GetStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StatusService_GetStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatusServiceServer).GetStatus(ctx, req.(*GetStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StatusService_ListStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatusServiceServer).ListStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StatusService_ListStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatusServiceServer).ListStatus(ctx, req.(*ListStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StatusService_AddStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatusServiceServer).AddStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StatusService_AddStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatusServiceServer).AddStatus(ctx, req.(*AddStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StatusService_UpdateStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatusServiceServer).UpdateStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StatusService_UpdateStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatusServiceServer).UpdateStatus(ctx, req.(*UpdateStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StatusService_DeleteStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatusServiceServer).DeleteStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StatusService_DeleteStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatusServiceServer).DeleteStatus(ctx, req.(*DeleteStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// StatusService_ServiceDesc is the grpc.ServiceDesc for StatusService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StatusService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "statusthing.v1.StatusService",
	HandlerType: (*StatusServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetStatus",
			Handler:    _StatusService_GetStatus_Handler,
		},
		{
			MethodName: "ListStatus",
			Handler:    _StatusService_ListStatus_Handler,
		},
		{
			MethodName: "AddStatus",
			Handler:    _StatusService_AddStatus_Handler,
		},
		{
			MethodName: "UpdateStatus",
			Handler:    _StatusService_UpdateStatus_Handler,
		},
		{
			MethodName: "DeleteStatus",
			Handler:    _StatusService_DeleteStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "statusthing/v1/services.proto",
}

const (
	NotesService_GetNote_FullMethodName    = "/statusthing.v1.NotesService/GetNote"
	NotesService_ListNotes_FullMethodName  = "/statusthing.v1.NotesService/ListNotes"
	NotesService_AddNote_FullMethodName    = "/statusthing.v1.NotesService/AddNote"
	NotesService_UpdateNote_FullMethodName = "/statusthing.v1.NotesService/UpdateNote"
	NotesService_DeleteNote_FullMethodName = "/statusthing.v1.NotesService/DeleteNote"
)

// NotesServiceClient is the client API for NotesService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NotesServiceClient interface {
	// GetNote gets a Note by its Id
	GetNote(ctx context.Context, in *GetNoteRequest, opts ...grpc.CallOption) (*GetNoteResponse, error)
	// ListNotes gets all Note for an Item
	ListNotes(ctx context.Context, in *ListNotesRequest, opts ...grpc.CallOption) (*ListNotesResponse, error)
	// AddNote adds a Note to an Item
	AddNote(ctx context.Context, in *AddNoteRequest, opts ...grpc.CallOption) (*AddNoteResponse, error)
	// UpdateNote updates an existing Note
	UpdateNote(ctx context.Context, in *UpdateNoteRequest, opts ...grpc.CallOption) (*UpdateNoteResponse, error)
	// DeleteNote deletes a Note from an Item
	DeleteNote(ctx context.Context, in *DeleteNoteRequest, opts ...grpc.CallOption) (*DeleteNoteResponse, error)
}

type notesServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNotesServiceClient(cc grpc.ClientConnInterface) NotesServiceClient {
	return &notesServiceClient{cc}
}

func (c *notesServiceClient) GetNote(ctx context.Context, in *GetNoteRequest, opts ...grpc.CallOption) (*GetNoteResponse, error) {
	out := new(GetNoteResponse)
	err := c.cc.Invoke(ctx, NotesService_GetNote_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notesServiceClient) ListNotes(ctx context.Context, in *ListNotesRequest, opts ...grpc.CallOption) (*ListNotesResponse, error) {
	out := new(ListNotesResponse)
	err := c.cc.Invoke(ctx, NotesService_ListNotes_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notesServiceClient) AddNote(ctx context.Context, in *AddNoteRequest, opts ...grpc.CallOption) (*AddNoteResponse, error) {
	out := new(AddNoteResponse)
	err := c.cc.Invoke(ctx, NotesService_AddNote_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notesServiceClient) UpdateNote(ctx context.Context, in *UpdateNoteRequest, opts ...grpc.CallOption) (*UpdateNoteResponse, error) {
	out := new(UpdateNoteResponse)
	err := c.cc.Invoke(ctx, NotesService_UpdateNote_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notesServiceClient) DeleteNote(ctx context.Context, in *DeleteNoteRequest, opts ...grpc.CallOption) (*DeleteNoteResponse, error) {
	out := new(DeleteNoteResponse)
	err := c.cc.Invoke(ctx, NotesService_DeleteNote_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NotesServiceServer is the server API for NotesService service.
// All implementations must embed UnimplementedNotesServiceServer
// for forward compatibility
type NotesServiceServer interface {
	// GetNote gets a Note by its Id
	GetNote(context.Context, *GetNoteRequest) (*GetNoteResponse, error)
	// ListNotes gets all Note for an Item
	ListNotes(context.Context, *ListNotesRequest) (*ListNotesResponse, error)
	// AddNote adds a Note to an Item
	AddNote(context.Context, *AddNoteRequest) (*AddNoteResponse, error)
	// UpdateNote updates an existing Note
	UpdateNote(context.Context, *UpdateNoteRequest) (*UpdateNoteResponse, error)
	// DeleteNote deletes a Note from an Item
	DeleteNote(context.Context, *DeleteNoteRequest) (*DeleteNoteResponse, error)
	mustEmbedUnimplementedNotesServiceServer()
}

// UnimplementedNotesServiceServer must be embedded to have forward compatible implementations.
type UnimplementedNotesServiceServer struct {
}

func (UnimplementedNotesServiceServer) GetNote(context.Context, *GetNoteRequest) (*GetNoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNote not implemented")
}
func (UnimplementedNotesServiceServer) ListNotes(context.Context, *ListNotesRequest) (*ListNotesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListNotes not implemented")
}
func (UnimplementedNotesServiceServer) AddNote(context.Context, *AddNoteRequest) (*AddNoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddNote not implemented")
}
func (UnimplementedNotesServiceServer) UpdateNote(context.Context, *UpdateNoteRequest) (*UpdateNoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateNote not implemented")
}
func (UnimplementedNotesServiceServer) DeleteNote(context.Context, *DeleteNoteRequest) (*DeleteNoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteNote not implemented")
}
func (UnimplementedNotesServiceServer) mustEmbedUnimplementedNotesServiceServer() {}

// UnsafeNotesServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NotesServiceServer will
// result in compilation errors.
type UnsafeNotesServiceServer interface {
	mustEmbedUnimplementedNotesServiceServer()
}

func RegisterNotesServiceServer(s grpc.ServiceRegistrar, srv NotesServiceServer) {
	s.RegisterService(&NotesService_ServiceDesc, srv)
}

func _NotesService_GetNote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetNoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotesServiceServer).GetNote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotesService_GetNote_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotesServiceServer).GetNote(ctx, req.(*GetNoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotesService_ListNotes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListNotesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotesServiceServer).ListNotes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotesService_ListNotes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotesServiceServer).ListNotes(ctx, req.(*ListNotesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotesService_AddNote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddNoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotesServiceServer).AddNote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotesService_AddNote_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotesServiceServer).AddNote(ctx, req.(*AddNoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotesService_UpdateNote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateNoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotesServiceServer).UpdateNote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotesService_UpdateNote_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotesServiceServer).UpdateNote(ctx, req.(*UpdateNoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotesService_DeleteNote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteNoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotesServiceServer).DeleteNote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotesService_DeleteNote_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotesServiceServer).DeleteNote(ctx, req.(*DeleteNoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// NotesService_ServiceDesc is the grpc.ServiceDesc for NotesService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NotesService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "statusthing.v1.NotesService",
	HandlerType: (*NotesServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetNote",
			Handler:    _NotesService_GetNote_Handler,
		},
		{
			MethodName: "ListNotes",
			Handler:    _NotesService_ListNotes_Handler,
		},
		{
			MethodName: "AddNote",
			Handler:    _NotesService_AddNote_Handler,
		},
		{
			MethodName: "UpdateNote",
			Handler:    _NotesService_UpdateNote_Handler,
		},
		{
			MethodName: "DeleteNote",
			Handler:    _NotesService_DeleteNote_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "statusthing/v1/services.proto",
}
