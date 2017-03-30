package tracking

import "github.com/SeerUK/tid/pkg/state"

// TimesheetFacade provides a simpler interface for common Timesheet-related tasks.
type TimesheetFacade struct {
	// trGateway is a TimesheetGateway used for accessing timesheet storage.
	trGateway state.TrackingGateway
}

// NewTimesheetFacade creates a new TimesheetFacade instance.
func NewTimesheetFacade(trackingGateway state.TrackingGateway) *TimesheetFacade {
	return &TimesheetFacade{
		trGateway: trackingGateway,
	}
}
