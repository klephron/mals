package wire

import (
	"fmt"
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
			if s, ok := v.(string); ok {
				c.Name = &s
			}
		case "assign":
			if s, ok := v.(string); ok {
				c.Assign = &s
			}
		default:
			m, ok := v.(map[string]any)
			if !ok {
				continue
			}

			switch k {
			case (&config.StepLspCompetion{}).StepDefinitionKind():
				d := config.StepLspCompetion{}
				if resource, ok := m["resource"].(string); ok {
					d.Resource = &resource
				}
				c.Definition = &d

			case (&config.StepJsonDumps{}).StepDefinitionKind():
				d := config.StepJsonDumps{}
				if input, ok := m["input"].(string); ok {
					d.Input = &input
				}
				c.Definition = &d

			case (&config.StepJsonParse{}).StepDefinitionKind():
				d := config.StepJsonParse{}
				if input, ok := m["input"].(string); ok {
					d.Input = &input
				}
				c.Definition = &d

			case (&config.StepJsonParseCompletion{}).StepDefinitionKind():
				d := config.StepJsonParseCompletion{}
				if input, ok := m["input"].(string); ok {
					d.Input = &input
				}
				c.Definition = &d

			case (&config.StepModelSimple{}).StepDefinitionKind():
				d := config.StepModelSimple{}
				if resource, ok := m["resource"].(string); ok {
					d.Resource = &resource
				}
				if prompt, ok := m["prompt"].(string); ok {
					d.Prompt = &prompt
				}
				c.Definition = &d

			case (&config.StepModelTemplate{}).StepDefinitionKind():
				d := config.StepModelTemplate{}
				if resource, ok := m["resource"].(string); ok {
					d.Resource = &resource
				}
				if prompt, ok := m["prompt"].(string); ok {
					d.Prompt = &prompt
				}
				c.Definition = &d

			case (&config.StepReturn{}).StepDefinitionKind():
				d := config.StepReturn{}
				if input, ok := m["input"].(string); ok {
					d.Input = &input
				}
				c.Definition = &d

			case (&config.StepIf{}).StepDefinitionKind():
				d := config.StepIf{}
				if condition, ok := m["condition"].(string); ok {
					d.Condition = &condition
				}
				if thenRaw, ok := m["then"].([]map[string]any); ok {
					thenSteps := make([]*config.Step, len(thenRaw))
					for i, stepMap := range thenRaw {
						step := Step(stepMap)
						if unwired, err := step.Unwire(); err == nil {
							thenSteps[i] = unwired
						}
					}
					d.Then = thenSteps
				} else {
					fmt.Println("something is wrong")
				}
				if elseRaw, ok := m["else"].([]map[string]any); ok {
					elseSteps := make([]*config.Step, len(elseRaw))
					for i, stepMap := range elseRaw {
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
				if condition, ok := m["condition"].(string); ok {
					d.Condition = &condition
				}
				if max, ok := m["max"].(int); ok {
					d.Max = &max
				}
				if doRaw, ok := m["do"].([]map[string]any); ok {
					doSteps := make([]*config.Step, len(doRaw))
					for i, stepMap := range doRaw {
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
