package util

import (
	"time"

	"github.com/SeerUK/tid/pkg/state"
	"github.com/SeerUK/tid/pkg/types"
	"github.com/SeerUK/tid/pkg/xtime"
)

// TimesheetFacade provides a simpler interface for common Timesheet-related tasks.
type TimesheetFacade struct {
	// entryFacade is an EntryFacade used for performing tasks related to entries.
	entryFacade *EntryFacade
	// trGateway is a TimesheetGateway used for accessing timesheet storage.
	trGateway state.TrackingGateway
}

// NewTimesheetFacade creates a new TimesheetFacade instance.
func NewTimesheetFacade(trGateway state.TrackingGateway, entryFacade *EntryFacade) *TimesheetFacade {
	return &TimesheetFacade{
		entryFacade: entryFacade,
		trGateway:   trGateway,
	}
}

// Delete attempts to delete a timesheet at the given date.
func (f *TimesheetFacade) Delete(date time.Time) (types.Timesheet, error) {
	sheet, err := f.trGateway.FindTimesheet(date.Format(xtime.DateFmt))
	if err != nil {
		return sheet, err
	}

	// Handle removing all of the entries on the timesheet. The entryFacade instance takes care of
	// updating things that need to be updated when deleting entries.
	for _, entry := range sheet.Entries {
		f.entryFacade.Delete(entry.Hash)
	}

	f.trGateway.DeleteTimesheet(sheet)

	return sheet, nil
}
