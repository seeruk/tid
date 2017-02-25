package tracking

import (
	"time"

	"github.com/SeerUK/tid/pkg/state"
	"github.com/SeerUK/tid/proto"
)

// KeyStatus is the key for the current tracking status in the store.
const KeyStatus = "status"

// KeyTimeSheetFmt is the date formatting string for timesheet keys in the store.
const KeyTimeSheetFmt = "2006-01-02"

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

// FindStatus attempts to get the current status.
func (g *Gateway) FindStatus() (*Status, error) {
	status := NewStatus(&proto.Status{})

	return status, g.store.Read(KeyStatus, status.Message)
}

// FindTimeSheet looks for a timesheet at the given date.
func (g *Gateway) FindTimeSheet(date time.Time) (proto.TimeSheet, error) {
	var sheet proto.TimeSheet

	return sheet, g.store.Read(date.Format(KeyTimeSheetFmt), &sheet)
}

// FindCurrentTimeSheet will find the timesheet for the current date.
func (g *Gateway) FindCurrentTimeSheet() (proto.TimeSheet, error) {
	return g.FindTimeSheet(time.Now().Local())
}

// PersistStatus persists a given status to the store.
func (g *Gateway) PersistStatus(status *Status) error {
	return g.store.Write(KeyStatus, status.Message)
}

// PersistTimesheet persists a given timesheet to the store.
func (g *Gateway) PersistTimesheet(date time.Time, sheet *proto.TimeSheet) error {
	return g.store.Write(date.Format(KeyTimeSheetFmt), sheet)
}
