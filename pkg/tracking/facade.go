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

	if err = errs.Errors(); err != nil {
		return entry, err
	}

	return entry, nil
}
