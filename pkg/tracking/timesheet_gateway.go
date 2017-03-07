package tracking

import (
	"fmt"
	"time"

	"github.com/SeerUK/tid/pkg/errhandling"
	"github.com/SeerUK/tid/pkg/state"
	"github.com/SeerUK/tid/pkg/types"
	"github.com/SeerUK/tid/proto"
)

const (
	// KeyEntryFmt is the formatting string for the entry keys in the store.
	KeyEntryFmt = "entry:%s"
	// KeyTimesheetFmt is the formatting string for the timesheet keys in the store.
	KeyTimesheetFmt = "sheet:%s"
)

// TimesheetGateway provides access to timesheet data in the database.
type TimesheetGateway interface {
	// FindEntry attempts to find an entry with the given key.
	FindEntry(hash string) (types.Entry, error)
	// FindEntryHashByShortHash attempts to find an entry's hash by it's short hash.
	FindEntryHashByShortHash(hash string) (string, error)
	// FindTimesheet attempts to find a timesheet for the given date.
	FindTimesheet(key string) (types.Timesheet, error)
	// FindOrCreateTimesheet attempts to find a timesheet for the given date, if one is not in the
	// store then a new timesheet object is instantiated.
	FindOrCreateTimesheet(key string) (types.Timesheet, error)
	// FindOrCreateTodaysTimesheet attempts to find the timesheet for the current date, if one is
	// not in the store then a new timesheet object is instantiated.
	FindOrCreateTodaysTimesheet() (types.Timesheet, error)
	// PersistEntry persists a given entry to the store.
	PersistEntry(entry types.Entry) error
	// PersistTimesheet persists a given timesheet to the store.
	PersistTimesheet(timesheet types.Timesheet) error
	// DeleteEntry attempts to delete and entry from the store.
	DeleteEntry(entry types.Entry) error
}

// storeTimesheetGateway is a functional TimesheetGateway.
type storeTimesheetGateway struct {
	// An underlying store to access.
	store state.Store
	// A SysGateway to lookup system info.
	sysGateway SysGateway
}

// NewStoreTimesheetGateway creates a new timesheet gateway.
func NewStoreTimesheetGateway(store state.Store, sysGateway SysGateway) TimesheetGateway {
	return &storeTimesheetGateway{
		store:      store,
		sysGateway: sysGateway,
	}
}

func (g *storeTimesheetGateway) FindEntry(hash string) (types.Entry, error) {
	entry := types.NewEntry()

	status, err := g.sysGateway.FindOrCreateStatus()
	if err != nil {
		return entry, err
	}

	if len(hash) == 7 {
		longHash, err := g.FindEntryHashByShortHash(hash)
		if err != nil {
			return entry, err
		}

		hash = longHash
	}

	message := &proto.TrackingEntry{}

	err = g.store.Read(fmt.Sprintf(KeyEntryFmt, hash), message)
	if err != nil {
		return entry, err
	}

	entry.FromMessage(message)
	entry.IsRunning = status.IsRunning && status.Entry == entry.Hash

	return entry, nil
}

func (g *storeTimesheetGateway) FindEntryHashByShortHash(hash string) (string, error) {
	var ref proto.TrackingEntryRef

	err := g.store.Read(fmt.Sprintf(KeyEntryFmt, hash), &ref)
	if err != nil {
		return "", err
	}

	return ref.Entry, nil
}

func (g *storeTimesheetGateway) FindTimesheet(sheetKey string) (types.Timesheet, error) {
	sheet := types.NewTimesheet()
	message := &proto.TrackingTimesheet{}

	err := g.store.Read(fmt.Sprintf(KeyTimesheetFmt, sheetKey), message)
	if err != nil {
		return sheet, err
	}

	sheet.FromMessage(message)

	return sheet, nil
}

func (g *storeTimesheetGateway) FindOrCreateTimesheet(sheetKey string) (types.Timesheet, error) {
	sheet := types.NewTimesheet()
	sheet.Key = sheetKey

	message := &proto.TrackingTimesheet{}

	err := g.store.Read(fmt.Sprintf(KeyTimesheetFmt, sheetKey), message)
	if err != nil && err != state.ErrNilResult {
		return sheet, err
	}

	if err == nil {
		sheet.FromMessage(message)
	}

	return sheet, nil
}

func (g *storeTimesheetGateway) FindOrCreateTodaysTimesheet() (types.Timesheet, error) {
	return g.FindOrCreateTimesheet(time.Now().Local().Format(types.TimesheetKeyDateFmt))
}

func (g *storeTimesheetGateway) PersistEntry(entry types.Entry) error {
	entryRef := &proto.TrackingEntryRef{
		Key:   entry.ShortHash(),
		Entry: entry.Hash,
	}

	// Persisting an entry is a 2-step process, as we need to also store the short-key so we can
	// look up the long key.
	errs := errhandling.NewErrorStack()
	errs.Add(g.store.Write(fmt.Sprintf(KeyEntryFmt, entry.ShortHash()), entryRef))
	errs.Add(g.store.Write(fmt.Sprintf(KeyEntryFmt, entry.Hash), entry.ToMessage()))

	return errs.Errors()
}

func (g *storeTimesheetGateway) PersistTimesheet(sheet types.Timesheet) error {
	return g.store.Write(fmt.Sprintf(KeyTimesheetFmt, sheet.Key), sheet.ToMessage())
}

func (g *storeTimesheetGateway) DeleteEntry(entry types.Entry) error {
	errs := errhandling.NewErrorStack()
	errs.Add(g.store.Delete(fmt.Sprintf(KeyEntryFmt, entry.ShortHash())))
	errs.Add(g.store.Delete(fmt.Sprintf(KeyEntryFmt, entry.Hash)))

	return errs.Errors()
}
