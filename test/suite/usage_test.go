package test

import (
	"mals/internal/usage"
	"mals/internal/util"
	"mals/pkg/config"
	"testing"
)

func Test_StepsFilter(t *testing.T) {
	tests := []struct {
		name          string
		steps         []*config.Step
		filter        usage.ConditionFilter
		expectedCount int
	}{
		{
			name:          "empty steps returns empty",
			steps:         []*config.Step{},
			filter:        usage.ConditionFilter{},
			expectedCount: 0,
		},
		{
			name: "all steps match",
			steps: []*config.Step{
				{Conditions: []*config.Condition{{Filetypes: []string{"go"}, Paths: []string{"/src"}}}},
				{Conditions: []*config.Condition{{Filetypes: []string{"go"}, Paths: []string{"/test"}}}},
			},
			filter: usage.ConditionFilter{
				Filetype: util.Ptr("go"),
				Path:     nil,
			},
			expectedCount: 2,
		},
		{
			name: "some steps match",
			steps: []*config.Step{
				{Conditions: []*config.Condition{{Filetypes: []string{"go"}, Paths: []string{"/src"}}}},
				{Conditions: []*config.Condition{{Filetypes: []string{"py"}, Paths: []string{"/src"}}}},
			},
			filter: usage.ConditionFilter{
				Filetype: util.Ptr("go"),
				Path:     util.Ptr("/src"),
			},
			expectedCount: 1,
		},
		{
			name: "no steps match",
			steps: []*config.Step{
				{Conditions: []*config.Condition{{Filetypes: []string{"py"}, Paths: []string{"/src"}}}},
			},
			filter: usage.ConditionFilter{
				Filetype: util.Ptr("go"),
				Path:     util.Ptr("/src"),
			},
			expectedCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := usage.StepsFilter(tt.steps, tt.filter)
			if len(result) != tt.expectedCount {
				t.Errorf("StepsFilter() returned %d steps, want %d", len(result), tt.expectedCount)
			}
		})
	}
}

func Test_UsagesFilter(t *testing.T) {
	tests := []struct {
		name          string
		usages        []*config.Usage
		condition     usage.ConditionFilter
		event         usage.EventFilter
		expectedCount int
		expectedSteps []int
	}{
		{
			name:          "empty usages returns empty",
			usages:        []*config.Usage{},
			condition:     usage.ConditionFilter{},
			event:         usage.EventFilter{},
			expectedCount: 0,
			expectedSteps: []int{},
		},
		{
			name: "usage matches condition and event",
			usages: []*config.Usage{
				{
					Name:       "test-usage",
					Events:     []config.Event{"save", "open"},
					Conditions: []*config.Condition{{Filetypes: []string{"go"}, Paths: []string{"/src"}}},
					Workflow: &config.Workflow{
						Steps: []*config.Step{
							{Conditions: []*config.Condition{{Filetypes: []string{"go"}, Paths: []string{"/src"}}}},
						},
					},
				},
			},
			condition: usage.ConditionFilter{
				Filetype: util.Ptr("go"),
				Path:     util.Ptr("/src"),
			},
			event: usage.EventFilter{
				Event: util.Ptr(config.Event("save")),
			},
			expectedCount: 1,
			expectedSteps: []int{1},
		},
		{
			name: "usage matches condition but not event",
			usages: []*config.Usage{
				{
					Name:       "test-usage",
					Events:     []config.Event{"save"},
					Conditions: []*config.Condition{{Filetypes: []string{"go"}, Paths: []string{"/src"}}},
					Workflow:   &config.Workflow{Steps: []*config.Step{}},
				},
			},
			condition: usage.ConditionFilter{
				Filetype: util.Ptr("go"),
				Path:     util.Ptr("/src"),
			},
			event: usage.EventFilter{
				Event: util.Ptr(config.Event("open")),
			},
			expectedCount: 0,
			expectedSteps: []int{},
		},
		{
			name: "usage matches event but not condition",
			usages: []*config.Usage{
				{
					Name:       "test-usage",
					Events:     []config.Event{"save"},
					Conditions: []*config.Condition{{Filetypes: []string{"py"}, Paths: []string{"/src"}}},
					Workflow:   &config.Workflow{Steps: []*config.Step{}},
				},
			},
			condition: usage.ConditionFilter{
				Filetype: util.Ptr("go"),
				Path:     util.Ptr("/src"),
			},
			event: usage.EventFilter{
				Event: util.Ptr(config.Event("save")),
			},
			expectedCount: 0,
			expectedSteps: []int{},
		},
		{
			name: "steps filtered correctly within matching usage",
			usages: []*config.Usage{
				{
					Name:       "test-usage",
					Events:     []config.Event{"save"},
					Conditions: []*config.Condition{{Filetypes: []string{"go"}, Paths: []string{"/src"}}},
					Workflow: &config.Workflow{
						Steps: []*config.Step{
							{Conditions: []*config.Condition{{Filetypes: []string{"go"}, Paths: []string{"/src"}}}},
							{Conditions: []*config.Condition{{Filetypes: []string{"py"}, Paths: []string{"/src"}}}},
						},
					},
				},
			},
			condition: usage.ConditionFilter{
				Filetype: util.Ptr("go"),
				Path:     util.Ptr("/src"),
			},
			event: usage.EventFilter{
				Event: util.Ptr(config.Event("save")),
			},
			expectedCount: 1,
			expectedSteps: []int{1},
		},
		{
			name: "multiple usages some match",
			usages: []*config.Usage{
				{
					Name:       "usage1",
					Events:     []config.Event{"save"},
					Conditions: []*config.Condition{{Filetypes: []string{"go"}, Paths: []string{"/src"}}},
					Workflow:   &config.Workflow{Steps: []*config.Step{{Conditions: []*config.Condition{{Filetypes: []string{"go"}, Paths: []string{"/src"}}}}}},
				},
				{
					Name:       "usage2",
					Events:     []config.Event{"open"},
					Conditions: []*config.Condition{{Filetypes: []string{"go"}, Paths: []string{"/src"}}},
					Workflow:   &config.Workflow{Steps: []*config.Step{}},
				},
			},
			condition: usage.ConditionFilter{
				Filetype: util.Ptr("go"),
				Path:     util.Ptr("/src"),
			},
			event: usage.EventFilter{
				Event: util.Ptr(config.Event("save")),
			},
			expectedCount: 1,
			expectedSteps: []int{1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := usage.UsagesFilter(tt.usages, tt.condition, tt.event)
			if len(result) != tt.expectedCount {
				t.Errorf("UsagesFilter() returned %d usages, want %d", len(result), tt.expectedCount)
			}
			if len(tt.expectedSteps) > 0 {
				for i, expected := range tt.expectedSteps {
					if i < len(result) {
						if len(result[i].Workflow.Steps) != expected {
							t.Errorf("Usage %d has %d steps, want %d", i, len(result[i].Workflow.Steps), expected)
						}
					}
				}
			}
		})
	}
}
