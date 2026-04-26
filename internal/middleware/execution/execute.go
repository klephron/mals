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
	memory := make(map[string]any)

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
			if definition.Condition == nil {
				return nil, fmt.Errorf("if condition is nil")
			}
			condition, err := s.renderBool(*definition.Condition, memory)
			if err != nil {
				return nil, err
			}
			if condition == nil {
				return nil, fmt.Errorf("condition is nil")
			}

			if *condition {
				current = current.Then
			} else {
				current = current.Else
			}
			continue

		case *config.StepFor:
			if definition.Condition == nil {
				current = current.Then
				continue
			}

			condition, err := s.renderBool(*definition.Condition, memory)
			if err != nil {
				return nil, err
			}
			if condition == nil {
				return nil, fmt.Errorf("condition is nil")
			}

			if *condition {
				current = current.Then
			} else {
				current = current.Else
			}
			continue

		case *config.StepReturn:
			if definition.Input == nil {
				return nil, nil
			}

			value, err := s.renderValue(*definition.Input, memory)
			if err != nil {
				return nil, err
			}
			return value, nil

		case *config.StepJsonDumps:
		case *config.StepJsonParse:
		case *config.StepJsonParseCompletion:
		case *config.StepLspCompletion:
		case *config.StepModelSimple:
		case *config.StepModelTemplate:
			// TODO

		default:
			return nil, fmt.Errorf("unexpected config.StepDefinition: %#v", definition)
		}

		// simple sequential steps
		current = current.Then
	}
}
