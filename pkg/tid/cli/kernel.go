package cli

import (
	"github.com/SeerUK/tid/pkg/state"
	"github.com/SeerUK/tid/pkg/util"
	"github.com/SeerUK/tid/pkg/types"
)

// TidKernel is the core "container", and provider of services and information to the `tid` cli
// application.
type TidKernel struct {
	// Backend provides an abstracted, but reasonably low-level interface to the underlying storage.
	Backend state.Backend
	// Factory abstracts the creation of services.
	Factory util.Factory
	// Config has all the configurations specified in the config file
	Config types.TomlConfig
}

// NewTidKernel creates a new TidKernel, with services attached.
func NewTidKernel(backend state.Backend, factory util.Factory, tomlConfig types.TomlConfig) *TidKernel {
	return &TidKernel{
		Backend: backend,
		Factory: factory,
		Config: tomlConfig,
	}
}
