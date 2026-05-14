package config

import (
	"mals/internal/util"
	"mals/pkg/config"
	"mals/pkg/core"
)

type Step map[string]any

type StepModelMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (o *Step) Wire(c *config.Step) error {
	return nil
}

func unwireModelParameters(raw map[string]any) core.ModelParameters {
	var p core.ModelParameters
	if schemaRaw, ok := raw["schema"].(string); ok {
		schema := core.ModelSchema(schemaRaw)
		switch schema {
		case core.ModelSchemaJsonCompletionItems:
			p.Schema = core.ModelSchemaJsonCompletionItems
		default:
			p.Schema = ""
		}
	}
	if v, ok := raw["temperature"].(float64); ok {
		p.Temperature = &v
	}
	if v, ok := raw["max_tokens"].(int); ok {
		p.MaxTokens = util.Ptr(int64(v))
	}
	return p
}

func unwireStepLspCompletion(vm map[string]any) *config.StepLspCompletion {
	d := config.StepLspCompletion{}
	if resource, ok := vm["resource"].(string); ok {
		d.Resource = &resource
	}
	return &d
}

func unwireStepJsonDumps(vm map[string]any) *config.StepJsonDumps {
	d := config.StepJsonDumps{}
	if input, ok := vm["input"].(string); ok {
		d.Input = &input
	}
	return &d
}

func unwireStepJsonParse(vm map[string]any) *config.StepJsonParse {
	d := config.StepJsonParse{}
	if input, ok := vm["input"].(string); ok {
		d.Input = &input
	}
	return &d
}

func unwireStepJsonParseCompletion(vm map[string]any) *config.StepJsonParseCompletion {
	d := config.StepJsonParseCompletion{}
	if input, ok := vm["input"].(string); ok {
		d.Input = &input
	}
	return &d
}

func unwireStepModelRaw(vm map[string]any) *config.StepModelRaw {
	d := config.StepModelRaw{}
	if resource, ok := vm["resource"].(string); ok {
		d.Resource = &resource
	}
	if prompt, ok := vm["prompt"].(string); ok {
		d.Prompt = &prompt
	}
	if raw, ok := vm["parameters"].(map[string]interface{}); ok {
		d.Parameters = unwireModelParameters(raw)
	}
	return &d
}

func unwireStepModel(vm map[string]any) *config.StepModel {
	d := config.StepModel{}
	if resource, ok := vm["resource"].(string); ok {
		d.Resource = &resource
	}
	if prompt, ok := vm["prompt"].(string); ok {
		d.Prompt = &prompt
	}
	if messagesRaw, ok := vm["messages"].([]any); ok {
		messages := make([]*config.StepModelMessage, len(messagesRaw))
		for i, item := range messagesRaw {
			msgMap, ok := item.(map[string]any)
			if !ok {
				continue
			}
			messages[i] = &config.StepModelMessage{}
			role, _ := msgMap["role"].(string)
			switch core.ModelRole(role) {
			case core.ModelRoleSystem:
				messages[i].Role = core.ModelRoleSystem
			case core.ModelRoleUser:
				messages[i].Role = core.ModelRoleUser
			case core.ModelRoleAssistant:
				messages[i].Role = core.ModelRoleAssistant
			}
			if content, ok := msgMap["content"].(string); ok {
				messages[i].Content = &content
			}
		}
		d.Messages = messages
	}
	if raw, ok := vm["parameters"].(map[string]interface{}); ok {
		d.Parameters = unwireModelParameters(raw)
	}
	return &d
}

func unwireStepReturn(v any) *config.StepReturn {
	d := config.StepReturn{}
	if input, ok := v.(string); ok {
		d.Input = &input
	}
	return &d
}

func unwireStepIf(vm map[string]any) *config.StepIf {
	d := config.StepIf{}
	if condition, ok := vm["condition"].(string); ok {
		d.Condition = &condition
	}
	if thenRaw, ok := vm["then"].([]any); ok {
		thenSteps := make([]*config.Step, len(thenRaw))
		for i, item := range thenRaw {
			stepMap, ok := item.(map[string]any)
			if !ok {
				continue
			}
			step := Step(stepMap)
			if unwired, err := step.Unwire(); err == nil {
				thenSteps[i] = unwired
			}
		}
		d.Then = thenSteps
	}
	if elseRaw, ok := vm["else"].([]any); ok {
		elseSteps := make([]*config.Step, len(elseRaw))
		for i, item := range elseRaw {
			stepMap, ok := item.(map[string]any)
			if !ok {
				continue
			}
			step := Step(stepMap)
			if unwired, err := step.Unwire(); err == nil {
				elseSteps[i] = unwired
			}
		}
		d.Else = elseSteps
	}
	return &d
}

func unwireStepFor(vm map[string]any) *config.StepFor {
	d := config.StepFor{}
	if condition, ok := vm["condition"].(string); ok {
		d.Condition = &condition
	}
	if max, ok := vm["max"].(int); ok {
		d.Max = &max
	}
	if doRaw, ok := vm["do"].([]any); ok {
		doSteps := make([]*config.Step, len(doRaw))
		for i, item := range doRaw {
			stepMap, ok := item.(map[string]any)
			if !ok {
				continue
			}
			step := Step(stepMap)
			if unwired, err := step.Unwire(); err == nil {
				doSteps[i] = unwired
			}
		}
		d.Do = doSteps
	}
	return &d
}

func (o *Step) Unwire() (*config.Step, error) {
	c := config.Step{}

	for k, v := range *o {
		switch k {
		case "name":
			if vs, ok := v.(string); ok {
				c.Name = &vs
			}
		case "assign":
			if vs, ok := v.(string); ok {
				c.Assign = &vs
			}
		default:
			vm, _ := v.(map[string]any)

			switch k {
			case (&config.StepLspCompletion{}).StepDefinitionKind():
				if vm != nil {
					c.Definition = unwireStepLspCompletion(vm)
				}
			case (&config.StepJsonDumps{}).StepDefinitionKind():
				if vm != nil {
					c.Definition = unwireStepJsonDumps(vm)
				}
			case (&config.StepJsonParse{}).StepDefinitionKind():
				if vm != nil {
					c.Definition = unwireStepJsonParse(vm)
				}
			case (&config.StepJsonParseCompletion{}).StepDefinitionKind():
				if vm != nil {
					c.Definition = unwireStepJsonParseCompletion(vm)
				}
			case (&config.StepModelRaw{}).StepDefinitionKind():
				if vm != nil {
					c.Definition = unwireStepModelRaw(vm)
				}
			case (&config.StepModel{}).StepDefinitionKind():
				if vm != nil {
					c.Definition = unwireStepModel(vm)
				}
			case (&config.StepReturn{}).StepDefinitionKind():
				c.Definition = unwireStepReturn(v)
			case (&config.StepIf{}).StepDefinitionKind():
				if vm != nil {
					c.Definition = unwireStepIf(vm)
				}
			case (&config.StepFor{}).StepDefinitionKind():
				if vm != nil {
					c.Definition = unwireStepFor(vm)
				}
			}
		}
	}

	return &c, nil
}
