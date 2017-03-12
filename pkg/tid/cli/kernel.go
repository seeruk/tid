package cli

import (
	"github.com/SeerUK/tid/pkg/state"
	"github.com/SeerUK/tid/pkg/tracking"
)

// TidKernel is the core "container", and provider of services and information to the `tid` cli
// application.
type TidKernel struct {
	// Backend provides an abstracted, but reasonably low-level interface to the underlying storage.
	Backend state.Backend
	// TrackingFactory abstracts the creation of tracking-related services.
	TrackingFactory tracking.Factory
}

// NewTidKernel creates a new TidKernel, with services attached.
func NewTidKernel(backend state.Backend, trackingFactory tracking.Factory) *TidKernel {
	return &TidKernel{
		Backend:         backend,
		TrackingFactory: trackingFactory,
	}
}
