package mapper

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func FromProtoTime(t *timestamppb.Timestamp) time.Time {
	if t == nil {
		return time.Time{}
	}
	return t.AsTime()
}
