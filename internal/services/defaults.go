package services

import (
	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"

	"github.com/segmentio/ksuid"
)

const (
	defaultColor = green
	green        = "#5DFC0A"
	red          = "#FF0000"
	yellow       = "#EEEB8D"
)

var defaultStatuses = []*statusthingv1.Status{
	// up status
	{Id: func() string { return ksuid.New().String() }(), Name: "UP", Kind: statusthingv1.StatusKind_STATUS_KIND_UP, Color: green},
	{Id: func() string { return ksuid.New().String() }(), Name: "DOWN", Kind: statusthingv1.StatusKind_STATUS_KIND_DOWN, Color: red},
	{Id: func() string { return ksuid.New().String() }(), Name: "WARNING", Kind: statusthingv1.StatusKind_STATUS_KIND_WARNING, Color: yellow},
}
