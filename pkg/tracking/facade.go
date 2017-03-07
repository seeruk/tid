package tracking

import (
	"github.com/SeerUK/tid/pkg/errhandling"
	"github.com/SeerUK/tid/pkg/types"
)

// Facade combines other interfaces to provide more complex functionality related to time tracking.
type Facade struct {
	// A SysGateway for manipulating the overarching tid system.
	sysGateway SysGateway
	// A TimesheetGateway for manipulating timesheet / entry data.
	timesheetGateway TimesheetGateway
}

// NewFacade creates a new instance of Facade.
func NewFacade(sysGateway SysGateway, timesheetGateway TimesheetGateway) *Facade {
	return &Facade{
		sysGateway:       sysGateway,
		timesheetGateway: timesheetGateway,
	}
}

// RemoveEntry attempts to remove an entry from existence.
func (f *Facade) RemoveEntry(entry types.Entry) error {
	// Remove from status, if applicable
	status, err := f.sysGateway.FindOrCreateStatus()
	if err != nil {
		return err
	}

	if status.Entry == entry.Hash {
		status.StopAndClear()
	}

	// Remove from timesheet
	sheet, err := f.timesheetGateway.FindOrCreateTimesheet(entry.Timesheet)
	if err != nil {
		return err
	}

	sheet.RemoveEntry(entry)

	errs := errhandling.NewErrorStack()
	errs.Add(f.sysGateway.PersistStatus(status))
	errs.Add(f.timesheetGateway.PersistTimesheet(sheet))

	if !errs.Empty() {
		return errs.Errors()
	}

	// Remove entry
	return f.timesheetGateway.DeleteEntry(entry)
}
