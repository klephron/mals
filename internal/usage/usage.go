package usage

import (
	"mals/pkg/config"
	"regexp"
)

func regexMatch(value *string, pattern string) bool {
	if value == nil {
		return true
	}
	matched, _ := regexp.MatchString(pattern, *value)
	return matched
}

func filetypeMatch(filetype *string, patterns []string) bool {
	if len(patterns) == 0 {
		return true
	}
	for _, pattern := range patterns {
		if regexMatch(filetype, pattern) {
			return true
		}
	}
	return false
}

func pathMatch(filetype *string, patterns []string) bool {
	if len(patterns) == 0 {
		return true
	}
	for _, pattern := range patterns {
		if regexMatch(filetype, pattern) {
			return true
		}
	}
	return false
}

func eventMatch(filetype *string, patterns []string) bool {
	if len(patterns) == 0 {
		return true
	}
	for _, pattern := range patterns {
		if regexMatch(filetype, pattern) {
			return true
		}
	}
	return false
}

func conditionMatch(filetype *string, path *string, event *string, condition *config.Condition) bool {
	if condition == nil {
		return true
	}
	return filetypeMatch(filetype, condition.Filetypes) &&
		pathMatch(path, condition.Paths) &&
		eventMatch(event, condition.Events)
}

func conditionAnyMatch(filetype *string, path *string, event *string, conditions []*config.Condition) bool {
	if len(conditions) == 0 {
		return true
	}
	for _, condition := range conditions {
		if conditionMatch(filetype, path, event, condition) {
			return true
		}
	}
	return false
}

func StepsFilter(filetype *string, path *string, event *string, steps []config.Step) []config.Step {
	filtered := make([]config.Step, 0)

	for _, step := range steps {
		var conditions []*config.Condition

		switch step := step.(type) {
		case *config.StepGeneric:
			conditions = step.Conditions
		case *config.StepModel:
			conditions = step.Conditions
		case *config.StepLsp:
			conditions = step.Conditions
		default:
			filtered = append(filtered, step)
			continue
		}

		if conditionAnyMatch(filetype, path, event, conditions) {
			filtered = append(filtered, step)
		}
	}

	return filtered
}

func UsagesFilter(filetype *string, path *string, event *string, usages []*config.Usage) []*config.Usage {
	filtered := make([]*config.Usage, 0)

	for _, usage := range usages {
		if !conditionAnyMatch(filetype, path, event, usage.Conditions) {
			continue
		}

		newUsage := &config.Usage{
			Name:       usage.Name,
			Conditions: usage.Conditions,
			Workflow: &config.Workflow{
				Name:  usage.Workflow.Name,
				Steps: StepsFilter(filetype, path, event, usage.Workflow.Steps),
			},
		}

		filtered = append(filtered, newUsage)
	}

	return filtered
}
