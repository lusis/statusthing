package storers

import (
	"github.com/segmentio/ksuid"

	v1 "github.com/lusis/statusthing/gen/go/statusthing/v1"
)

const (
	green  = "#5DFC0A"
	red    = "#FF0000"
	yellow = "#EEEB8D"
)

// DefaultStatuses are safe default values that can be populated
var DefaultStatuses = []*v1.Status{
	// up status
	{Id: func() string { return ksuid.New().String() }(), Name: "UP", Kind: v1.StatusKind_STATUS_KIND_UP, Color: green},
	{Id: func() string { return ksuid.New().String() }(), Name: "DOWN", Kind: v1.StatusKind_STATUS_KIND_DOWN, Color: red},
	{Id: func() string { return ksuid.New().String() }(), Name: "WARNING", Kind: v1.StatusKind_STATUS_KIND_WARNING, Color: yellow},
	{Id: func() string { return ksuid.New().String() }(), Name: "CREATED", Kind: v1.StatusKind_STATUS_KIND_CREATED, Color: green},
	{Id: func() string { return ksuid.New().String() }(), Name: "OFFLINE", Kind: v1.StatusKind_STATUS_KIND_OFFLINE, Color: red},
	{Id: func() string { return ksuid.New().String() }(), Name: "ONLINE", Kind: v1.StatusKind_STATUS_KIND_ONLINE, Color: green},
	{Id: func() string { return ksuid.New().String() }(), Name: "OBSERVING", Kind: v1.StatusKind_STATUS_KIND_OBSERVING, Color: yellow},
	{Id: func() string { return ksuid.New().String() }(), Name: "INVESTIGATING", Kind: v1.StatusKind_STATUS_KIND_INVESTIGATING, Color: yellow},
}
