package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bufbuild/connect-go"
	"github.com/go-chi/chi"
	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	v1connect "github.com/lusis/statusthing/gen/go/statusthing/v1/statusthingv1connect"
	"github.com/lusis/statusthing/internal/filters"
	"github.com/lusis/statusthing/internal/serrors"
	"github.com/lusis/statusthing/internal/services"
	"github.com/lusis/statusthing/internal/storers/memdb"

	"github.com/stretchr/testify/require"
)

func TestImplements(t *testing.T) {
	require.Implements(t, (*v1connect.ItemsServiceHandler)(nil), new(APIHandler), "should satisfy items rpc interface")
	require.Implements(t, (*v1connect.NotesServiceHandler)(nil), new(APIHandler), "should satisfy notes rpc interface")
	require.Implements(t, (*v1connect.StatusServiceHandler)(nil), new(APIHandler), "should sastify statuses rpc interface")
}

func TestNew(t *testing.T) {
	t.Parallel()
	t.Run("nil-svc", func(t *testing.T) {
		res, err := NewAPIHandler(nil)
		require.ErrorIs(t, err, serrors.ErrNilVal)
		require.Nil(t, res)
	})
}
func TestAddStatus(t *testing.T) {
	t.Parallel()
	t.Run("happy-path", func(t *testing.T) {
		// arrange
		ctx := context.TODO()
		status := statusthingv1.Status{
			Name:        t.Name(),
			Kind:        statusthingv1.StatusKind_STATUS_KIND_UNAVAILABLE,
			Color:       "green",
			Description: "my-status",
		}
		api, _, httpSrv, err := apiTestSetup(t)
		defer httpSrv.Close()
		require.NoError(t, err)

		// act
		res, err := api.AddStatus(ctx, connect.NewRequest(&statusthingv1.AddStatusRequest{
			Name:        status.GetName(),
			Kind:        status.GetKind(),
			Color:       status.GetColor(),
			Description: status.GetDescription(),
		}))

		// assert
		require.NoError(t, err)
		require.NotNil(t, res.Msg.GetStatus())
		require.Equal(t, status.GetName(), res.Msg.Status.GetName())
		require.Equal(t, status.GetKind(), res.Msg.GetStatus().GetKind())
		require.Equal(t, status.GetColor(), res.Msg.GetStatus().GetColor())
		require.Equal(t, status.GetDescription(), res.Msg.GetStatus().GetDescription())
		require.NotEmpty(t, res.Msg.GetStatus().GetId())
	})
	t.Run("invalid-name", func(t *testing.T) {
		// arrange
		ctx := context.TODO()
		status := statusthingv1.Status{
			Name: "",
			Kind: statusthingv1.StatusKind_STATUS_KIND_UNAVAILABLE,
		}
		api, _, httpSrv, err := apiTestSetup(t)
		defer httpSrv.Close()
		require.NoError(t, err)

		// act
		res, err := api.AddStatus(ctx, connect.NewRequest(&statusthingv1.AddStatusRequest{
			Name: status.GetName(),
			Kind: status.GetKind(),
		}))

		// assert
		require.ErrorIs(t, err, serrors.ErrEmptyString)
		require.Nil(t, res)
	})
}
func TestAllStatus(t *testing.T) {
	t.Parallel()
	t.Run("happy-path", func(t *testing.T) {
		// arrange
		ctx := context.TODO()
		api, _, httpSrv, err := apiTestSetup(t)
		defer httpSrv.Close()
		require.NoError(t, err)
		api.sts.CreateDefaultStatuses()

		// act
		res, err := api.ListStatus(ctx, connect.NewRequest(&statusthingv1.ListStatusRequest{}))

		// assert
		require.NoError(t, err)
		require.Len(t, res.Msg.GetStatuses(), 3)
	})
	t.Run("with-kinds", func(t *testing.T) {
		// arrange
		ctx := context.TODO()
		api, _, httpSrv, err := apiTestSetup(t)
		defer httpSrv.Close()
		require.NoError(t, err)
		api.sts.CreateDefaultStatuses()

		// act
		// the defaults status include one of each up,down,warning
		for _, k := range []statusthingv1.StatusKind{
			statusthingv1.StatusKind_STATUS_KIND_UP,
			statusthingv1.StatusKind_STATUS_KIND_DOWN,
			statusthingv1.StatusKind_STATUS_KIND_WARNING,
		} {
			res, err := api.ListStatus(ctx, connect.NewRequest(&statusthingv1.ListStatusRequest{Kinds: []statusthingv1.StatusKind{k}}))

			// assert
			require.NoError(t, err)
			require.Lenf(t, res.Msg.GetStatuses(), 1, "should have one item of kind: %s", k.String())
		}
	})
	t.Run("no-statuses", func(t *testing.T) {
		// arrange
		ctx := context.TODO()
		api, _, httpSrv, err := apiTestSetup(t)
		defer httpSrv.Close()
		require.NoError(t, err)

		// act
		res, err := api.ListStatus(ctx, connect.NewRequest(&statusthingv1.ListStatusRequest{}))

		// assert
		require.NoError(t, err)
		require.Len(t, res.Msg.GetStatuses(), 0)
	})
}
func TestUpdateStatus(t *testing.T) {
	t.Parallel()
	t.Run("happy-path", func(t *testing.T) {
		ctx := context.TODO()
		api, _, httpSrv, err := apiTestSetup(t)
		defer httpSrv.Close()
		require.NoError(t, err)

		istatus, err := api.AddStatus(ctx, connect.NewRequest(&statusthingv1.AddStatusRequest{
			Name: t.Name(),
			Kind: statusthingv1.StatusKind_STATUS_KIND_AVAILABLE,
		}))
		require.NoError(t, err)
		require.NotNil(t, istatus)
		require.Equal(t, t.Name(), istatus.Msg.GetStatus().GetName())
		require.Equal(t, statusthingv1.StatusKind_STATUS_KIND_AVAILABLE, istatus.Msg.GetStatus().GetKind())

		// this doesn't return anything message wise that we care about
		_, uerr := api.UpdateStatus(ctx, connect.NewRequest(&statusthingv1.UpdateStatusRequest{
			StatusId:    istatus.Msg.GetStatus().GetId(),
			Name:        "new-name",
			Description: "new-description",
			Color:       "new-color",
			Kind:        statusthingv1.StatusKind_STATUS_KIND_DOWN,
		}))
		require.NoError(t, uerr)
	})
	t.Run("not-found", func(t *testing.T) {
		ctx := context.TODO()
		api, _, httpSrv, err := apiTestSetup(t)
		defer httpSrv.Close()
		require.NoError(t, err)

		// this doesn't return anything message wise that we care about
		_, uerr := api.UpdateStatus(ctx, connect.NewRequest(&statusthingv1.UpdateStatusRequest{
			StatusId:    "non-existent",
			Name:        "new-name",
			Description: "new-description",
			Color:       "new-color",
			Kind:        statusthingv1.StatusKind_STATUS_KIND_DOWN,
		}))
		require.ErrorIs(t, uerr, serrors.ErrNotFound)
	})
}
func TestDeleteStatus(t *testing.T) {
	t.Parallel()
	t.Run("happy-path", func(t *testing.T) {
		ctx := context.TODO()
		api, _, httpSrv, err := apiTestSetup(t)
		defer httpSrv.Close()
		require.NoError(t, err)

		err = api.sts.CreateDefaultStatuses()
		require.NoError(t, err)

		created := api.sts.GetCreatedDefaults()
		require.Len(t, created, 3)
		first := created[0]

		req := connect.NewRequest(&statusthingv1.DeleteStatusRequest{StatusId: first.GetId()})

		_, delerr := api.DeleteStatus(ctx, req)
		require.NoError(t, delerr)
	})
	t.Run("not-found", func(t *testing.T) {
		ctx := context.TODO()
		api, _, httpSrv, err := apiTestSetup(t)
		defer httpSrv.Close()
		require.NoError(t, err)

		req := connect.NewRequest(&statusthingv1.DeleteStatusRequest{StatusId: "not-a-valid-status"})

		delres, delerr := api.DeleteStatus(ctx, req)
		require.ErrorIs(t, delerr, serrors.ErrNotFound)
		require.Nil(t, delres)
	})
}
func TestGetItems(t *testing.T) {
	t.Parallel()

	type createItem struct {
		item *statusthingv1.Item
		opts []filters.FilterOption
	}
	type createStatus struct {
		status *statusthingv1.Status
		opts   []filters.FilterOption
	}
	testCases := map[string]struct {
		createItems    []createItem
		createStatuses []createStatus
		reqOpts        []filters.FilterOption
		reslen         int
		err            error
	}{
		"no-options": {
			reslen: 1,
			createStatuses: []createStatus{
				{
					status: &statusthingv1.Status{
						Name: t.Name(),
						Kind: statusthingv1.StatusKind_STATUS_KIND_AVAILABLE,
					},
					opts: []filters.FilterOption{
						filters.WithStatusID("my-custom-status-id"),
					},
				},
			},
			createItems: []createItem{
				{
					item: &statusthingv1.Item{
						Name: t.Name(),
					},
					opts: []filters.FilterOption{
						filters.WithStatusID("my-custom-status-id"),
					},
				},
			},
		},
		"with-kinds-absent": {
			reslen: 0,
			createStatuses: []createStatus{
				{status: &statusthingv1.Status{
					Name: t.Name() + "_too",
					Kind: statusthingv1.StatusKind_STATUS_KIND_CREATED,
				}},
			},
			createItems: []createItem{
				{item: &statusthingv1.Item{Name: t.Name()}},
			},
			reqOpts: []filters.FilterOption{filters.WithStatusKinds(statusthingv1.StatusKind_STATUS_KIND_DOWN)},
		},
		"with-kinds-present": {
			reslen: 1,
			createStatuses: []createStatus{
				{
					status: &statusthingv1.Status{
						Name: "another-status",
						Kind: statusthingv1.StatusKind_STATUS_KIND_DOWN,
					},
				},
			},
			createItems: []createItem{
				{
					item: &statusthingv1.Item{Name: "1"},
				},
				{
					item: &statusthingv1.Item{Name: "2"},
				},
				{
					item: &statusthingv1.Item{Name: "3"},
					opts: []filters.FilterOption{
						filters.WithStatus(&statusthingv1.Status{Name: "my-custom-status", Kind: statusthingv1.StatusKind_STATUS_KIND_CREATED}),
					},
				},
			},
			reqOpts: []filters.FilterOption{filters.WithStatusKinds(statusthingv1.StatusKind_STATUS_KIND_CREATED)},
		},
		"with-status-ids": {
			reslen: 1,
			createStatuses: []createStatus{
				{
					status: &statusthingv1.Status{
						Name: t.Name() + "_too",
						Kind: statusthingv1.StatusKind_STATUS_KIND_DOWN,
					},
				},
				{
					status: &statusthingv1.Status{
						Name: t.Name() + "_too",
						Kind: statusthingv1.StatusKind_STATUS_KIND_CREATED,
					},
					opts: []filters.FilterOption{filters.WithStatusID("with-status-ids")},
				},
			},
			createItems: []createItem{
				{
					item: &statusthingv1.Item{
						Name: t.Name() + "_item",
					},
					opts: []filters.FilterOption{filters.WithStatusID("with-status-ids")},
				},
				{
					item: &statusthingv1.Item{
						Name: t.Name() + "_item_2",
					},
				},
			},
			reqOpts: []filters.FilterOption{filters.WithStatusIDs("with-status-ids")},
		},
	}
	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			ctx := context.TODO()
			api, _, httpSrv, err := apiTestSetup(t)
			defer httpSrv.Close()
			require.NoError(t, err)

			req := &statusthingv1.ListItemsRequest{}
			f, err := filters.New(tc.reqOpts...)
			require.NoError(t, err, "should parse request options")
			if f.StatusKinds() != nil {
				req.Kinds = f.StatusKinds()
			}
			if f.StatusIDs() != nil {
				req.StatusIds = f.StatusIDs()
			}

			if len(tc.createStatuses) != 0 {
				for idx, s := range tc.createStatuses {
					_, err := api.sts.AddStatus(ctx, fmt.Sprintf("%s-%d", s.status.GetName(), idx), s.status.GetKind(), s.opts...)
					require.NoError(t, err)
				}
			}
			if len(tc.createItems) != 0 {
				for _, i := range tc.createItems {
					_, err := api.sts.AddItem(ctx, i.item.GetName(), i.opts...)
					require.NoError(t, err)
				}
			}
			res, err := api.ListItems(ctx, connect.NewRequest(req))
			if tc.err != nil {
				require.Error(t, err)
				require.Nil(t, res.Msg.GetItems())
			} else {
				require.NoError(t, err)
				require.NotNil(t, res.Msg.GetItems())
				require.Len(t, res.Msg.GetItems(), tc.reslen)
			}
		})
	}
}
func TestAddItem(t *testing.T) {
	t.Parallel()
	type testItem struct {
		itemName string
		opts     []filters.FilterOption
	}
	type testStatus struct {
		name string
		kind statusthingv1.StatusKind
		opts []filters.FilterOption
	}
	type testcase struct {
		items    []testItem
		statuses []testStatus
		reqOpts  []filters.FilterOption
		err      error
	}
	testcases := map[string]testcase{
		"happy-path": {
			items: []testItem{
				{itemName: "happy-path-item"},
			},
		},
		"with-description": {
			items: []testItem{
				{
					itemName: "with-description",
					opts: []filters.FilterOption{
						filters.WithDescription("my-desc"),
					},
				},
			},
		},
		"with-statusid-valid": {
			items: []testItem{
				{
					itemName: "with-statusid-valid",
					opts: []filters.FilterOption{
						filters.WithStatusID("my-test-status"),
					},
				},
			},
			statuses: []testStatus{
				{
					name: "my-test-status",
					kind: statusthingv1.StatusKind_STATUS_KIND_AVAILABLE,
					opts: []filters.FilterOption{filters.WithStatusID("my-test-status")},
				},
			},
		},
		"with-statusid-invalid": {
			items: []testItem{
				{
					itemName: "with-statusid-invalid",
					opts: []filters.FilterOption{
						filters.WithStatusID("my-status-id"),
					},
				},
			},
			err: serrors.ErrNotFound,
		},
		"with-new-status": {
			items: []testItem{
				{
					itemName: "with-new-status",
					opts: []filters.FilterOption{
						filters.WithStatus(&statusthingv1.Status{Name: "my-new-status", Kind: statusthingv1.StatusKind_STATUS_KIND_AVAILABLE}),
					},
				},
			},
		},
		"with-initial-note": {
			items: []testItem{
				{
					itemName: "with-note",
					opts: []filters.FilterOption{
						filters.WithNoteText("my-note-text"),
					},
				},
			},
		},
	}
	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			ctx := context.TODO()
			api, _, httpSrv, err := apiTestSetup(t)
			defer httpSrv.Close()
			require.NoError(t, err)

			// create any statuses if requested
			for _, s := range tc.statuses {
				// we're just going to go to the status service to add these
				_, err = api.sts.AddStatus(ctx, s.name, s.kind, s.opts...)
				require.NoError(t, err)
			}

			for _, i := range tc.items {
				f, err := filters.New(i.opts...)
				require.NoError(t, err)

				req := &statusthingv1.AddItemRequest{
					Name: i.itemName,
				}

				if f.Description() != "" {
					req.Description = f.Description()
				}
				if f.Status() != nil {
					req.InitialStatus = f.Status()
				}
				if f.StatusID() != "" {
					req.InitialStatusId = f.StatusID()
				}
				if f.NoteText() != "" {
					req.InitialNoteText = f.NoteText()
				}
				res, err := api.AddItem(ctx, connect.NewRequest(req))

				if tc.err != nil {
					require.ErrorIs(t, err, tc.err)
					require.Nil(t, res)
				} else {
					resitem := res.Msg.GetItem()
					require.NoError(t, err)
					require.NotNil(t, resitem)
					require.NotEmpty(t, resitem.GetId())
					require.Equal(t, i.itemName, resitem.GetName())
					require.Equal(t, f.Description(), resitem.GetDescription())
					if f.StatusID() != "" {
						require.Equal(t, f.StatusID(), resitem.GetStatus().GetId())
					}
					if f.Status() != nil {
						require.NotNil(t, resitem.GetStatus())
						require.Equal(t, f.Status().GetName(), resitem.GetStatus().GetName())
						require.Equal(t, f.Status().GetKind(), resitem.GetStatus().GetKind())
					}
					if f.NoteText() != "" {
						require.NotNil(t, resitem.GetNotes())
						// we only support 1 note on creation
						require.Len(t, resitem.GetNotes(), 1)
						require.Equal(t, f.NoteText(), resitem.GetNotes()[0].GetText())
					}
				}
			}
		})
	}
}
func TestUpdateItem(t *testing.T) {
	t.Parallel()
	t.Run("happy-path", func(t *testing.T) {
		ctx := context.TODO()
		api, _, httpSrv, err := apiTestSetup(t)
		defer httpSrv.Close()
		require.NoError(t, err)

		// have to add an item first
		res, err := api.AddItem(
			ctx,
			connect.NewRequest(&statusthingv1.AddItemRequest{
				Name:        t.Name(),
				Description: t.Name(),
			}))
		require.NoError(t, err)
		item := res.Msg.GetItem()
		require.NotNil(t, item)
		require.NotEmpty(t, item.GetId())
		require.Equal(t, t.Name(), item.GetName())
		require.Equal(t, t.Name(), item.GetDescription())

		sres, serr := api.AddStatus(ctx, connect.NewRequest(&statusthingv1.AddStatusRequest{
			Name:        t.Name() + "_name",
			Description: t.Name() + "_desc",
			Color:       t.Name() + "_color",
			Kind:        statusthingv1.StatusKind_STATUS_KIND_AVAILABLE,
		}))
		require.NoError(t, serr, "status should create")
		require.NotNil(t, sres.Msg.GetStatus())
		status := sres.Msg.GetStatus()
		_, err2 := api.UpdateItem(ctx, connect.NewRequest(&statusthingv1.UpdateItemRequest{
			ItemId:      item.GetId(),
			Name:        "new_name",
			Description: "new_description",
			StatusId:    status.GetId(),
		}))
		require.NoError(t, err2)

		// now we check the update
		res3, err3 := api.ListItems(ctx, connect.NewRequest(&statusthingv1.ListItemsRequest{}))
		require.NoError(t, err3)
		require.Len(t, res3.Msg.GetItems(), 1)

		res4, err4 := api.GetItem(ctx, connect.NewRequest(&statusthingv1.GetItemRequest{ItemId: res3.Msg.GetItems()[0].GetId()}))
		require.NoError(t, err4)
		require.NotNil(t, res4.Msg.GetItem())
	})

	t.Run("edit-errors", func(t *testing.T) {
		// arrange
		ctx := context.TODO()
		api, _, httpSrv, err := apiTestSetup(t)
		defer httpSrv.Close()
		require.NoError(t, err)

		// have to add an item first
		res, err := api.AddItem(
			ctx,
			connect.NewRequest(&statusthingv1.AddItemRequest{
				Name:        t.Name(),
				Description: t.Name(),
			}))
		require.NoError(t, err)
		item := res.Msg.GetItem()
		require.NotNil(t, item)
		require.NotEmpty(t, item.GetId())

		// act
		_, err2 := api.UpdateItem(ctx, connect.NewRequest(&statusthingv1.UpdateItemRequest{
			ItemId:   item.GetId(),
			StatusId: "invalid-status-id",
		}))

		// assert
		require.ErrorIs(t, err2, serrors.ErrNotFound)
	})
}
func TestDeleteItem(t *testing.T) {
	t.Run("existing", func(t *testing.T) {
		ctx := context.TODO()
		api, _, httpSrv, err := apiTestSetup(t)
		defer httpSrv.Close()
		require.NoError(t, err)

		// have to add an item first
		res, err := api.AddItem(
			context.TODO(),
			connect.NewRequest(&statusthingv1.AddItemRequest{
				Name:        t.Name(),
				Description: t.Name(),
			}))
		require.NoError(t, err)
		require.NotNil(t, res)
		itemID := res.Msg.GetItem().GetId()
		_, delerr := api.DeleteItem(ctx, connect.NewRequest(&statusthingv1.DeleteItemRequest{
			ItemId: itemID,
		}))
		require.NoError(t, delerr)
	})
	t.Run("not-existing", func(t *testing.T) {
		ctx := context.TODO()
		api, _, httpSrv, err := apiTestSetup(t)
		defer httpSrv.Close()
		require.NoError(t, err)

		delres, delerr := api.DeleteItem(ctx, connect.NewRequest(&statusthingv1.DeleteItemRequest{
			ItemId: "not-there",
		}))
		require.ErrorIs(t, delerr, serrors.ErrNotFound)
		require.Nil(t, delres)
	})
}
func TestGetNotes(t *testing.T) {
	t.Parallel()
	t.Run("happy-path", func(t *testing.T) {
		// arrange
		ctx := context.TODO()
		api, _, httpSrv, err := apiTestSetup(t)
		defer httpSrv.Close()
		require.NoError(t, err)

		_, err = api.sts.AddItem(ctx, t.Name(), filters.WithItemID(t.Name()), filters.WithNoteText("yoooooo"))
		require.NoError(t, err)

		// act
		res, err := api.ListNotes(ctx, connect.NewRequest(&statusthingv1.ListNotesRequest{ItemId: t.Name()}))

		// assert
		require.NoError(t, err)
		require.Len(t, res.Msg.GetNotes(), 1)
		require.Equal(t, "yoooooo", res.Msg.GetNotes()[0].GetText())
	})

	t.Run("no-such-item", func(t *testing.T) {
		// arrange
		ctx := context.TODO()
		api, _, httpSrv, err := apiTestSetup(t)
		defer httpSrv.Close()
		require.NoError(t, err)

		// act
		res, err := api.ListNotes(ctx, connect.NewRequest(&statusthingv1.ListNotesRequest{ItemId: t.Name()}))

		// assert
		require.ErrorIs(t, err, serrors.ErrNotFound)
		require.Nil(t, res)
	})
}
func TestAddNote(t *testing.T) {
	t.Parallel()
	t.Run("happy-path", func(t *testing.T) {
		// arrange
		ctx := context.TODO()
		api, _, httpSrv, err := apiTestSetup(t)
		defer httpSrv.Close()
		require.NoError(t, err)

		_, err = api.sts.AddItem(ctx, t.Name(), filters.WithItemID(t.Name()))
		require.NoError(t, err)

		// act
		res, err := api.AddNote(ctx, connect.NewRequest(&statusthingv1.AddNoteRequest{ItemId: t.Name(), NoteText: t.Name()}))

		// assert
		require.NoError(t, err)
		require.NotNil(t, res.Msg.GetNote())
		require.NotEmpty(t, res.Msg.GetNote().GetId())
	})
	t.Run("not-found", func(t *testing.T) {
		// arrange
		ctx := context.TODO()
		api, _, httpSrv, err := apiTestSetup(t)
		defer httpSrv.Close()
		require.NoError(t, err)

		// act
		res, err := api.AddNote(ctx, connect.NewRequest(&statusthingv1.AddNoteRequest{ItemId: "invalid-yall", NoteText: t.Name()}))

		// assert
		require.ErrorIs(t, err, serrors.ErrNotFound)
		require.Nil(t, res)
	})
}
func TestDeleteNote(t *testing.T) {
	t.Parallel()
	t.Run("happy-path", func(t *testing.T) {
		// arrange
		ctx := context.TODO()
		api, _, httpSrv, err := apiTestSetup(t)
		defer httpSrv.Close()
		require.NoError(t, err)

		item, err := api.sts.AddItem(ctx, t.Name())
		require.NoError(t, err, "item create should not err")
		require.NotNil(t, item, "created item should not be nil")
		note, err := api.sts.AddNote(ctx, item.GetId(), t.Name())
		require.NoError(t, err, "note create should not error")
		require.NotNil(t, note, "created note should not be nil")

		// act
		res, err := api.DeleteNote(ctx, connect.NewRequest(&statusthingv1.DeleteNoteRequest{NoteId: note.GetId()}))

		// assert
		require.NoError(t, err)
		require.Nil(t, res.Msg)
	})
	t.Run("not-found", func(t *testing.T) {
		// arrange
		ctx := context.TODO()
		api, _, httpSrv, err := apiTestSetup(t)
		defer httpSrv.Close()
		require.NoError(t, err)

		// act
		res, err := api.DeleteNote(ctx, connect.NewRequest(&statusthingv1.DeleteNoteRequest{NoteId: t.Name()}))

		// assert
		require.ErrorIs(t, err, serrors.ErrNotFound)
		require.Nil(t, res)
	})
}
func TestUpdateNote(t *testing.T) {
	t.Run("happy-path", func(t *testing.T) {
		// arrange
		ctx := context.TODO()
		api, _, httpSrv, err := apiTestSetup(t)
		defer httpSrv.Close()
		require.NoError(t, err)

		item, err := api.sts.AddItem(ctx, t.Name())
		require.NoError(t, err, "create item should not fail")
		require.NotNil(t, item, "created item should not be nil")

		note, err := api.sts.AddNote(ctx, item.GetId(), t.Name())
		require.NoError(t, err, "create note should not error")
		require.NotNil(t, note, "created note should not be nil")

		// act
		res, err := api.UpdateNote(ctx, connect.NewRequest(&statusthingv1.UpdateNoteRequest{NoteId: note.GetId(), NoteText: "new-note-text"}))

		// assert
		require.NoError(t, err)
		require.Nil(t, res.Msg)
	})
	t.Run("not-found", func(t *testing.T) {
		// arrange
		ctx := context.TODO()
		api, _, httpSrv, err := apiTestSetup(t)
		defer httpSrv.Close()
		require.NoError(t, err)

		// act
		res, err := api.UpdateNote(ctx, connect.NewRequest(&statusthingv1.UpdateNoteRequest{NoteId: t.Name(), NoteText: t.Name()}))

		// assert
		require.ErrorIs(t, err, serrors.ErrNotFound)
		require.Nil(t, res)
	})
}
func TestGetNote(t *testing.T) {
	t.Run("happy-path", func(t *testing.T) {
		// arrange
		ctx := context.TODO()
		api, _, httpSrv, err := apiTestSetup(t)
		defer httpSrv.Close()
		require.NoError(t, err)

		item, err := api.sts.AddItem(ctx, t.Name())
		require.NoError(t, err, "create item should not fail")
		require.NotNil(t, item, "created item should not be nil")

		note, err := api.sts.AddNote(ctx, item.GetId(), t.Name())
		require.NoError(t, err, "create note should not error")
		require.NotNil(t, note, "created note should not be nil")

		// act
		res, err := api.GetNote(ctx, connect.NewRequest(&statusthingv1.GetNoteRequest{NoteId: note.GetId()}))

		// assert
		require.NoError(t, err)
		require.NotNil(t, res.Msg.GetNote())
		require.Equal(t, t.Name(), res.Msg.GetNote().GetText())
	})
	t.Run("not-found", func(t *testing.T) {
		// arrange
		ctx := context.TODO()
		api, _, httpSrv, err := apiTestSetup(t)
		defer httpSrv.Close()
		require.NoError(t, err)

		// act
		res, err := api.GetNote(ctx, connect.NewRequest(&statusthingv1.GetNoteRequest{NoteId: t.Name()}))

		// assert
		require.ErrorIs(t, err, serrors.ErrNotFound)
		require.Nil(t, res)
	})
}
func TestGetItem(t *testing.T) {
	t.Parallel()
	t.Run("happy-path", func(t *testing.T) {
		// arrange
		ctx := context.TODO()
		api, _, httpSrv, err := apiTestSetup(t)
		defer httpSrv.Close()
		require.NoError(t, err)

		item, err := api.sts.AddItem(ctx, t.Name(), filters.WithItemID(t.Name()))
		require.NoError(t, err, "create item should not fail")
		require.NotNil(t, item, "created item should not be nil")

		// act
		res, err := api.GetItem(ctx, connect.NewRequest(&statusthingv1.GetItemRequest{
			ItemId: t.Name(),
		}))
		require.NoError(t, err)
		require.NotNil(t, res)
	})
	t.Run("not-found", func(t *testing.T) {
		// arrange
		ctx := context.TODO()
		api, _, httpSrv, err := apiTestSetup(t)
		defer httpSrv.Close()
		require.NoError(t, err)

		// act
		res, err := api.GetItem(ctx, connect.NewRequest(&statusthingv1.GetItemRequest{
			ItemId: t.Name(),
		}))
		require.ErrorIs(t, err, serrors.ErrNotFound)
		require.Nil(t, res)
	})
}
func apiTestSetup(t *testing.T) (*APIHandler, *http.Client, *httptest.Server, error) {
	// test setup
	store, err := memdb.New()
	if err != nil {
		return nil, nil, nil, err
	}
	if store == nil {
		return nil, nil, nil, serrors.ErrNilVal
	}
	sts, err := services.NewStatusThingService(store)
	if err != nil {
		return nil, nil, nil, err
	}
	if sts == nil {
		return nil, nil, nil, serrors.ErrNilVal
	}
	api, err := NewAPIHandler(sts)
	if err != nil {
		return nil, nil, nil, err
	}
	if api == nil {
		return nil, nil, nil, serrors.ErrNilVal
	}
	ispath, ishandler := v1connect.NewItemsServiceHandler(api)
	spath, shandler := v1connect.NewStatusServiceHandler(api)
	npath, nhandler := v1connect.NewNotesServiceHandler(api)

	rtr := chi.NewRouter()
	rtr.Handle(ispath, ishandler)
	rtr.Handle(spath, shandler)
	rtr.Handle(npath, nhandler)

	srv := httptest.NewServer(rtr)
	client := srv.Client()
	return api, client, srv, nil
}
