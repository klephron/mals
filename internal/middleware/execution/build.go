package execution

import (
	"fmt"
	"mals/pkg/config"
)

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

	case *config.StepFor:
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

	s.graph = start

	return nil
}
