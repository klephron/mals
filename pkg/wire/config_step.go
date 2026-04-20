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
			}
		}
	}

	return &c, nil
}
