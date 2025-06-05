package client

import (
	"mals-engine/internal/model"
	"mals-engine/internal/workspace"
	"path/filepath"
)

func (s *Client) FindWorkspaceExact(path string) (*workspace.Workspace, bool) {
	if workspace, exists := s.workspaces[path]; exists {
		return workspace, true
	}
	return nil, false
}

func (s *Client) FindWorkspace(path string) (*workspace.Workspace, bool) {
	currentPath := path

	for {
		if workspace, exists := s.workspaces[currentPath]; exists {
			return workspace, true
		}

		parentPath := filepath.Dir(currentPath)

		if parentPath == currentPath || parentPath == "." {
			break
		}

		currentPath = parentPath
	}

	return nil, false
}

func (s *Client) NewWorkspace(path string, m model.ModelService) (*workspace.Workspace, bool) {
	if workspace, exists := s.FindWorkspaceExact(path); exists {
		return workspace, false
	}

	workspace := workspace.NewWorkspace(path, m)
	s.workspaces[path] = workspace

	return workspace, true
}
