package execution

import (
	"fmt"
	"mals/internal/middleware/workspace"
	"mals/internal/plane"
	"mals/pkg/config"

	"github.com/puzpuzpuz/xsync/v4"
)

type executionNode struct {
	Step *config.Step
	Then *executionNode
	Else *executionNode
}

type ExecutionEnvironment struct {
	plane plane.Plane

	cfg *executionNode

	resources  *xsync.Map[string, *config.HandlerLspResource]
	workspaces []*workspace.Workspace
	fileUri    *string
	memory     map[string]any
}

func New(plane plane.Plane) *ExecutionEnvironment {
	return &ExecutionEnvironment{
		plane:      plane,
		cfg:        nil,
		resources:  nil,
		workspaces: nil,
		fileUri:    nil,
		memory:     make(map[string]any),
	}
}

// recursively builds CFG (cuz execution graph is rather small no optimizations are done)
// returns (start, end pair execution)
//
// * if: Then - then, Else - else
// * while: Then - do, Else - <statement_after>
// * return: Then - nil, Else - nil
// * simple: Then - <statement_after>
func (s *ExecutionEnvironment) buildRecursive(steps []*config.Step, i int) (*executionNode, *executionNode, error) {
	if i >= len(steps) {
		return nil, nil, nil
	}

	curr := steps[i]

	switch cd := curr.Definition.(type) {
	case *config.StepIf:
		execThenStart, execThenEnd, execThenErr := s.buildRecursive(cd.Then, 0)
		execElseStart, execElseEnd, execElseErr := s.buildRecursive(cd.Else, 0)
		if execThenErr != nil {
			return nil, nil, execThenErr
		}
		if execElseErr != nil {
			return nil, nil, execElseErr
		}

		if execThenStart == nil {
			return nil, nil, fmt.Errorf("build step %v then is nil", curr)
		}

		execCurrStart := &executionNode{
			Step: curr,
			Then: execThenStart,
			Else: execElseStart,
		}

		execNextStart, execNextEnd, execNextErr := s.buildRecursive(steps, i+1)
		if execNextErr != nil {
			return nil, nil, execNextErr
		}

		if execThenEnd != nil {
			execThenEnd.Then = execNextStart
		}
		if execElseEnd != nil {
			execElseEnd.Then = execNextStart
		}

		execCurrEnd := execNextEnd
		if execCurrEnd == nil {
			execCurrEnd = execThenEnd // any branch is valid
		}
		return execCurrStart, execCurrEnd, nil

	case *config.StepWhile:
		execDoStart, execDoEnd, execDoErr := s.buildRecursive(cd.Do, 0)

		if execDoErr != nil {
			return nil, nil, execDoErr
		}

		if execDoStart == nil {
			return nil, nil, fmt.Errorf("build step %v do is nil", curr)
		}

		execNextStart, execNextEnd, execNextErr := s.buildRecursive(steps, i+1)
		if execNextErr != nil {
			return nil, nil, execNextErr
		}

		execCurrStart := &executionNode{
			Step: curr,
			Then: execDoStart,
			Else: execNextStart,
		}

		if execDoEnd != nil {
			execDoEnd.Then = execCurrStart
		}

		execCurrEnd := execNextEnd
		if execCurrEnd == nil {
			execCurrEnd = execCurrStart // loop re-entry
		}

		return execCurrStart, execCurrEnd, nil

	case *config.StepReturn:
		execCurrStart := &executionNode{
			Step: curr,
			Then: nil,
			Else: nil,
		}

		return execCurrStart, execCurrStart, nil

	case *config.StepJsonDumps:
	case *config.StepJsonParse:
	case *config.StepJsonParseCompletion:
	case *config.StepLspCompletion:
	case *config.StepModelSimple:
	case *config.StepModelTemplate:
		// simple sequential steps: fall through to linear node construction below
	default:
		return nil, nil, fmt.Errorf("unexpected config.StepDefinition: %#v", cd)
	}

	execNextStart, execNextEnd, execNextErr := s.buildRecursive(steps, i+1)
	if execNextErr != nil {
		return nil, nil, execNextErr
	}

	execCurrStart := &executionNode{
		Step: curr,
		Then: execNextStart,
		Else: nil,
	}

	execCurrEnd := execNextEnd
	if execCurrEnd == nil {
		execCurrEnd = execCurrStart
	}

	return execCurrStart, execCurrEnd, nil
}

func (s *ExecutionEnvironment) Build(steps []*config.Step) error {
	if len(steps) == 0 {
		return fmt.Errorf("execution is nil or empty")
	}

	start, _, err := s.buildRecursive(steps, 0)
	if err != nil {
		return err
	}

	s.cfg = start

	return nil
}

func (s *ExecutionEnvironment) SetResources(resources *xsync.Map[string, *config.HandlerLspResource]) {
	s.resources = resources
}

func (s *ExecutionEnvironment) SetWorkspaces(workspaces []*workspace.Workspace) {
	s.workspaces = workspaces
}

func (s *ExecutionEnvironment) SetFileUri(uri string) {
	s.fileUri = &uri
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

func (s *ExecutionEnvironment) ResetMemory() {
	clear(s.memory)
}

func (s *ExecutionEnvironment) ResetContext() {
	s.ResetMemory()
	s.ResetFileUri()
	s.ResetWorkspace()
	s.ResetResources()
}

func (s *ExecutionEnvironment) get(path ...string) {
	//
}

func (s *ExecutionEnvironment) Execute() (any, error) {
	return nil, nil
}
