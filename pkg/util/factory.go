package util

import (
	"fmt"

	"github.com/SeerUK/tid/pkg/state"
)

// Factory abstracts the creation of services.
type Factory interface {
	// BuildEntryFacade builds an EntryFacade instance.
	BuildEntryFacade() *EntryFacade
	// BuildTimesheetFacade builds an TimesheetFacade instance.
	BuildTimesheetFacade() *TimesheetFacade
	// BuildTrackingFacade builds a TrackingFacade instance.
	BuildTrackingFacade() *TrackingFacade
	// BuildSysGateway builds a SysGateway instance.
	BuildSysGateway() state.SysGateway
	// BuildTrackingGateway builds a TimesheetGateway instance.
	BuildTrackingGateway() state.TrackingGateway
}

// standardFactory provides a standard, simple, functional implementation of the
// Factory interface.
type standardFactory struct {
	// backend keeps the reference to the storage backend to re-use.
	backend state.Backend
	// sysGateway keeps the reference to a SysGateway instance to re-use.
	sysGateway state.SysGateway
	// trackingGateway keeps the reference to a TimesheetGateway to re-use.
	trackingGateway state.TrackingGateway
}

// NewStandardFactory creates a new Factory instance.
func NewStandardFactory(backend state.Backend) Factory {
	return &standardFactory{
		backend: backend,
	}
}

func (f *standardFactory) BuildEntryFacade() *EntryFacade {
	return NewEntryFacade(f.BuildSysGateway(), f.BuildTrackingGateway())
}

func (f *standardFactory) BuildTimesheetFacade() *TimesheetFacade {
	return NewTimesheetFacade(f.BuildTrackingGateway(), f.BuildEntryFacade())
}

func (f *standardFactory) BuildTrackingFacade() *TrackingFacade {
	return NewTrackingFacade(f.BuildSysGateway(), f.BuildTrackingGateway())
}

func (f *standardFactory) BuildSysGateway() state.SysGateway {
	sysStore := f.getStore(f.backend, state.BackendBucketSys)

	return state.NewStoreSysGateway(sysStore)
}

func (f *standardFactory) BuildTrackingGateway() state.TrackingGateway {
	sysGateway := f.BuildSysGateway()

	status, err := sysGateway.FindOrCreateStatus()
	if err != nil {
		panic(err)
	}

	tsStore := f.getStore(f.backend, fmt.Sprintf(
		state.BackendBucketWorkspaceFmt,
		status.Workspace,
	))

	return state.NewStoreTrackingGateway(tsStore, sysGateway)
}

// getStore gets the application data store.
func (f *standardFactory) getStore(backend state.Backend, bucketName string) state.Store {
	return state.NewBackendStore(backend, bucketName)
}
