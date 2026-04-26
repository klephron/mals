package execution

import (
	"fmt"
	"mals/pkg/config"
)

func (s *ExecutionEnvironment) Execute() (any, error) {
	if s.graph == nil {
		return nil, fmt.Errorf("nothing to be started")
	}

	return s.execute()
}

func (s *ExecutionEnvironment) execute() (any, error) {
	current := s.graph

	for {
		if current == nil {
			return nil, nil
		}

		if current.Step == nil {
			return nil, fmt.Errorf("step is nil")
		}

		switch definition := current.Step.Definition.(type) {
		case nil:
			return nil, fmt.Errorf("step definition is nil")
		case *config.StepIf:
		case *config.StepJsonDumps:
		case *config.StepJsonParse:
		case *config.StepJsonParseCompletion:
		case *config.StepLspCompletion:
		case *config.StepModelSimple:
		case *config.StepModelTemplate:
		case *config.StepReturn:
		case *config.StepWhile:
		default:
			return nil, fmt.Errorf("unexpected config.StepDefinition: %#v", definition)
		}

	}
}
