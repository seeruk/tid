package types

import (
	"github.com/SeerUK/tid/proto"
)

const StatusDefaultWorkspace = "default"

// Status represents the status of what is being tracked currently.
type Status struct {
	// Whether or not a (any) entry's timer is running.
	IsRunning bool
	// The date of the timesheet currently being tracked.
	Timesheet string
	// The hash of the entry currently being tracked.
	Entry string
	// THe name of the workspace currently being tracked.
	Workspace string
}

// NewStatus create a new instance of Status.
func NewStatus() Status {
	// @todo: Do we even need this really?
	return Status{}
}

// FromMessage reads a `proto.SysStatus` message into this Status.
func (s *Status) FromMessage(message *proto.SysStatus) {
	s.IsRunning = message.IsRunning
	s.Timesheet = message.Timesheet
	s.Entry = message.Entry

	if message.Workspace != "" {
		s.Workspace = message.Workspace
	} else {
		s.Workspace = StatusDefaultWorkspace
	}
}

// ToMessage converts this Status into a `proto.SysStatus`.
func (s *Status) ToMessage() *proto.SysStatus {
	return &proto.SysStatus{
		IsRunning: s.IsRunning,
		Timesheet: s.Timesheet,
		Entry:     s.Entry,
		Workspace: s.Workspace,
	}
}

// Start updates the status to reflect that a given timesheet and entry are being tracked.
func (s *Status) Start(sheet Timesheet, entry Entry) {
	s.IsRunning = true
	s.Timesheet = sheet.Key
	s.Entry = entry.Hash
}

// Stop updates the status to reflect that tracking has ended (at least temporarily).
func (s *Status) Stop() {
	s.IsRunning = false
}

// StopAndClear updates the status to reflect that tracking has ended, and we should no longer know
// about a timesheet or entry.
func (s *Status) StopAndClear() {
	s.Stop()
	s.Timesheet = ""
	s.Entry = ""
}
