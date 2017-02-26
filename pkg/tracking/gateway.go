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
	KeyTimesheetDateFmt = "sheet:2006-01-02"
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
func (g *Gateway) FindEntry(entryKey string) (*Entry, error) {
	entry := NewEntry("")

	return entry, g.store.Read(fmt.Sprintf(KeyEntryFmt, entryKey), entry.Message)
}

// FindStatus attempts to find the current status.
func (g *Gateway) FindStatus() (*Status, error) {
	status := NewStatus()

	return status, g.store.Read(KeyStatus, status.Message)
}

// FindTimesheet attempts to find a timesheet with the given date.
func (g *Gateway) FindTimesheet(sheetKey string) (*Timesheet, error) {
	sheet := NewTimesheet(&proto.TrackingTimesheet{
		Key: sheetKey,
	})

	return sheet, g.store.Read(fmt.Sprintf(KeyTimesheetFmt, sheetKey), sheet.Message)
}

// FindTodaysTimesheet attempts to find the timesheet for the current date.
func (g *Gateway) FindTodaysTimesheet() (*Timesheet, error) {
	return g.FindTimesheet(time.Now().Local().Format(KeyTimesheetDateFmt))
}

// PersistEntry persists a given entry to the store.
func (g *Gateway) PersistEntry(entry *Entry) error {
	entryRef := &proto.TrackingEntryRef{
		Key:   entry.ShortKey(),
		Entry: entry.Key(),
	}

	// Persisting an entry is a 2-step process, as we need to also store the short-key so we can
	// look up the long key.
	errs := errhandling.NewErrorStack()
	errs.Add(g.store.Write(fmt.Sprintf(KeyEntryFmt, entry.ShortKey()), entryRef))
	errs.Add(g.store.Write(fmt.Sprintf(KeyEntryFmt, entry.Key()), entry.Message))

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
