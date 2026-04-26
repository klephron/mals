package execution

import (
	"fmt"
	"mals/internal/middleware/document"
	"mals/internal/middleware/workspace"
	"mals/internal/plane"
	"mals/internal/scope"
	"mals/pkg/config"
	"path"
	"reflect"

	"github.com/puzpuzpuz/xsync/v4"
)

type executionNode struct {
	Step *config.Step
	Then *executionNode
	Else *executionNode
}

type executionResource struct {
	scope *scope.Scope
	spec  config.HandlerLspResourceSpec
}

type ExecutionSetResource struct {
	Name  string
	Scope *scope.Scope
	Spec  config.HandlerLspResourceSpec
}

type ExecutionEnvironment struct {
	plane plane.Plane

	graph *executionNode

	resources  *xsync.Map[string, *executionResource]
	workspaces []*workspace.Workspace
	fileUri    *string
	fileLine   *uint32
	fileChar   *uint32
}

func New(plane plane.Plane) *ExecutionEnvironment {
	return &ExecutionEnvironment{
		plane:      plane,
		graph:      nil,
		resources:  nil,
		workspaces: nil,
		fileUri:    nil,
		fileLine:   nil,
		fileChar:   nil,
	}
}

func (s *ExecutionEnvironment) SetResources(resources []ExecutionSetResource) {
	if s.resources == nil {
		s.resources = xsync.NewMap[string, *executionResource]()
	}
	for _, resource := range resources {
		s.resources.Store(resource.Name, &executionResource{
			scope: resource.Scope,
			spec:  resource.Spec,
		})
	}
}

func (s *ExecutionEnvironment) SetWorkspaces(workspaces []*workspace.Workspace) {
	s.workspaces = workspaces
}

func (s *ExecutionEnvironment) SetFileUri(uri string) {
	s.fileUri = &uri
}

func (s *ExecutionEnvironment) SetFileCursor(line, char uint32) {
	s.fileLine = &line
	s.fileChar = &char
}

func (s *ExecutionEnvironment) ResetResources() {
	s.resources = nil
}

func (s *ExecutionEnvironment) ResetWorkspace() {
	s.workspaces = nil
}

func (s *ExecutionEnvironment) ResetFileUri() {
	s.fileUri = nil
}

func (s *ExecutionEnvironment) ResetFileCursor() {
	s.fileLine = nil
	s.fileChar = nil
}

func (s *ExecutionEnvironment) ResetContext() {
	s.ResetFileCursor()
	s.ResetFileUri()
	s.ResetWorkspace()
	s.ResetResources()
}

func (s *ExecutionEnvironment) get(memory map[string]any, segments ...string) (any, error) {
	if len(segments) == 0 {
		return nil, fmt.Errorf("get: empty path")
	}

	if segments[0] == "file" {
		if len(segments) != 2 {
			return nil, fmt.Errorf("get: file path requires exactly one subkey, got %d segment(s)", len(segments)-1)
		}
		return s.getFile(segments[1])
	}

	var current any = memory
	for i, key := range segments {
		var err error
		current, err = traverseGet(current, key)
		if err != nil {
			return nil, fmt.Errorf("get: segment[%d] %q: %w", i, key, err)
		}
	}
	return current, nil
}

func (s *ExecutionEnvironment) set(memory map[string]any, value any, segments ...string) error {
	if len(segments) == 0 {
		return fmt.Errorf("set: empty path")
	}

	if segments[0] == "file" {
		if len(segments) < 2 {
			return fmt.Errorf("set: file path requires a subkey")
		}
		return fmt.Errorf("set: file.%s is read-only", segments[1])
	}

	if memory == nil {
		return fmt.Errorf("set: memory is nil")
	}

	return traverseSet(reflect.ValueOf(memory), segments, value)
}

func (s *ExecutionEnvironment) getCurrentDocument() (*document.Document, error) {
	if s.fileUri == nil {
		return nil, fmt.Errorf("no active file")
	}
	for _, ws := range s.workspaces {
		doc, err := ws.DocumentGet(*s.fileUri)
		if err == nil {
			return doc, nil
		}
	}
	return nil, fmt.Errorf("document not found: %s", *s.fileUri)
}

func (s *ExecutionEnvironment) getFile(key string) (any, error) {
	if s.fileUri == nil {
		return nil, fmt.Errorf("no active file")
	}
	switch key {
	case "name":
		return path.Base(*s.fileUri), nil

	case "path":
		return *s.fileUri, nil

	case "text":
		doc, err := s.getCurrentDocument()
		if err != nil {
			return nil, fmt.Errorf("file.text: %w", err)
		}
		return doc.Text(), nil

	case "text_before":
		if s.fileLine == nil || s.fileChar == nil {
			return nil, fmt.Errorf("file.text_before: no active cursor")
		}
		doc, err := s.getCurrentDocument()
		if err != nil {
			return nil, fmt.Errorf("file.text_before: %w", err)
		}
		return doc.TextBefore(*s.fileLine, *s.fileChar), nil

	case "text_after":
		if s.fileLine == nil || s.fileChar == nil {
			return nil, fmt.Errorf("file.text_after: no active cursor")
		}
		doc, err := s.getCurrentDocument()
		if err != nil {
			return nil, fmt.Errorf("file.text_after: %w", err)
		}
		return doc.TextAfter(*s.fileLine, *s.fileChar), nil

	default:
		return nil, fmt.Errorf("file: unknown field %q", key)
	}
}
