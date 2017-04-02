package cli

import (
	"github.com/SeerUK/tid/pkg/state"
	"github.com/SeerUK/tid/pkg/util"
)

// TidKernel is the core "container", and provider of services and information to the `tid` cli
// application.
type TidKernel struct {
	// Backend provides an abstracted, but reasonably low-level interface to the underlying storage.
	Backend state.Backend
	// Factory abstracts the creation of services.
	Factory util.Factory
}

// NewTidKernel creates a new TidKernel, with services attached.
func NewTidKernel(backend state.Backend, factory util.Factory) *TidKernel {
	return &TidKernel{
		Backend: backend,
		Factory: factory,
	}
}
