package types

import "github.com/SeerUK/tid/proto"

// WorkspaceIndex represents an index of all workspaces.
type WorkspaceIndex struct {
	// Workspaces is an array of workspace names.
	Workspaces []string
}

// NewWorkspaceIndex creates a new instance of WorkspaceIndex.
func NewWorkspaceIndex() WorkspaceIndex {
	return WorkspaceIndex{}
}

// FromMessage reads a `proto.SysWorkspaceIndex` message into this WorkspaceIndex.
func (i *WorkspaceIndex) FromMessage(message *proto.SysWorkspaceIndex) {
	i.Workspaces = message.Workspaces
}

// ToMessage converts this WorkspaceIndex into a `proto.SysWorkspaceIndex`.
func (i *WorkspaceIndex) ToMessage() *proto.SysWorkspaceIndex {
	return &proto.SysWorkspaceIndex{
		Workspaces: i.Workspaces,
	}
}
