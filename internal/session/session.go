// Package session ...
package session

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"golang.org/x/exp/slog"
)

const (
	// LoggedInKey is the session key for being logged in
	LoggedInKey = "loggedin"
)

// Sessions is the global session manager
var Sessions *scs.SessionManager

// NewSession creates a new session manager
func NewSession() {
	s := scs.New()
	s.Cookie.HttpOnly = true
	s.ErrorFunc = func(w http.ResponseWriter, r *http.Request, err error) {
		slog.Error("error in session", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
	Sessions = s
}
