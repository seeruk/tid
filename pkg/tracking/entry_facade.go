package tracking

import (
	"errors"
	"time"

	"github.com/SeerUK/tid/pkg/errhandling"
	"github.com/SeerUK/tid/pkg/state"
	"github.com/SeerUK/tid/pkg/types"
)

// EntryFacade provides a simpler interface for common Entry-related tasks.
type EntryFacade struct {
	// sysGateway is a SysGateway used for accessing system storage.
	sysGateway state.SysGateway
	// tsGateway is a TimesheetGateway used for accessing timesheet storage.
	tsGateway state.TimesheetGateway
}

// NewEntryFacade creates a new EntryFacade instance.
func NewEntryFacade(sysGateway state.SysGateway, timesheetGateway state.TimesheetGateway) *EntryFacade {
	return &EntryFacade{
		sysGateway: sysGateway,
		tsGateway:  timesheetGateway,
	}
}

// Create creates and persists a new entry with the given details.
func (f *EntryFacade) Create(start time.Time, dur time.Duration, note string) (types.Entry, error) {
	entry := types.NewEntry()

	sheet, err := f.tsGateway.FindOrCreateTimesheet(start.Format(types.TimesheetKeyDateFmt))
	if err != nil {
		return entry, err
	}

	entry.Duration = dur
	entry.Note = note
	entry.Timesheet = sheet.Key

	sheet.AppendEntry(entry)

	errs := errhandling.NewErrorStack()
	errs.Add(f.tsGateway.PersistTimesheet(sheet))
	errs.Add(f.tsGateway.PersistEntry(entry))

	if !errs.Empty() {
		return entry, errs.Errors()
	}

	return entry, nil
}

// UpdateDuration updates an entry with the given hash with the given duration.
func (f *EntryFacade) UpdateDuration(hash string, duration time.Duration) (types.Entry, error) {
	entry, err := f.tsGateway.FindEntry(hash)
	if err != nil {
		return entry, err
	}

	if duration < 0 {
		return entry, errors.New("tracking: Duration cannot be less than 0")
	}

	entry.Duration = duration

	return entry, f.tsGateway.PersistEntry(entry)
}

// UpdateDurationByOffset updates an entry with the given hash, offsetting the duration by the given
// offset duration.
func (f *EntryFacade) UpdateDurationByOffset(hash string, offset time.Duration) (types.Entry, error) {
	entry, err := f.tsGateway.FindEntry(hash)
	if err != nil {
		return entry, err
	}

	status, err := f.sysGateway.FindOrCreateStatus()
	if err != nil {
		return entry, err
	}

	if status.IsRunning && status.Entry == entry.Hash {
		entry.UpdateDuration()
	}

	duration := entry.Duration + offset

	return f.UpdateDuration(hash, duration)
}

// UpdateNote updates an entry with the given hash with the given note.
func (f *EntryFacade) UpdateNote(hash string, note string) (types.Entry, error) {
	entry, err := f.tsGateway.FindEntry(hash)
	if err != nil {
		return entry, err
	}

	entry.Note = note

	return entry, f.tsGateway.PersistEntry(entry)
}

// Delete deletes persisted data for a timesheet entry with the given hash.
func (f *EntryFacade) Delete(hash string) (types.Entry, error) {
	entry, err := f.tsGateway.FindEntry(hash)
	if err != nil {
		return entry, err
	}

	// Remove from status, if applicable
	status, err := f.sysGateway.FindOrCreateStatus()
	if err != nil {
		return entry, err
	}

	if status.Entry == entry.Hash {
		status.StopAndClear()
	}

	// Remove from timesheet
	sheet, err := f.tsGateway.FindOrCreateTimesheet(entry.Timesheet)
	if err != nil {
		return entry, err
	}

	sheet.RemoveEntry(entry)

	errs := errhandling.NewErrorStack()
	errs.Add(f.sysGateway.PersistStatus(status))
	errs.Add(f.tsGateway.PersistTimesheet(sheet))

	if !errs.Empty() {
		return entry, errs.Errors()
	}

	err = f.tsGateway.DeleteEntry(entry)
	if err != nil {
		return entry, err
	}

	return entry, nil
}
