package tracking

import "github.com/SeerUK/tid/pkg/state"

// TimesheetFacade provides a simpler interface for common Timesheet-related tasks.
type TimesheetFacade struct {
	// tsGateway is a TimesheetGateway used for accessing timesheet storage.
	tsGateway state.TimesheetGateway
}

// NewTimesheetFacade creates a new TimesheetFacade instance.
func NewTimesheetFacade(timesheetGateway state.TimesheetGateway) *TimesheetFacade {
	return &TimesheetFacade{
		tsGateway: timesheetGateway,
	}
}
