package state

import (
	"errors"
	"fmt"
	"time"

	"github.com/SeerUK/tid/pkg/errhandling"
	"github.com/SeerUK/tid/pkg/types"
	"github.com/SeerUK/tid/proto"
)

const (
	// KeyEntryFmt is the formatting string for the entry keys in the store.
	KeyEntryFmt = "entry:%s"
	// KeyTimesheetFmt is the formatting string for the timesheet keys in the store.
	KeyTimesheetFmt = "sheet:%s"
)

// @todo: Consider splitting this up, we have facades for common tasks more granularly than this.

// TrackingGateway provides access to timesheet data in the database.
type TrackingGateway interface {
	// FindEntry attempts to find an entry with the given key.
	FindEntry(hash string) (types.Entry, error)
	// FindEntryHashByShortHash attempts to find an entry's hash by it's short hash.
	FindEntryHashByShortHash(hash string) (string, error)
	// FindEntriesInDateRange attempts to find all of the entries within a given start and end date.
	FindEntriesInDateRange(start time.Time, end time.Time) ([]types.Entry, error)
	// FindTimesheet attempts to find a timesheet for the given date.
	FindTimesheet(key string) (types.Timesheet, error)
	// FindOrCreateTimesheet attempts to find a timesheet for the given date, if one is not in the
	// store then a new timesheet object is instantiated.
	FindOrCreateTimesheet(key string) (types.Timesheet, error)
	// FindOrCreateTodaysTimesheet attempts to find the timesheet for the current date, if one is
	// not in the store then a new timesheet object is instantiated.
	FindOrCreateTodaysTimesheet() (types.Timesheet, error)
	// FindTimesheetsInDateRange attempts to find all timesheets within a given start and end date.
	FindTimesheetsInDateRange(start time.Time, end time.Time) ([]types.Timesheet, error)
	// PersistEntry persists a given entry to the store.
	PersistEntry(entry types.Entry) error
	// PersistTimesheet persists a given timesheet to the store.
	PersistTimesheet(timesheet types.Timesheet) error
	// DeleteEntry attempts to delete and entry from the store.
	DeleteEntry(entry types.Entry) error
}

// storeTimesheetGateway is a functional TimesheetGateway.
type storeTrackingGateway struct {
	// An underlying store to access.
	store Store
	// A SysGateway to lookup system info.
	sysGateway SysGateway
}

// NewStoreTrackingGateway creates a new timesheet gateway.
func NewStoreTrackingGateway(store Store, sysGateway SysGateway) TrackingGateway {
	return &storeTrackingGateway{
		store:      store,
		sysGateway: sysGateway,
	}
}

func (g *storeTrackingGateway) FindEntry(hash string) (types.Entry, error) {
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

	if entry.IsRunning {
		entry.UpdateDuration()
	}

	return entry, nil
}

func (g *storeTrackingGateway) FindEntryHashByShortHash(hash string) (string, error) {
	var ref proto.TrackingEntryRef

	err := g.store.Read(fmt.Sprintf(KeyEntryFmt, hash), &ref)
	if err != nil {
		return "", err
	}

	return ref.Entry, nil
}

func (g *storeTrackingGateway) FindEntriesInDateRange(start time.Time, end time.Time) ([]types.Entry, error) {
	var entries []types.Entry

	timesheets, err := g.FindTimesheetsInDateRange(start, end)
	if err != nil {
		return entries, err
	}

	for _, sheet := range timesheets {
		for _, hash := range sheet.Entries {
			entry, err := g.FindEntry(hash)
			if err != nil {
				return entries, err
			}

			entries = append(entries, entry)
		}
	}

	return entries, nil
}

func (g *storeTrackingGateway) FindTimesheet(sheetKey string) (types.Timesheet, error) {
	sheet := types.NewTimesheet()
	message := &proto.TrackingTimesheet{}

	err := g.store.Read(fmt.Sprintf(KeyTimesheetFmt, sheetKey), message)
	if err != nil {
		return sheet, err
	}

	sheet.FromMessage(message)

	return sheet, nil
}

func (g *storeTrackingGateway) FindOrCreateTimesheet(sheetKey string) (types.Timesheet, error) {
	sheet := types.NewTimesheet()
	sheet.Key = sheetKey

	message := &proto.TrackingTimesheet{}

	err := g.store.Read(fmt.Sprintf(KeyTimesheetFmt, sheetKey), message)
	if err != nil && err != ErrStoreNilResult {
		return sheet, err
	}

	if err == nil {
		sheet.FromMessage(message)
	}

	return sheet, nil
}

func (g *storeTrackingGateway) FindOrCreateTodaysTimesheet() (types.Timesheet, error) {
	return g.FindOrCreateTimesheet(time.Now().Local().Format(types.TimesheetKeyDateFmt))
}

func (g *storeTrackingGateway) FindTimesheetsInDateRange(start time.Time, end time.Time) ([]types.Timesheet, error) {
	var sheets []types.Timesheet

	if start.After(end) {
		return sheets, errors.New("tracking: The start date must be before the end date")
	}

	for current := start; !current.After(end); current = current.AddDate(0, 0, 1) {
		key := current.Format(types.TimesheetKeyDateFmt)

		sheet, err := g.FindTimesheet(key)
		if err != nil && err != ErrStoreNilResult {
			return sheets, err
		}

		if err == ErrStoreNilResult {
			continue
		}

		sheets = append(sheets, sheet)
	}

	return sheets, nil
}

func (g *storeTrackingGateway) PersistEntry(entry types.Entry) error {
	entryRef := &proto.TrackingEntryRef{
		Key:   entry.ShortHash(),
		Entry: entry.Hash,
	}

	// Every time we do anything to an entry, we should update when it was updated. This helps keep
	// things properly in sync.
	entry.Updated = time.Now()

	// Persisting an entry is a 2-step process, as we need to also store the short-key so we can
	// look up the long key.
	errs := errhandling.NewErrorStack()
	errs.Add(g.store.Write(fmt.Sprintf(KeyEntryFmt, entry.ShortHash()), entryRef))
	errs.Add(g.store.Write(fmt.Sprintf(KeyEntryFmt, entry.Hash), entry.ToMessage()))

	return errs.Errors()
}

func (g *storeTrackingGateway) PersistTimesheet(sheet types.Timesheet) error {
	return g.store.Write(fmt.Sprintf(KeyTimesheetFmt, sheet.Key), sheet.ToMessage())
}

func (g *storeTrackingGateway) DeleteEntry(entry types.Entry) error {
	errs := errhandling.NewErrorStack()
	errs.Add(g.store.Delete(fmt.Sprintf(KeyEntryFmt, entry.ShortHash())))
	errs.Add(g.store.Delete(fmt.Sprintf(KeyEntryFmt, entry.Hash)))

	return errs.Errors()
}
