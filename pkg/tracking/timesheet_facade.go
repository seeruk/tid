package tracking

// TimesheetFacade provides a simpler interface for common Timesheet-related tasks.
type TimesheetFacade struct {
	// tsGateway is a TimesheetGateway used for accessing timesheet storage.
	tsGateway TimesheetGateway
}

// NewTimesheetFacade creates a new TimesheetFacade instance.
func NewTimesheetFacade(timesheetGateway TimesheetGateway) *TimesheetFacade {
	return &TimesheetFacade{
		tsGateway: timesheetGateway,
	}
}
