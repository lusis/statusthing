package filters

import (
	"sync"
	"time"

	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
)

// FilterOption is a function type for configuring [Filters]
type FilterOption func(*Filters) error

// Filters represent a set of reusable options that can be used in function signatures
type Filters struct {
	// for synchronization
	l sync.RWMutex
	// itemID is for providing a custom [statusthingv1.Item] id
	itemID *string
	// noteID is for providing a custom [statusthingv1.Note] id
	noteID *string
	// statusID is for providing a custom [statusthingv1.Status] id
	statusID *string
	// color is for setting the color of something
	color *string
	// description is for setting the description of something
	description *string
	// statusKind is for setting the [statusthingv1.StatusKind] of a [statusthingv1.Status]
	statusKind statusthingv1.StatusKind
	// timestamps allows setting custom timestamps
	timestamps *statusthingv1.Timestamps
	// noteText allows providing note text for things like updates
	noteText *string
	// thingIDs stores a slice of [statusthingv1.Item] ids
	thingIDs []string
	// statusIDs stores a slice of [statusthingv1.Item] ids
	statusIDs []string
	// noteIDs stores a slice of [statusthingv1.Note] ids
	noteIDs []string
	// statusKinds stores a slice of [statusthingv1.StatusKind]
	statusKinds []statusthingv1.StatusKind
	// status stores a [statusthingv1.Status]
	status *statusthingv1.Status
	// name stores a custom name value
	name *string
	// userid stores a custom userid value
	userid *string
	// firstname stores a custom firstname value
	firstname *string
	// lastname stores a custom lastname value
	lastname *string
	// emailaddress stores a custom email address
	emailaddress *string
	// lastlogin store the lastlogin
	lastlogin *time.Time
	// avatarURL stores the avatar url
	avatarURL *string
	// avatar stores the avatar image
	avatar []byte
	// stores a new password for a password change
	password *string
}

// New returns a new [Filters] configured with the provided [FilterOption]
// you should never create a [Filters] via struct literal for thread safety reasons. Always use [filters.New]
func New(opts ...FilterOption) (*Filters, error) {
	f := &Filters{l: sync.RWMutex{}}
	f.l.Lock()
	for _, opt := range opts {
		if err := opt(f); err != nil {
			f.l.Unlock()
			return nil, err
		}
	}
	f.l.Unlock()
	return f, nil
}

func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
