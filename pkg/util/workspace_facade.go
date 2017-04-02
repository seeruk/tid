package util

import (
	"fmt"

	"github.com/SeerUK/tid/pkg/state"
)

// WorkspaceFacade provides a simpler interface for common general workspace-related tasks.
type WorkspaceFacade struct {
	// backend is a lower-level backend storage interface.
	backend state.Backend
	// sysGateway is a SysGateway used for accessing system storage.
	sysGateway state.SysGateway
}

// NewWorkspaceFacade creates a new WorkspaceFacade instance.
func NewWorkspaceFacade(backend state.Backend, sysGateway state.SysGateway) *WorkspaceFacade {
	return &WorkspaceFacade{
		backend:    backend,
		sysGateway: sysGateway,
	}
}

// Create attempts to create a new workspace.
func (f *WorkspaceFacade) Create(name string) error {
	index, err := f.sysGateway.FindWorkspaceIndex()
	if err != nil {
		return err
	}

	for _, workspace := range index.Workspaces {
		if workspace == name {
			return fmt.Errorf("util: Workspace '%s' already exists", workspace)
		}
	}

	bucketName := fmt.Sprintf(
		state.BackendBucketWorkspaceFmt,
		name,
	)

	err = f.backend.CreateBucketIfNotExists(bucketName)
	if err != nil {
		return err
	}

	index.Workspaces = append(index.Workspaces, name)

	return f.sysGateway.PersistWorkspaceIndex(index)
}

// Delete attempts to delete a workspace.
func (f *WorkspaceFacade) Delete(name string) error {
	index, err := f.sysGateway.FindWorkspaceIndex()
	if err != nil {
		return err
	}

	exists := false

	// Remove the workspace from the index.
	for i, ws := range index.Workspaces {
		if ws == name {
			index.Workspaces = append(index.Workspaces[:i], index.Workspaces[i+1:]...)
			exists = true
			break
		}
	}

	if !exists {
		return fmt.Errorf("util: Workspace '%s' does not exist", name)
	}

	err = f.sysGateway.PersistWorkspaceIndex(index)
	if err != nil {
		return err
	}

	return f.backend.DeleteBucket(fmt.Sprintf(
		state.BackendBucketWorkspaceFmt,
		name,
	))
}
