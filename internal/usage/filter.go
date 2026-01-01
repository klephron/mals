package usage

import (
	"mals/pkg/config"
)

func StepsFilter(steps []*config.Step, filter ConditionFilter) []*config.Step {
	filtered := make([]*config.Step, 0)

	for _, step := range steps {
		conditions := step.Conditions

		if ConditionMatchAny(conditions, filter) {
			filtered = append(filtered, step)
		}
	}

	return filtered
}

func UsagesFilter(usages []*config.Usage, condition ConditionFilter, event EventFilter) []*config.Usage {
	filtered := make([]*config.Usage, 0)

	for _, usage := range usages {
		if ConditionMatchAny(usage.Conditions, condition) && EventMatchAny(usage.Events, event) {
			newUsage := &config.Usage{
				Name:       usage.Name,
				Conditions: usage.Conditions,
				Workflow: &config.Workflow{
					Steps: StepsFilter(usage.Workflow.Steps, condition),
				},
			}

			filtered = append(filtered, newUsage)
		}
	}

	return filtered
}
