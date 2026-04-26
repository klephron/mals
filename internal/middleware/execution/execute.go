package execution

import (
	"encoding/json"
	"fmt"
	"mals/internal/lsp/protocol"
	"mals/internal/model"
	"mals/internal/util"
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

		next, value, err := s.executeStep(current, memory)
		if err != nil {
			return nil, err
		}
		if next == nil {
			return value, nil
		}
		if current.Step.Assign != nil {
			if err := s.set(memory, value, *current.Step.Assign); err != nil {
				return nil, err
			}
		}
		current = next
	}
}

func (s *ExecutionEnvironment) executeStep(current *executionNode, memory map[string]any) (*executionNode, any, error) {
	step := current.Step

	switch definition := step.Definition.(type) {
	case nil:
		return nil, nil, fmt.Errorf("step definition is nil")

	case *config.StepIf:
		return s.executeStepIf(current, definition, memory)

	case *config.StepFor:
		return s.executeStepFor(current, definition, memory)

	case *config.StepReturn:
		value, err := s.executeStepReturn(definition, memory)
		return nil, value, err

	case *config.StepLspCompletion:
		value, err := s.executeStepLspCompletion(definition)
		return current.Then, value, err

	case *config.StepJsonDumps:
		value, err := s.executeStepJsonDumps(definition, memory)
		return current.Then, value, err

	case *config.StepJsonParse:
		value, err := s.executeStepJsonParse(definition, memory)
		return current.Then, value, err

	case *config.StepJsonParseCompletion:
		value, err := s.executeStepJsonParseCompletion(definition, memory)
		return current.Then, value, err

	case *config.StepModelRaw:
		value, err := s.executeStepModelRaw(definition)
		return current.Then, value, err

	case *config.StepModel:
		value, err := s.executeStepModel(definition, memory)
		return current.Then, value, err

	default:
		return nil, nil, fmt.Errorf("unexpected config.StepDefinition: %#v", definition)
	}
}

func (s *ExecutionEnvironment) executeStepIf(current *executionNode, def *config.StepIf, memory map[string]any) (*executionNode, any, error) {
	if def.Condition == nil {
		return nil, nil, fmt.Errorf("if condition is nil")
	}
	condition, err := s.renderBool(*def.Condition, memory)
	if err != nil {
		return nil, nil, err
	}
	if condition == nil {
		return nil, nil, fmt.Errorf("condition is nil")
	}
	if *condition {
		return current.Then, nil, nil
	}
	return current.Else, nil, nil
}

func (s *ExecutionEnvironment) executeStepFor(current *executionNode, def *config.StepFor, memory map[string]any) (*executionNode, any, error) {
	if def.Condition == nil {
		return current.Then, nil, nil
	}
	condition, err := s.renderBool(*def.Condition, memory)
	if err != nil {
		return nil, nil, err
	}
	if condition == nil {
		return nil, nil, fmt.Errorf("condition is nil")
	}
	if *condition {
		return current.Then, nil, nil
	}
	return current.Else, nil, nil
}

func (s *ExecutionEnvironment) executeStepReturn(def *config.StepReturn, memory map[string]any) (any, error) {
	if def.Input == nil {
		return nil, nil
	}
	return s.renderValue(*def.Input, memory)
}

func (s *ExecutionEnvironment) executeStepLspCompletion(def *config.StepLspCompletion) (any, error) {
	if def.Resource == nil {
		return nil, fmt.Errorf("lsp resource not defined")
	}
	resource, ok := s.resources.Load(*def.Resource)
	if !ok {
		return nil, fmt.Errorf("lsp resource %v not found", *def.Resource)
	}

	lspSpec, ok := resource.spec.(*config.HandlerLspResourceSpecLsp)
	if !ok {
		return nil, fmt.Errorf("lsp resource %v is not of type %T", *def.Resource, (*config.HandlerLspResourceSpecLsp)(nil))
	}

	lspName, token, err := s.plane.Scope().LspAcquire(lspSpec.Name, resource.scope)
	if err != nil {
		return nil, err
	}
	defer s.plane.Scope().LspRelease(lspName, token)

	if s.fileUri == nil || s.fileLine == nil || s.fileChar == nil {
		return nil, fmt.Errorf("file uri/line/char is not set")
	}

	params := &protocol.CompletionParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{URI: *s.fileUri},
			Position:     protocol.Position{Line: *s.fileLine, Character: *s.fileChar},
		},
	}
	lspList, err := s.plane.Lsp().TextDocumentCompletion(lspName, params)
	if err != nil {
		return nil, err
	}

	items := make([]protocol.CompletionItem, len(lspList.Items))
	for i, item := range lspList.Items {
		items[i] = protocol.CompletionItem{
			Label:         strings.TrimSpace(item.Label),
			Detail:        fmt.Sprintf("%v(%v)", info.MiddlewareServerName, lspName),
			Documentation: &protocol.Or_CompletionItem_documentation{Value: fmt.Sprintf("%v", lspName)},
		}
	}

	return items, nil
}

func (s *ExecutionEnvironment) executeStepJsonDumps(def *config.StepJsonDumps, memory map[string]any) (any, error) {
	if def.Input == nil {
		return nil, fmt.Errorf("no input")
	}
	value, err := s.renderValue(*def.Input, memory)
	if err != nil {
		return nil, err
	}

	out, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	return string(out), nil
}

func (s *ExecutionEnvironment) executeStepJsonParse(def *config.StepJsonParse, memory map[string]any) (any, error) {
	if def.Input == nil {
		return nil, fmt.Errorf("no input")
	}
	value, err := s.renderString(*def.Input, memory)
	if err != nil {
		return nil, err
	}

	if value == nil {
		return nil, nil
	}

	var output any
	if err := json.Unmarshal([]byte(*value), &output); err != nil {
		return nil, err
	}

	return output, nil
}

func (s *ExecutionEnvironment) executeStepJsonParseCompletion(def *config.StepJsonParseCompletion, memory map[string]any) (any, error) {
	if def.Input == nil {
		return nil, fmt.Errorf("no input")
	}
	value, err := s.renderString(*def.Input, memory)
	if err != nil {
		return nil, err
	}

	if value == nil {
		return nil, nil
	}

	var output []protocol.CompletionItem
	if err := json.Unmarshal([]byte(*value), &output); err != nil {
		return nil, err
	}

	return output, nil
}

func (s *ExecutionEnvironment) executeStepModelRaw(def *config.StepModelRaw) (any, error) {
	if def.Resource == nil {
		return nil, fmt.Errorf("model resource not defined")
	}
	if def.Prompt == nil {
		return nil, fmt.Errorf("model prompt not defined")
	}

	resource, ok := s.resources.Load(*def.Resource)
	if !ok {
		return nil, fmt.Errorf("model resource %v not found", *def.Resource)
	}

	modelSpec, ok := resource.spec.(*config.HandlerLspResourceSpecModel)
	if !ok {
		return nil, fmt.Errorf("model resource %v is not of type %T", *def.Resource, (*config.HandlerLspResourceSpecModel)(nil))
	}

	modelName, token, err := s.plane.Scope().ModelAcquire(modelSpec.Name, resource.scope)
	if err != nil {
		return nil, err
	}
	defer s.plane.Scope().ModelRelease(modelName, token)

	task := model.NewTask(*def.Prompt, nil, nil, nil)

	return s.plane.Model().TaskExecClient(modelName, task, s.clientName)
}

func (s *ExecutionEnvironment) executeStepModel(def *config.StepModel, memory map[string]any) (any, error) {
	if def.Resource == nil {
		return nil, fmt.Errorf("model resource not defined")
	}
	if def.Prompt == nil {
		return nil, fmt.Errorf("model prompt not defined")
	}

	prompt, err := s.renderString(*def.Prompt, memory)
	if err != nil {
		return nil, err
	}
	if prompt == nil {
		prompt = util.Ptr("")
	}

	resource, ok := s.resources.Load(*def.Resource)
	if !ok {
		return nil, fmt.Errorf("model resource %v not found", *def.Resource)
	}

	modelSpec, ok := resource.spec.(*config.HandlerLspResourceSpecModel)
	if !ok {
		return nil, fmt.Errorf("model resource %v is not of type %T", *def.Resource, (*config.HandlerLspResourceSpecModel)(nil))
	}

	modelName, token, err := s.plane.Scope().ModelAcquire(modelSpec.Name, resource.scope)
	if err != nil {
		return nil, err
	}
	defer s.plane.Scope().ModelRelease(modelName, token)

	task := model.NewTask(*prompt, nil, nil, nil)

	return s.plane.Model().TaskExecClient(modelName, task, s.clientName)
}
