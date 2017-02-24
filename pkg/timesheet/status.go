package timesheet

import "github.com/SeerUK/tid/proto"

func IsActive(status proto.Status) bool {
	return status.State == proto.Status_STARTED || status.State == proto.Status_PAUSED
}
