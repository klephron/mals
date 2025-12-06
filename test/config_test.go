package test

import (
	"bytes"
	"encoding/json"
	. "mals/pkg/config"
	"testing"
)

func Test_Unmarshall(t *testing.T) {
	tests := []struct {
		name          string
		input         []byte
		expected      *Config
		expectedError bool
	}{
		{
			name:  "empty",
			input: []byte("{}"),
			expected: &Config{
				Loggers:   []Log{},
				Listeners: []*Listener{},
				Models:    []*Model{},
				Lsps:      []*Lsp{},
				Usages:    []*Usage{},
			},
			expectedError: false,
		},
		{
			name:          "empty error",
			input:         []byte("{"),
			expected:      nil,
			expectedError: true,
		},
		{
			name:  "empty explicit",
			input: []byte(`{"models":[], "lsps":[], "usages":[]}`),
			expected: &Config{
				Loggers:   []Log{},
				Listeners: []*Listener{},
				Models:    []*Model{},
				Lsps:      []*Lsp{},
				Usages:    []*Usage{},
			},
			expectedError: false,
		},
		{
			name:  "log",
			input: []byte(`{"loggers": [{"name": "l1", "type": "file", "file": "logFile", "level": "debug"}], "models":[], "lsps":[], "usages":[]}`),
			expected: &Config{
				Loggers: []Log{
					&LogFile{
						LogGeneric: NewLogGeneric("l1"),
						File:       "logFile",
						Level:      "debug",
					},
				},
				Listeners: []*Listener{},
				Models:    []*Model{},
				Lsps:      []*Lsp{},
				Usages:    []*Usage{},
			},
			expectedError: false,
		},
		{
			name:          "log error no type",
			input:         []byte(`{"loggers": [{"file": "logFile", "level": "debug"}], "models":[], "lsps":[], "usages":[]}`),
			expected:      nil,
			expectedError: true,
		},
		{
			name:  "listeners empty",
			input: []byte(`{"listeners": [], "models":[], "lsps":[], "usages":[]}`),
			expected: &Config{
				Loggers:   []Log{},
				Listeners: []*Listener{},
				Models:    []*Model{},
				Lsps:      []*Lsp{},
				Usages:    []*Usage{},
			},
			expectedError: false,
		},
		{
			name:  "listeners",
			input: []byte(`{"listeners": [{"name": "l1", "type": "rest", "port": 12},{"type": "lsp", "port": 999}], "models":[], "lsps":[], "usages":[]}`),
			expected: &Config{
				Loggers: []Log{},
				Listeners: []*Listener{
					&Listener{
						Name: "l1",
						Type: "rest",
						Port: 12,
					},
					&Listener{
						Type: "lsp",
						Port: 999,
					},
				},
				Models: []*Model{},
				Lsps:   []*Lsp{},
				Usages: []*Usage{},
			},
			expectedError: false,
		},
		{
			name:  "model spec OpenAI",
			input: []byte(`{"models":[{"name": "1", "spec": "openai", "settings": {"url": "something", "max_tokens": 12, "temperature": 1}}], "lsps":[], "usages":[]}`),
			expected: &Config{
				Loggers:   []Log{},
				Listeners: []*Listener{},
				Models: []*Model{
					&Model{
						Name: "1",
						Settings: &ModelSpecOpenAI{
							Url:         "something",
							MaxTokens:   12,
							Temperature: 1,
						},
					},
				},
				Lsps:   []*Lsp{},
				Usages: []*Usage{},
			},
			expectedError: false,
		},
		{
			name:  "model spec OpenAI default",
			input: []byte(`{"models":[{"name": "1", "spec": "openai", "settings": {}}]}`),
			expected: &Config{
				Loggers:   []Log{},
				Listeners: []*Listener{},
				Models: []*Model{
					&Model{
						Name:     "1",
						Settings: &ModelSpecOpenAI{},
					},
				},
				Lsps:   []*Lsp{},
				Usages: []*Usage{},
			},
			expectedError: false,
		},
		{
			name:  "lsp spec stdio",
			input: []byte(`{"models":[], "lsps":[{"name": "1", "spec": "stdio", "settings": {"cmd": ["something", "arg1", "arg2"]}}], "usages":[]}`),
			expected: &Config{
				Loggers:   []Log{},
				Listeners: []*Listener{},
				Models:    []*Model{},
				Lsps: []*Lsp{
					&Lsp{
						Name: "1",
						Settings: &LspSpecStdio{
							Cmd: []string{"something", "arg1", "arg2"},
						},
					},
				},
				Usages: []*Usage{},
			},
			expectedError: false,
		},
		{
			name:  "lsp spec stdio empty cmd",
			input: []byte(`{"models":[], "lsps":[{"name": "1", "spec": "stdio", "settings": {"cmd": []}}], "usages":[]}`),
			expected: &Config{
				Loggers:   []Log{},
				Listeners: []*Listener{},
				Models:    []*Model{},
				Lsps: []*Lsp{
					&Lsp{
						Name: "1",
						Settings: &LspSpecStdio{
							Cmd: []string{},
						},
					},
				},
				Usages: []*Usage{},
			},
			expectedError: false,
		},
		{
			name:  "lsp spec stdio nil cmd",
			input: []byte(`{"lsps":[{"name": "1", "spec": "stdio", "settings": {}}]}`),
			expected: &Config{
				Loggers:   []Log{},
				Listeners: []*Listener{},
				Models:    []*Model{},
				Lsps: []*Lsp{
					&Lsp{
						Name: "1",
						Settings: &LspSpecStdio{
							Cmd: []string{},
						},
					},
				},
				Usages: []*Usage{},
			},
			expectedError: false,
		},
		{
			name:  "usage condition nil workflow nil",
			input: []byte(`{"usages":[{"name": "1"}]}`),
			expected: &Config{
				Loggers:   []Log{},
				Listeners: []*Listener{},
				Models:    []*Model{},
				Lsps:      []*Lsp{},
				Usages: []*Usage{
					&Usage{
						Name:       "1",
						Conditions: []*Condition{},
						Workflow:   nil,
					},
				},
			},
			expectedError: false,
		},
		{
			name:  "usage condition empty workflow nil",
			input: []byte(`{"usages":[{"name": "1", "conditions": []}]}`),
			expected: &Config{
				Loggers:   []Log{},
				Listeners: []*Listener{},
				Models:    []*Model{},
				Lsps:      []*Lsp{},
				Usages: []*Usage{
					&Usage{
						Name:       "1",
						Conditions: []*Condition{},
						Workflow:   nil,
					},
				},
			},
			expectedError: false,
		},
		{
			name:  "usage condition empty",
			input: []byte(`{"usages":[{"name": "1", "conditions": [{}]}]}`),
			expected: &Config{
				Loggers:   []Log{},
				Listeners: []*Listener{},
				Models:    []*Model{},
				Lsps:      []*Lsp{},
				Usages: []*Usage{
					&Usage{
						Name: "1",
						Conditions: []*Condition{
							&Condition{
								Filetypes: []string{},
								Paths:     []string{},
								Types:     []string{},
							},
						},
						Workflow: nil,
					},
				},
			},
			expectedError: false,
		},
		{
			name:  "usage condition",
			input: []byte(`{"usages":[{"name": "1", "conditions": [{"filetypes": ["s1"], "paths": ["s2", "s3"], "types": ["s4", "s5", "s6"]}]}]}`),
			expected: &Config{
				Loggers:   []Log{},
				Listeners: []*Listener{},
				Models:    []*Model{},
				Lsps:      []*Lsp{},
				Usages: []*Usage{
					&Usage{
						Name: "1",
						Conditions: []*Condition{
							&Condition{
								Filetypes: []string{"s1"},
								Paths:     []string{"s2", "s3"},
								Types:     []string{"s4", "s5", "s6"},
							},
						},
						Workflow: nil,
					},
				},
			},
			expectedError: false,
		},
		{
			name:  "usage workflow empty",
			input: []byte(`{"usages":[{"name": "1", "workflow": {}}]}`),
			expected: &Config{
				Loggers:   []Log{},
				Listeners: []*Listener{},
				Models:    []*Model{},
				Lsps:      []*Lsp{},
				Usages: []*Usage{
					&Usage{
						Name:       "1",
						Conditions: []*Condition{},
						Workflow: &Workflow{
							Steps: []Step{},
						},
					},
				},
			},
			expectedError: false,
		},
		{
			name:  "usage workflow name",
			input: []byte(`{"usages":[{"workflow": {"name": "w1"}}]}`),
			expected: &Config{
				Loggers:   []Log{},
				Listeners: []*Listener{},
				Models:    []*Model{},
				Lsps:      []*Lsp{},
				Usages: []*Usage{
					&Usage{
						Conditions: []*Condition{},
						Workflow: &Workflow{
							Name:  "w1",
							Steps: []Step{},
						},
					},
				},
			},
			expectedError: false,
		},
		{
			name:  "usage workflow steps empty",
			input: []byte(`{"usages":[{"workflow": {"steps":[]}}]}`),
			expected: &Config{
				Loggers:   []Log{},
				Listeners: []*Listener{},
				Models:    []*Model{},
				Lsps:      []*Lsp{},
				Usages: []*Usage{
					&Usage{
						Conditions: []*Condition{},
						Workflow: &Workflow{
							Steps: []Step{},
						},
					},
				},
			},
			expectedError: false,
		},
		{
			name:  "usage workflow step lsp empty",
			input: []byte(`{"usages":[{"workflow": {"steps":[{"lsp": "1"}]}}]}`),
			expected: &Config{
				Loggers:   []Log{},
				Listeners: []*Listener{},
				Models:    []*Model{},
				Lsps:      []*Lsp{},
				Usages: []*Usage{
					&Usage{
						Conditions: []*Condition{},
						Workflow: &Workflow{
							Steps: []Step{
								StepLsp{
									StepGeneric: StepGeneric{
										Conditions: []*Condition{},
									},
									Lsp: "1",
								},
							},
						},
					},
				},
			},
			expectedError: false,
		},
		{
			name:  "usage workflow step lsp",
			input: []byte(`{"usages":[{"workflow": {"steps":[{"name": "1", "conditions": [], "lsp": "l1", "template": "t1"}]}}]}`),
			expected: &Config{
				Loggers:   []Log{},
				Listeners: []*Listener{},
				Models:    []*Model{},
				Lsps:      []*Lsp{},
				Usages: []*Usage{
					&Usage{
						Conditions: []*Condition{},
						Workflow: &Workflow{
							Steps: []Step{
								StepLsp{
									StepGeneric: StepGeneric{
										Name:       "1",
										Conditions: []*Condition{},
									},
									Lsp:      "l1",
									Template: "t1",
								},
							},
						},
					},
				},
			},
			expectedError: false,
		},
		{
			name:  "usage workflow step model",
			input: []byte(`{"usages":[{"workflow": {"steps":[{"name": "1", "conditions": [], "model": "m1", "template": "t1"}]}}]}`),
			expected: &Config{
				Loggers:   []Log{},
				Listeners: []*Listener{},
				Models:    []*Model{},
				Lsps:      []*Lsp{},
				Usages: []*Usage{
					&Usage{
						Conditions: []*Condition{},
						Workflow: &Workflow{
							Steps: []Step{
								StepModel{
									StepGeneric: StepGeneric{
										Name:       "1",
										Conditions: []*Condition{},
									},
									Model:    "m1",
									Template: "t1",
								},
							},
						},
					},
				},
			},
			expectedError: false,
		},
		{
			name:          "usage workflow step error none set",
			input:         []byte(`{"usages":[{"workflow": {"steps":[{"name": "1", "conditions": [], "template": "t1"}]}}]}`),
			expected:      nil,
			expectedError: true,
		},
		{
			name:          "usage workflow step error multiple set",
			input:         []byte(`{"usages":[{"workflow": {"steps":[{"name": "1", "conditions": [], "template": "t1", "model": "m1", "lsp": "l1"}]}}]}`),
			expected:      nil,
			expectedError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var config Config
			err := json.Unmarshal(tt.input, &config)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Unmarshal() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unmarshal() unexpected error: %v", err)
				return
			}

			configJson, err := json.Marshal(config)
			expectedJson, err := json.Marshal(tt.expected)

			if bytes.Compare(configJson, expectedJson) != 0 {
				t.Errorf("expected %v, got %v", string(expectedJson), string(configJson))
				return
			}
		})
	}
}
