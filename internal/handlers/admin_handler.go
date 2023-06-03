package handlers

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/go-chi/chi"

	"github.com/lusis/htmxtools"
	"github.com/lusis/statusthing/assets"
	v1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
	"github.com/lusis/statusthing/internal/filters"
	"github.com/lusis/statusthing/internal/serrors"
	"github.com/lusis/statusthing/internal/services"
	"github.com/lusis/statusthing/internal/session"
	"github.com/lusis/statusthing/internal/templating"
	"github.com/lusis/statusthing/internal/validation"

	"golang.org/x/exp/slog"
)

const (
	contentDivTarget    = "#content"
	contentDivID        = "content"
	loginUIBlock        = "login-ui"
	listItemsBlock      = "list-items-ui"
	listStatusBlock     = "list-status-ui"
	hxCurrentURLHeader  = "hx-current-url"
	hxTriggerHeader     = "hx-trigger"
	hxTriggerNameHeader = "hx-trigger-name"
	hxTargetHeader      = "hx-target"
	hxRedirectHeader    = "hx-redirect"
	hxRequestHeader     = "hx-request"
	hxLocationHeader    = "hx-location"
	hxReplaceURLHeader  = "hx-replace-url"
	defaultUIDir        = "./assets/ui/"
	defaultTemplateDir  = "./assets/templates/"
)

// AdminHandler is the http handler for the admin site
// the admin handler serves static content from assets.UIFs and templates from assets.TemplateFS
type AdminHandler struct {
	sts            *services.StatusThingService
	uiFS           fs.FS
	templateFS     fs.FS
	templates      *template.Template
	funcmap        template.FuncMap
	mux            chi.Router
	templateLoader templating.TemplateLoader
}

// NewAdminHandler returns a new admin handler
func NewAdminHandler(sts *services.StatusThingService, mux chi.Router, reloadable bool) (*AdminHandler, error) {
	funcMap := template.FuncMap{
		"items": func() ([]*v1.Item, error) {
			return sts.FindItems(context.TODO())
		},
		"statuses": func() ([]*v1.Status, error) {
			res, err := sts.FindStatus(context.TODO())
			if err != nil {
				return nil, err
			}
			return res, nil
		},
		"notes": func(itemID string) ([]*v1.Note, error) {
			return sts.FindNotes(context.TODO(), itemID)
		},
		"kinds": func() []string {
			return templating.AllStatusKind
		},
	}

	var uifs fs.FS
	var templatefs fs.FS
	var loader templating.TemplateLoader

	if reloadable {
		uifs = os.DirFS(defaultUIDir)
		templatefs = os.DirFS(defaultTemplateDir)
		l, err := templating.NewReloadingFSTemplateLoader(templatefs, "*.html", funcMap)
		if err != nil {
			return nil, err
		}
		loader = l
	} else {
		ui, err := fs.Sub(assets.UIFs, "ui")
		if err != nil {
			return nil, err
		}
		uifs = ui
		tfs, err := fs.Sub(assets.TemplateFS, "templates")
		if err != nil {
			return nil, err
		}
		templatefs = tfs
		templates, err := template.New("").Funcs(funcMap).ParseFS(templatefs, "*.html")
		if err != nil {
			return nil, err
		}
		loader, err = templating.NewDefaultTemplateLoader(templates)
		if err != nil {
			return nil, err
		}
	}

	return newAdminHandler(sts, uifs, templatefs, loader, mux)
}

func newAdminHandler(sts *services.StatusThingService, uifs fs.FS, templatefs fs.FS, loader templating.TemplateLoader, mux chi.Router) (*AdminHandler, error) {
	if sts == nil {
		return nil, serrors.NewError("sts", serrors.ErrNilVal)
	}
	if uifs == nil {
		return nil, serrors.NewError("uifs", serrors.ErrNilVal)
	}
	if templatefs == nil {
		return nil, serrors.NewError("templatefs", serrors.ErrNilVal)
	}
	if loader == nil {
		return nil, serrors.NewError("loader", serrors.ErrNilVal)
	}

	handler := &AdminHandler{
		sts:            sts,
		uiFS:           uifs,
		templateFS:     templatefs,
		templateLoader: loader,
	}

	ourmux := chi.NewRouter()
	session.NewSession()
	ourmux.Use(session.Sessions.LoadAndSave)
	ourmux.Use(htmxtools.Wrap)
	ourmux.Get("/*", handler.templateHandler(http.FileServer(http.FS(uifs))))
	ourmux.Post("/login", hxonly(handler.login))
	ourmux.Post("/add-status", hxonly(handler.addStatus))
	ourmux.Post("/add-item", hxonly(handler.addItem))
	ourmux.Post("/delete-item", hxonly(handler.deleteItem))
	ourmux.Post("/delete-status", hxonly(handler.deleteStatus))
	ourmux.Post("/edit-item", hxonly(handler.addItem))
	ourmux.Post("/edit-status", hxonly(handler.addStatus))
	ourmux.Route("/statuses", func(r chi.Router) {})
	ourmux.Route("/items", func(r chi.Router) {})
	handler.mux = ourmux
	mux.Mount("/", ourmux)
	return handler, nil
}

func hxonly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		htmxreq := htmxtools.RequestFromContext(r.Context())
		if htmxreq == nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next(w, r)
	}
}
func (ah *AdminHandler) login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		slog.Error("unable to parse form", "error", err)
		return
	}
	u, p := r.Form.Get("username"), r.Form.Get("password")
	if !validation.ValidString(u) || !validation.ValidString(p) {
		w.Header().Add(buildHXLocation(loginUIBlock))
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if u == "admin" && p == "password" {
		session.Sessions.Put(r.Context(), session.LoggedInKey, true)
		if err := session.Sessions.RenewToken(r.Context()); err != nil {
			slog.Error("unable to renew token", "error", err)
			w.Header().Add(buildHXLocation(loginUIBlock))
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.Header().Add(hxRedirectHeader, "/")
		w.WriteHeader(http.StatusOK)
	} else {
		w.Header().Add(buildHXLocation(loginUIBlock))
		w.WriteHeader(http.StatusForbidden)
	}
}

func (ah *AdminHandler) addItem(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		slog.Error("unable to parse form", "error", err)
		return
	}
	vars := r.Form
	name := vars.Get("name")
	statusid := vars.Get("status")
	description := vars.Get("description")
	if !validation.ValidString(name) {
		http.Error(w, "name required", http.StatusFailedDependency)
		return
	}

	opts := []filters.FilterOption{}
	if validation.ValidString(statusid) {
		opts = append(opts, filters.WithStatusID(statusid))
	}
	if validation.ValidString(description) {
		opts = append(opts, filters.WithDescription(description))
	}
	res, err := ah.sts.AddItem(r.Context(), name, opts...)
	if err != nil {
		slog.Error("unable to add item", "error", err)
	} else {
		slog.Info("created item", "item_id", res.GetId())
		w.Header().Add(htmxtools.RedirectResponse.String(), "items.html")
		return
	}
}

func (ah *AdminHandler) addStatus(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		slog.Error("unable to parse form", "error", err)
		return
	}
	vars := r.Form
	statusKind := v1.StatusKind_value[vars.Get("kind")]
	f := []filters.FilterOption{}
	if validation.ValidString(vars.Get("description")) {
		f = append(f, filters.WithDescription(vars.Get("description")))
	}
	if validation.ValidString(vars.Get("color")) {
		f = append(f, filters.WithColor(vars.Get("color")))
	}
	res, err := ah.sts.AddStatus(r.Context(), vars.Get("name"), v1.StatusKind(statusKind), f...)
	if err != nil {
		slog.Error("unable to add status", "error", err)
	} else {
		slog.Info("added status", "status_id", res.GetId())
		w.Header().Add(htmxtools.RedirectResponse.String(), "status.html")
		w.WriteHeader(http.StatusAccepted)
		return
	}
}

func (ah *AdminHandler) deleteItem(w http.ResponseWriter, r *http.Request) {
	htmxreq := htmxtools.RequestFromContext(r.Context())
	if htmxreq == nil {
		return
	}

	id := htmxreq.Trigger
	if !validation.ValidString(id) {
		return
	}

	err := ah.sts.RemoveItem(r.Context(), id)
	if err != nil {
		slog.Error("unable to delete item", "error", err)
	} else {
		slog.Info("deleted item", "item_id", id)
		locresp := &htmxtools.HXLocationResponse{
			Path:   listItemsBlock,
			Target: contentDivTarget,
		}
		resph := &htmxtools.HTMXResponse{
			Location: locresp.String(),
		}
		if err := resph.AddToResponse(w); err != nil {
			slog.Error("error adding headers", "error", err)
			return
		}
		w.Header().Add(htmxtools.RedirectResponse.String(), "items.html")
		w.WriteHeader(http.StatusAccepted)
		return
	}
}

func (ah *AdminHandler) deleteStatus(w http.ResponseWriter, r *http.Request) {
	htmxreq := htmxtools.RequestFromContext(r.Context())
	if htmxreq == nil {
		return
	}

	id := htmxreq.Trigger
	if !validation.ValidString(id) {
		return
	}

	err := ah.sts.RemoveStatus(r.Context(), id)
	if err != nil {
		slog.Error("unable to delete status", "error", err)
	} else {
		slog.Info("deleted status", "status_id", id)
		w.Header().Add(htmxtools.RedirectResponse.String(), "status.html")
		w.WriteHeader(http.StatusAccepted)
		return
	}
}

func (ah *AdminHandler) templateHandler(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// cut down on typos in the most convoluted way....
		sd := siteData{
			ContentDiv: contentDivID,
		}
		// TODO: this is super brittle right now
		path := strings.TrimLeft(r.URL.Path, "/")
		if !validation.ValidString(path) {
			// load index template by default
			path = "index.html"
		}
		t := ah.templateLoader.Lookup(path)
		if t == nil {
			slog.Info("no such template. falling back to fileserver", "template", path)
			next.ServeHTTP(w, r)
		} else {
			loggedIn := session.Sessions.GetBool(r.Context(), session.LoggedInKey)
			sd.LoggedIn = loggedIn
			sd.UserID = session.Sessions.Token(r.Context())
			slog.Info("session", session.LoggedInKey, loggedIn)
			if err := t.Execute(w, sd); err != nil {
				slog.Error("unable to execute template", "error", err)
			}
		}
	}
}

type hxRequest struct {
	Trigger     string
	Target      string
	TriggerName string
	CurrentURL  string
}
type siteData struct {
	LoggedIn   bool
	UserID     string
	ContentDiv string
	HXRequest  hxRequest
}

func buildHXLocation(path string) (string, string) {
	return hxLocationHeader, fmt.Sprintf(`{"path":"%s", "target":"%s"}`, path, contentDivTarget)
}

func buildHXReplace(url string) (string, string) {
	return hxReplaceURLHeader, url
}
