package cli

import (
	"fmt"

	"github.com/SeerUK/tid/pkg/state"
	"github.com/SeerUK/tid/pkg/tracking"
)

// TidKernel is the core "container", and provider of services and information to the `tid` cli
// application.
//
// @todo: Set up to use factories. Factories should be used to instantiate the kernel, and then they
// could be swapped out to provide different implementations of interfaces, and offer a more lazy
// loaded style core.
type TidKernel struct {
	// Backend provides an abstracted, but reasonably low-level interface to the underlying storage.
	Backend state.Backend
	// -- Tracking
	// Facade provides useful methods to accomplish complex tasks with a simpler interface.
	Facade *tracking.Facade
	// SysGateway is a gateway to core system information.
	SysGateway tracking.SysGateway
	// TimesheetGateway is a gateway to timesheet information.
	TimesheetGateway tracking.TimesheetGateway
}

// NewTidKernel creates a new TidKernel, with services attached.
func NewTidKernel(backend state.Backend, trackingFactory TrackingFactory) *TidKernel {
	return &TidKernel{
		Backend:          backend,
		Facade:           trackingFactory.BuildFacade(),
		SysGateway:       trackingFactory.BuildSysGateway(),
		TimesheetGateway: trackingFactory.BuildTimesheetGateway(),
	}
}

// TrackingFactory abstracts the creation of tracking-related services.
type TrackingFactory interface {
	// BuildFacade builds a tracking.Facade instance.
	BuildFacade() *tracking.Facade
	// BuildSysGateway builds a tracking.SysGateway instance.
	BuildSysGateway() tracking.SysGateway
	// BuildTimesheetGateway builds a tracking.TimesheetGateway instance.
	BuildTimesheetGateway() tracking.TimesheetGateway
}

// standardTrackingFactory provides a standard, simple, functional implementation of the
// TrackingFactory interface.
//
// @todo: Move me.
type standardTrackingFactory struct {
	// backend keeps the reference to the storage backend to re-use.
	backend state.Backend
	// facade keeps the reference to a tracking.Facade instance to re-use.
	facade tracking.Facade
	// sysGateway keeps the reference to a tracking.SysGateway instance to re-use.
	sysGateway tracking.SysGateway
	// timesheetGateway keeps the reference to a tracking.TimesheetGateway to re-use.
	timesheetGateway tracking.TimesheetGateway
}

// NewStandardTrackingFactory creates a new TrackingFactory instance.
func NewStandardTrackingFactory(backend state.Backend) TrackingFactory {
	return &standardTrackingFactory{
		backend: backend,
	}
}

func (f *standardTrackingFactory) BuildFacade() *tracking.Facade {
	return tracking.NewFacade(f.BuildSysGateway(), f.BuildTimesheetGateway())
}

func (f *standardTrackingFactory) BuildSysGateway() tracking.SysGateway {
	sysStore := f.getStore(f.backend, state.BackendBucketSys)

	return tracking.NewStoreSysGateway(sysStore)
}

func (f *standardTrackingFactory) BuildTimesheetGateway() tracking.TimesheetGateway {
	sysGateway := f.BuildSysGateway()

	status, err := sysGateway.FindOrCreateStatus()
	if err != nil {
		panic(err)
	}

	tsStore := f.getStore(f.backend, fmt.Sprintf(
		state.BackendBucketWorkspaceFmt,
		status.Workspace,
	))

	return tracking.NewStoreTimesheetGateway(tsStore, sysGateway)
}

// getStore gets the application data store.
func (f *standardTrackingFactory) getStore(backend state.Backend, bucketName string) state.Store {
	return state.NewBackendStore(backend, bucketName)
}
