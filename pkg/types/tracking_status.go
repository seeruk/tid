package types

import (
	"github.com/SeerUK/tid/proto"
)

// TrackingStatusDefaultWorkspace is the default workspace to fallback on if one isn't present in
// the DB, or if the currently active workspace is deleted.
const TrackingStatusDefaultWorkspace = "default"

// TrackingStatus represents the status of what is being tracked currently.
type TrackingStatus struct {
	// Whether or not a (any) entry's timer is running.
	IsRunning bool
	// The date of the timesheet currently being tracked.
	Timesheet string
	// The hash of the entry currently being tracked.
	Entry string
	// The name of the workspace currently being tracked.
	Workspace string
}

// NewTrackingStatus creates a new instance of TrackingStatus.
func NewTrackingStatus() TrackingStatus {
	return TrackingStatus{}
}

// FromMessage reads a `proto.SysTrackingStatus` message into this TrackingStatus.
func (s *TrackingStatus) FromMessage(message *proto.SysTrackingStatus) {
	s.IsRunning = message.IsRunning
	s.Timesheet = message.Timesheet
	s.Entry = message.Entry

	if message.Workspace != "" {
		s.Workspace = message.Workspace
	} else {
		s.Workspace = TrackingStatusDefaultWorkspace
	}
}

// ToMessage converts this TrackingStatus into a `proto.SysTrackingStatus`.
func (s *TrackingStatus) ToMessage() *proto.SysTrackingStatus {
	return &proto.SysTrackingStatus{
		IsRunning: s.IsRunning,
		Timesheet: s.Timesheet,
		Entry:     s.Entry,
		Workspace: s.Workspace,
	}
}

// Start updates the status to reflect that a given timesheet and entry are being tracked.
func (s *TrackingStatus) Start(sheet Timesheet, entry Entry) {
	s.IsRunning = true
	s.Timesheet = sheet.Key
	s.Entry = entry.Hash
}

// Stop updates the status to reflect that tracking has ended (at least temporarily).
func (s *TrackingStatus) Stop() {
	s.IsRunning = false
}

// StopAndClear updates the status to reflect that tracking has ended, and we should no longer know
// about a timesheet or entry.
func (s *TrackingStatus) StopAndClear() {
	s.Stop()
	s.Timesheet = ""
	s.Entry = ""
}
