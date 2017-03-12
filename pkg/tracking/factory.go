package tracking

import (
	"fmt"

	"github.com/SeerUK/tid/pkg/state"
)

// Factory abstracts the creation of tracking-related services.
type Factory interface {
	// BuildEntryFacade builds an EntryFacade instance.
	BuildEntryFacade() *EntryFacade
	// BuildSysGateway builds a SysGateway instance.
	BuildSysGateway() SysGateway
	// BuildTimesheetGateway builds a TimesheetGateway instance.
	BuildTimesheetGateway() TimesheetGateway
}

// standardFactory provides a standard, simple, functional implementation of the
// TrackingFactory interface.
type standardFactory struct {
	// backend keeps the reference to the storage backend to re-use.
	backend state.Backend
	// sysGateway keeps the reference to a SysGateway instance to re-use.
	sysGateway SysGateway
	// timesheetGateway keeps the reference to a TimesheetGateway to re-use.
	timesheetGateway TimesheetGateway
}

// NewStandardFactory creates a new Factory instance.
func NewStandardFactory(backend state.Backend) Factory {
	return &standardFactory{
		backend: backend,
	}
}

func (f *standardFactory) BuildEntryFacade() *EntryFacade {
	return NewEntryFacade(f.BuildSysGateway(), f.BuildTimesheetGateway())
}

func (f *standardFactory) BuildSysGateway() SysGateway {
	sysStore := f.getStore(f.backend, state.BackendBucketSys)

	return NewStoreSysGateway(sysStore)
}

func (f *standardFactory) BuildTimesheetGateway() TimesheetGateway {
	sysGateway := f.BuildSysGateway()

	status, err := sysGateway.FindOrCreateStatus()
	if err != nil {
		panic(err)
	}

	tsStore := f.getStore(f.backend, fmt.Sprintf(
		state.BackendBucketWorkspaceFmt,
		status.Workspace,
	))

	return NewStoreTimesheetGateway(tsStore, sysGateway)
}

// getStore gets the application data store.
func (f *standardFactory) getStore(backend state.Backend, bucketName string) state.Store {
	return state.NewBackendStore(backend, bucketName)
}
