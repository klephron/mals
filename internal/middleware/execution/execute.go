package execution

import (
	"encoding/json"
	"fmt"
	"mals/internal/lsp/protocol"
	"mals/pkg/config"
	"mals/pkg/info"
	"strings"
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

		var assignValue any

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

		case *config.StepLspCompletion:
			if definition.Resource == nil {
				return nil, fmt.Errorf("lsp resource not defined")
			}

			name := *definition.Resource
			resource, ok := s.resources.Load(name)
			if !ok {
				return nil, fmt.Errorf("lsp resource %v not found", name)
			}

			_, ok = resource.spec.(*config.HandlerLspResourceSpecLsp)
			if !ok {
				return nil, fmt.Errorf("lsp resource %v is not of type %T", name, (*config.HandlerLspResourceSpecLsp)(nil))
			}

			lspName, token, err := s.plane.Scope().LspAcquire(name, resource.scope)
			if err != nil {
				return nil, err
			}
			defer s.plane.Scope().LspRelease(lspName, token)

			if s.fileUri == nil || s.fileLine == nil || s.fileChar == nil {
				return nil, fmt.Errorf("file uri/line/char is not set")
			}

			params := &protocol.CompletionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: *s.fileUri,
					},
					Position: protocol.Position{
						Line:      *s.fileLine,
						Character: *s.fileChar,
					},
				},
			}

			lspList, err := s.plane.Lsp().TextDocumentCompletion(lspName, params)
			if err != nil {
				return nil, err
			}

			lspItems := make([]protocol.CompletionItem, len(lspList.Items))
			for i, s := range lspList.Items {
				lspItems[i] = protocol.CompletionItem{
					Label:         strings.TrimSpace(s.Label),
					Detail:        fmt.Sprintf("%v(%v)", info.MiddlewareServerName, lspName),
					Documentation: &protocol.Or_CompletionItem_documentation{Value: fmt.Sprintf("%v", lspName)},
				}
			}

			assignValue = lspItems

		case *config.StepJsonDumps:
			var input string
			if definition.Input == nil {
				return nil, fmt.Errorf("no input")
			}
			input = *definition.Input

			value, err := s.renderValue(input, memory)
			if err != nil {
				return nil, err
			}

			outputBytes, err := json.Marshal(value)
			if err != nil {
				return nil, err
			}

			assignValue = string(outputBytes)

		case *config.StepJsonParse:
			var input string
			if definition.Input == nil {
				return nil, fmt.Errorf("no input")
			}
			input = *definition.Input

			value, err := s.renderString(input, memory)
			if err != nil {
				return nil, err
			}

			if value != nil {
				var output any
				if err := json.Unmarshal([]byte(*value), &output); err != nil {
					return nil, err
				}
				assignValue = output
			} else {
				assignValue = nil
			}

		case *config.StepJsonParseCompletion:
			var input string
			if definition.Input == nil {
				return nil, fmt.Errorf("no input")
			}
			input = *definition.Input

			value, err := s.renderString(input, memory)
			if err != nil {
				return nil, err
			}

			if value != nil {
				var output []protocol.CompletionItem
				if err := json.Unmarshal([]byte(*value), &output); err != nil {
					return nil, err
				}
				assignValue = output
			} else {
				assignValue = nil
			}

		case *config.StepModelSimple:
		case *config.StepModelTemplate:
			// TODO

		default:
			return nil, fmt.Errorf("unexpected config.StepDefinition: %#v", definition)
		}

		// simple sequential steps
		if current.Step.Assign != nil {
			if err := s.set(memory, assignValue, *current.Step.Assign); err != nil {
				return nil, err
			}
		}

		current = current.Then
	}
}
