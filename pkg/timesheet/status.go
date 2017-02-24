package timesheet

import "github.com/SeerUK/tid/proto"

func ResetStatus(status *proto.Status) {
	status.State = proto.Status_STOPPED
	status.TimeSheetEntry = &proto.TimeSheetEntryRef{}
}
