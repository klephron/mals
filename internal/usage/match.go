package usage

import (
	"mals/pkg/config"
	"regexp"
)

func matchRegex(value string, filter *string) bool {
	if filter == nil {
		return true
	}
	matched, _ := regexp.MatchString(*filter, value)
	return matched
}

func matchRegexAny(values []string, filter *string) bool {
	if filter == nil || values == nil {
		return true
	}
	for _, value := range values {
		if matchRegex(value, filter) {
			return true
		}
	}
	return false
}

func ConditionMatch(condition *config.Condition, filter ConditionFilter) bool {
	if condition == nil {
		return true
	}

	return matchRegexAny(condition.Filetypes, filter.Filetype) &&
		matchRegexAny(condition.Paths, filter.Path)
}

func ConditionMatchAny(conditions []*config.Condition, filter ConditionFilter) bool {
	if conditions == nil {
		return true
	}
	for _, condition := range conditions {
		if ConditionMatch(condition, filter) {
			return true
		}
	}
	return false
}

func EventMatch(event config.Event, filter EventFilter) bool {
	return matchRegex(string(event), (*string)(filter.Event))
}

func EventMatchAny(events []config.Event, filter EventFilter) bool {
	if events == nil {
		return true
	}
	for _, event := range events {
		if EventMatch(event, filter) {
			return true
		}
	}
	return false
}
