package state

import (
	"log"
	"path/filepath"
)

type State struct {
	logger     *log.Logger
	Workspaces map[string]*Workspace // path should be cleaned
}

func NewState(logger *log.Logger) *State {
	return &State{
		logger: logger,
		Workspaces: make(map[string]*Workspace),
	}
}

func (s *State) FindWorkspaceExact(path string) (*Workspace, bool) {
	if workspace, exists := s.Workspaces[path]; exists {
		return workspace, true
	}
	return nil, false
}

func (s *State) FindWorkspace(path string) (*Workspace, bool) {
	currentPath := path

	for {
		if workspace, exists := s.Workspaces[currentPath]; exists {
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

func (s *State) NewWorkspace(path string) (*Workspace, bool) {
	if workspace, exists := s.FindWorkspaceExact(path); exists {
		return workspace, false
	}

	workspace := NewWorkspace(path)
	s.Workspaces[path] = workspace

	return workspace, true
}
