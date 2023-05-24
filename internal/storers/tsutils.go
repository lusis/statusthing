package storers

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// Consistent timestamp -> datastore mapping
func tsToInt64(ts *timestamppb.Timestamp) int {
	if ts == nil {
		return 0
	}
	return int(ts.AsTime().UnixNano())
}

// Consistent datastore -> timestamp mapping
func intToTs(i int) *timestamppb.Timestamp {
	// zero int val is a nil for us
	if i == 0 {
		return nil
	}
	return timestamppb.New(time.Unix(0, int64(i)).UTC())
}
