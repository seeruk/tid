package tracking

import (
	"github.com/SeerUK/tid/proto"
)

// Status represents the status of what is being tracked currently.
type Status struct {
	// Whether or a timer is active.
	IsActive bool
	// The date of the timesheet currently being tracked.
	Timesheet string
	// The hash of the entry currently being tracked.
	Entry string
}

// NewStatus create a new instance of Status.
func NewStatus() Status {
	// @todo: Do we even need this really?
	return Status{}
}

// FromMessage reads a `proto.TrackingStatus` message into this Status.
func (s *Status) FromMessage(message *proto.TrackingStatus) {
	s.IsActive = message.IsActive
	s.Timesheet = message.Timesheet
	s.Entry = message.Entry
}

// ToMessage converts this Status into a `proto.TrackingStatus`.
func (s *Status) ToMessage() *proto.TrackingStatus {
	return &proto.TrackingStatus{
		IsActive:  s.IsActive,
		Timesheet: s.Timesheet,
		Entry:     s.Entry,
	}
}

// Start updates the status to reflect that a given timesheet and entry are being tracked.
func (s *Status) Start(sheet *Timesheet, entry Entry) {
	s.IsActive = true
	s.Timesheet = sheet.Key()
	s.Entry = entry.Hash
}

// Stop updates the status to reflect that tracking has ended (at least temporarily).
func (s *Status) Stop() {
	s.IsActive = false
}

// StopAndClear updates the status to reflect that tracking has ended, and we should no longer know
// about a timesheet or entry.
func (s *Status) StopAndClear() {
	s.Stop()
	s.Timesheet = ""
	s.Entry = ""
}
