package tracking

import (
	"errors"

	"github.com/SeerUK/tid/pkg/errhandling"
	"github.com/SeerUK/tid/pkg/state"
	"github.com/SeerUK/tid/pkg/types"
)

// Facade provides a simpler interface for common general tracking-related tasks.
type Facade struct {
	// sysGateway is a SysGateway used for accessing system storage.
	sysGateway state.SysGateway
	// trGateway is a TrackingGateway used for accessing tracking storage.
	trGateway state.TrackingGateway
}

// NewFacade creates a new TrackingFacade instance.
func NewFacade(sysGateway state.SysGateway, trGateway state.TrackingGateway) *Facade {
	return &Facade{
		sysGateway: sysGateway,
		trGateway:  trGateway,
	}
}

// Start a new entry, with the given details.
func (f *Facade) Start(note string) (types.Entry, error) {
	var entry types.Entry

	status, err := f.sysGateway.FindOrCreateStatus()
	if err != nil {
		return entry, err
	}

	if status.IsRunning {
		return entry, errors.New("tracking: Stop your existing timer before starting a new one")
	}

	sheet, err := f.trGateway.FindOrCreateTodaysTimesheet()
	if err != nil {
		return entry, err
	}

	entry = types.NewEntry()
	entry.Note = note
	entry.Timesheet = sheet.Key

	sheet.AppendEntry(entry)

	status.Start(sheet, entry)

	errs := errhandling.NewErrorStack()
	errs.Add(f.sysGateway.PersistStatus(status))
	errs.Add(f.trGateway.PersistEntry(entry))
	errs.Add(f.trGateway.PersistTimesheet(sheet))

	if err = errs.Errors(); err != nil {
		return entry, err
	}

	return entry, nil
}

// Stop the currently active entry.
func (f *Facade) Stop() (types.Entry, error) {
	var entry types.Entry

	status, err := f.sysGateway.FindOrCreateStatus()
	if err != nil {
		return entry, err
	}

	if !status.IsRunning {
		return entry, errors.New("tracking: There is no active timer running")
	}

	entry, err = f.trGateway.FindEntry(status.Entry)
	if err != nil {
		return entry, err
	}

	status.Stop()

	errs := errhandling.NewErrorStack()
	errs.Add(f.sysGateway.PersistStatus(status))
	errs.Add(f.trGateway.PersistEntry(entry))

	if err = errs.Errors(); err != nil {
		return entry, err
	}

	return entry, nil
}

// Resume an entry with the given hash. If an empty hash is given, resume the currently active
// timer. If no timer is active, error.
func (f *Facade) Resume(hash string) (types.Entry, error) {
	var entry types.Entry

	status, err := f.sysGateway.FindOrCreateStatus()
	if err != nil {
		return entry, err
	}

	if hash == "" {
		if status.Entry == "" {
			return entry, errors.New("tracking: No timer to resume")
		}

		hash = status.Entry
	}

	entry, err = f.trGateway.FindEntry(hash)
	if err != nil {
		return entry, err
	}

	sheet, err := f.trGateway.FindOrCreateTimesheet(entry.Timesheet)
	if err != nil {
		return entry, err
	}

	status.Start(sheet, entry)

	errs := errhandling.NewErrorStack()
	errs.Add(f.sysGateway.PersistStatus(status))
	errs.Add(f.trGateway.PersistEntry(entry))
	errs.Add(f.trGateway.PersistTimesheet(sheet))

	if err = errs.Errors(); err != nil {
		return entry, err
	}

	return entry, nil
}
