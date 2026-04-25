package wire

import (
	"mals/pkg/config"
)

func (o *Step) Wire(c *config.Step) error {
	return nil
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
			switch k {
			case (&config.StepLspCompletion{}).StepDefinitionKind():
				d := config.StepLspCompletion{}
				vm, ok := v.(map[string]any)
				if !ok {
					continue
				}
				if resource, ok := vm["resource"].(string); ok {
					d.Resource = &resource
				}
				c.Definition = &d

			case (&config.StepJsonDumps{}).StepDefinitionKind():
				d := config.StepJsonDumps{}
				vm, ok := v.(map[string]any)
				if !ok {
					continue
				}
				if input, ok := vm["input"].(string); ok {
					d.Input = &input
				}
				c.Definition = &d

			case (&config.StepJsonParse{}).StepDefinitionKind():
				d := config.StepJsonParse{}
				vm, ok := v.(map[string]any)
				if !ok {
					continue
				}
				if input, ok := vm["input"].(string); ok {
					d.Input = &input
				}
				c.Definition = &d

			case (&config.StepJsonParseCompletion{}).StepDefinitionKind():
				d := config.StepJsonParseCompletion{}
				vm, ok := v.(map[string]any)
				if !ok {
					continue
				}
				if input, ok := vm["input"].(string); ok {
					d.Input = &input
				}
				c.Definition = &d

			case (&config.StepModelSimple{}).StepDefinitionKind():
				d := config.StepModelSimple{}
				vm, ok := v.(map[string]any)
				if !ok {
					continue
				}
				if resource, ok := vm["resource"].(string); ok {
					d.Resource = &resource
				}
				if prompt, ok := vm["prompt"].(string); ok {
					d.Prompt = &prompt
				}
				c.Definition = &d

			case (&config.StepModelTemplate{}).StepDefinitionKind():
				d := config.StepModelTemplate{}
				vm, ok := v.(map[string]any)
				if !ok {
					continue
				}
				if resource, ok := vm["resource"].(string); ok {
					d.Resource = &resource
				}
				if prompt, ok := vm["prompt"].(string); ok {
					d.Prompt = &prompt
				}
				c.Definition = &d

			case (&config.StepReturn{}).StepDefinitionKind():
				d := config.StepReturn{}
				if input, ok := v.(string); ok {
					d.Input = &input
				}
				c.Definition = &d

			case (&config.StepIf{}).StepDefinitionKind():
				d := config.StepIf{}
				vm, ok := v.(map[string]any)
				if !ok {
					continue
				}
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
				c.Definition = &d

			case (&config.StepWhile{}).StepDefinitionKind():
				d := config.StepWhile{}
				vm, ok := v.(map[string]any)
				if !ok {
					continue
				}
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
				c.Definition = &d
			}
		}
	}

	return &c, nil
}
