// Package testutils contains helpers for use in testing
package testutils

import (
	"testing"

	statusthingv1 "github.com/lusis/statusthing/gen/go/statusthing/v1"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// MakeTimestamps creates a [statusthingv1.Timestamps] based on the current time, optionally populating the deleted value
func MakeTimestamps(deleted bool) *statusthingv1.Timestamps {
	now := timestamppb.Now()
	res := &statusthingv1.Timestamps{
		Created: now,
		Updated: now,
	}
	if deleted {
		res.Deleted = now
	}
	return res
}

// MakeItem makes a valid minimal [statusthingv1.Item] for tests
// uval is generally the name of the current test (t.Name()) if determinism is needed
// but can be any value to use as the base for any string values
func MakeItem(uval string) *statusthingv1.Item {
	return &statusthingv1.Item{
		Id:         uval + "_item_id",
		Name:       uval + "_item_name",
		Timestamps: MakeTimestamps(false),
	}
}

// MakeNote makes a valid minimal [statusthingv1.Note] for tests
// uval is generally the name of the current test (t.Name()) if determinism is needed
// but can be any value to use as the base for any string values
func MakeNote(uval string) *statusthingv1.Note {
	return &statusthingv1.Note{
		Id:         uval + "_note_id",
		Text:       uval + "_note_text",
		Timestamps: MakeTimestamps(false),
	}
}

// MakeStatus makes a valid minimal [statusthingv1.Status] for tests
// uval is generally the name of the current test (t.Name()) if determinism is needed
// but can be any value to use as the base for any string values
func MakeStatus(uval string) *statusthingv1.Status {
	return &statusthingv1.Status{
		Id:         uval + "_status_id",
		Name:       uval + "_status_name",
		Kind:       statusthingv1.StatusKind_STATUS_KIND_AVAILABLE,
		Timestamps: MakeTimestamps(false),
	}
}

// MakeUser makes a valid minimal [statusthingv1.User] for tests
// uval is generally the name of the current test (t.Name()) if determinism is needed
// but can be any value to use as the base for any string values
func MakeUser(uval string) *statusthingv1.User {
	return &statusthingv1.User{
		Id:         uval + "_id",
		Username:   uval + "_username",
		Password:   uval + "_password",
		Timestamps: MakeTimestamps(false),
	}
}

// LogAll is used to log items at the end of a test if desired
func LogAll(t *testing.T, logged map[string]any) {
	for name, l := range logged {
		t.Logf("%s: +%v", name, l)
	}
}

// TimestampsEqual cuts down on error-prone copy/paste when needing to test timestamp equality
func TimestampsEqual(expected *statusthingv1.Timestamps, actual *statusthingv1.Timestamps) bool {
	return (expected.GetCreated().AsTime() == actual.GetCreated().AsTime() &&
		expected.GetUpdated().AsTime() == actual.GetUpdated().AsTime() &&
		expected.GetDeleted().AsTime() == actual.GetDeleted().AsTime())
}
