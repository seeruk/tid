package tracking

import (
	"fmt"
	"time"

	"github.com/SeerUK/tid/pkg/errhandling"
	"github.com/SeerUK/tid/pkg/state"
	"github.com/SeerUK/tid/proto"
)

const (
	// KeyStatus is the key for the current tracking status in the store.
	KeyStatus = "status"
	// KeyEntryFmt is the formatting string for the entry keys in the store.
	KeyEntryFmt = "entry:%s"
	// KeyTimesheetDateFmt is the date formatting string for timesheet keys in the store.
	KeyTimesheetDateFmt = "2006-01-02"
	// KeyTimesheetFmt is the formatting string for the timesheet keys in the store.
	KeyTimesheetFmt = "sheet:%s"
)

// Gateway provides access to timesheet data in the database.
type Gateway struct {
	// The underlying storage to access.
	store state.Store
}

// NewGateway creates a new timesheet gateway.
func NewGateway(store state.Store) Gateway {
	return Gateway{
		store: store,
	}
}

// FindEntry attempts to find an entry with the given key.
func (g *Gateway) FindEntry(hash string) (*Entry, error) {
	entry := NewEntry("")
	key := hash

	if len(hash) == 7 {
		longHash, err := g.FindEntryHashByShortHash(hash)
		if err != nil {
			return nil, err
		}

		key = longHash
	}

	return entry, g.store.Read(fmt.Sprintf(KeyEntryFmt, key), entry.Message)
}

// FindEntryHashByShortHash attempts to find an entry's hash by it's short hash.
func (g *Gateway) FindEntryHashByShortHash(hash string) (string, error) {
	var ref proto.TrackingEntryRef

	err := g.store.Read(fmt.Sprintf(KeyEntryFmt, hash), &ref)
	if err != nil {
		return "", err
	}

	return ref.Entry, nil
}

// FindOrCreateStatus attempts to find the current status.
func (g *Gateway) FindOrCreateStatus() (*Status, error) {
	status := NewStatus()

	err := g.store.Read(KeyStatus, status.Message)
	if err != nil && err != state.ErrNilResult {
		return nil, err
	}

	return status, nil
}

// FindTimesheet attempts to find a timesheet for the given date.
func (g *Gateway) FindTimesheet(sheetKey string) (*Timesheet, error) {
	sheet := NewTimesheet(&proto.TrackingTimesheet{
		Key: sheetKey,
	})

	return sheet, g.store.Read(fmt.Sprintf(KeyTimesheetFmt, sheetKey), sheet.Message)
}

// FindOrCreateTimesheet attempts to find a timesheet for the given date, or returns a timesheet.
func (g *Gateway) FindOrCreateTimesheet(sheetKey string) (*Timesheet, error) {
	sheet := NewTimesheet(&proto.TrackingTimesheet{
		Key: sheetKey,
	})

	err := g.store.Read(fmt.Sprintf(KeyTimesheetFmt, sheetKey), sheet.Message)
	if err != nil && err != state.ErrNilResult {
		return nil, err
	}

	return sheet, nil
}

// FindOrCreateTodaysTimesheet attempts to find the timesheet for the current date.
func (g *Gateway) FindOrCreateTodaysTimesheet() (*Timesheet, error) {
	return g.FindOrCreateTimesheet(time.Now().Local().Format(KeyTimesheetDateFmt))
}

// PersistEntry persists a given entry to the store.
func (g *Gateway) PersistEntry(entry *Entry) error {
	entryRef := &proto.TrackingEntryRef{
		Key:   entry.ShortHash(),
		Entry: entry.Hash(),
	}

	// Persisting an entry is a 2-step process, as we need to also store the short-key so we can
	// look up the long key.
	errs := errhandling.NewErrorStack()
	errs.Add(g.store.Write(fmt.Sprintf(KeyEntryFmt, entry.ShortHash()), entryRef))
	errs.Add(g.store.Write(fmt.Sprintf(KeyEntryFmt, entry.Hash()), entry.Message))

	return errs.Errors()
}

// PersistStatus persists a given status to the store.
func (g *Gateway) PersistStatus(status *Status) error {
	return g.store.Write(KeyStatus, status.Message)
}

// PersistTimesheet persists a given timesheet to the store.
func (g *Gateway) PersistTimesheet(sheet *Timesheet) error {
	return g.store.Write(fmt.Sprintf(KeyTimesheetFmt, sheet.Key()), sheet.Message)
}

// RemoveEntry attempts to remove an entry.
func (g *Gateway) RemoveEntry(entry *Entry) error {
	// Remove from status, if applicable
	status, err := g.FindOrCreateStatus()
	if err != nil {
		return err
	}

	if status.Ref() != nil && status.Ref().Entry == entry.Hash() {
		status.StopAndClear()
	}

	// Remove from timesheet
	sheet, err := g.FindOrCreateTimesheet(entry.Timesheet())
	if err != nil {
		return err
	}

	sheet.RemoveEntry(entry)

	errs := errhandling.NewErrorStack()
	errs.Add(g.PersistStatus(status))
	errs.Add(g.PersistTimesheet(sheet))

	if !errs.Empty() {
		return errs.Errors()
	}

	// Remove entry
	errs = errhandling.NewErrorStack()
	errs.Add(g.store.Delete(fmt.Sprintf(KeyEntryFmt, entry.ShortHash())))
	errs.Add(g.store.Delete(fmt.Sprintf(KeyEntryFmt, entry.Hash())))

	return errs.Errors()
}
