// Package templating contains code related to using go templates
package templating

import (
	"io/fs"
	"text/template"

	"golang.org/x/exp/slog"
)

// AllStatusKind is a reusable list of all our current status kinds in a quick slice form
var AllStatusKind = []string{
	"STATUS_KIND_UP",
	"STATUS_KIND_DOWN",
	"STATUS_KIND_WARNING",
	"STATUS_KIND_AVAILABLE",
	"STATUS_KIND_UNAVAILABLE",
	"STATUS_KIND_INVESTIGATING",
	"STATUS_KIND_OBSERVING",
	"STATUS_KIND_CREATED",
	"STATUS_KIND_ONLINE",
	"STATUS_KIND_OFFLINE",
	"STATUS_KIND_DECOMM",
}

// TemplateLoader is something that can lookup templates
type TemplateLoader interface {
	Lookup(s string) *template.Template
}

// DefaultTemplateLoader is the default template loader
type DefaultTemplateLoader struct {
	template *template.Template
}

// NewDefaultTemplateLoader returns a new template loader with defaults
func NewDefaultTemplateLoader(templates *template.Template) (*DefaultTemplateLoader, error) {
	return &DefaultTemplateLoader{template: templates}, nil
}

// Lookup implements the interface
func (l *DefaultTemplateLoader) Lookup(s string) *template.Template {
	return l.template.Lookup(s)
}

// ReloadingFSTemplateLoader is a template loader that reloads templates on each request from an fs.FS
type ReloadingFSTemplateLoader struct {
	fs      fs.FS
	pattern string
	funcmap template.FuncMap
}

// NewReloadingFSTemplateLoader returns a new reloadable template loader
func NewReloadingFSTemplateLoader(fs fs.FS, pattern string, funcmap template.FuncMap) (*ReloadingFSTemplateLoader, error) {
	return &ReloadingFSTemplateLoader{
		funcmap: funcmap,
		fs:      fs,
		pattern: pattern,
	}, nil
}

// Lookup implements the interface
func (r *ReloadingFSTemplateLoader) Lookup(s string) *template.Template {
	if r.fs == nil {
		return nil
	}
	slog.Info("start reloading templates")
	t, err := template.New("").Funcs(r.funcmap).ParseFS(r.fs, r.pattern)
	if err != nil {
		return nil
	}
	slog.Info("end reloading templates")
	return t.Lookup(s)
}
