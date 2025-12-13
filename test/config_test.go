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
				Listeners: []Listener{},
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
				Listeners: []Listener{},
				Models:    []*Model{},
				Lsps:      []*Lsp{},
				Usages:    []*Usage{},
			},
			expectedError: false,
		},
		{
			name:  "log",
			input: []byte(`{"loggers": [{"name": "l1", "kind": "file", "file": "logFile", "level": "debug"}], "models":[], "lsps":[], "usages":[]}`),
			expected: &Config{
				Loggers: []Log{
					&LogFile{
						LogGeneric: NewLogGeneric("l1"),
						File:       "logFile",
						Level:      "debug",
					},
				},
				Listeners: []Listener{},
				Models:    []*Model{},
				Lsps:      []*Lsp{},
				Usages:    []*Usage{},
			},
			expectedError: false,
		},
		{
			name:          "log error no kind",
			input:         []byte(`{"loggers": [{"file": "logFile", "level": "debug"}], "models":[], "lsps":[], "usages":[]}`),
			expected:      nil,
			expectedError: true,
		},
		{
			name:  "listeners empty",
			input: []byte(`{"listeners": [], "models":[], "lsps":[], "usages":[]}`),
			expected: &Config{
				Loggers:   []Log{},
				Listeners: []Listener{},
				Models:    []*Model{},
				Lsps:      []*Lsp{},
				Usages:    []*Usage{},
			},
			expectedError: false,
		},
		{
			name:  "listeners",
			input: []byte(`{"listeners": [{"name": "l1", "kind": "api", "ipc": "tcp", "port": 12},{"kind": "lsp", "ipc": "stdio"}], "models":[], "lsps":[], "usages":[]}`),
			expected: &Config{
				Loggers: []Log{},
				Listeners: []Listener{
					&ListenerTcp{
						ListenerGeneric: NewListenerGeneric("l1", "api"),
						Port:            12,
					},
					&ListenerStdio{
						ListenerGeneric: NewListenerGeneric("", "lsp"),
					},
				},
				Models: []*Model{},
				Lsps:   []*Lsp{},
				Usages: []*Usage{},
			},
			expectedError: false,
		},
		{
			name:          "listener error no ipc 1",
			input:         []byte(`{"listeners": [{"name": "l1", "kind": "api", "port": 12}], "models":[], "lsps":[], "usages":[]}`),
			expected:      nil,
			expectedError: true,
		},
		{
			name:          "listener error no ipc 2",
			input:         []byte(`{"listeners": [{"name": "l1", "kind": "lsp"}], "models":[], "lsps":[], "usages":[]}`),
			expected:      nil,
			expectedError: true,
		},
		{
			name:  "model kind OpenAI",
			input: []byte(`{"models":[{"name": "1", "kind": "openai", "settings": {"url": "something", "max_tokens": 12, "temperature": 1}}], "lsps":[], "usages":[]}`),
			expected: &Config{
				Loggers:   []Log{},
				Listeners: []Listener{},
				Models: []*Model{
					&Model{
						Name: "1",
						Settings: &ModelSettingsOpenAI{
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
			name:  "model kind OpenAI default",
			input: []byte(`{"models":[{"name": "1", "kind": "openai", "settings": {}}]}`),
			expected: &Config{
				Loggers:   []Log{},
				Listeners: []Listener{},
				Models: []*Model{
					&Model{
						Name:     "1",
						Settings: &ModelSettingsOpenAI{},
					},
				},
				Lsps:   []*Lsp{},
				Usages: []*Usage{},
			},
			expectedError: false,
		},
		{
			name:  "lsp kind stdio",
			input: []byte(`{"models":[], "lsps":[{"name": "1", "kind": "stdio", "settings": {"cmd": ["something", "arg1", "arg2"]}}], "usages":[]}`),
			expected: &Config{
				Loggers:   []Log{},
				Listeners: []Listener{},
				Models:    []*Model{},
				Lsps: []*Lsp{
					&Lsp{
						Name: "1",
						Settings: &LspSettingsStdio{
							Cmd: []string{"something", "arg1", "arg2"},
						},
					},
				},
				Usages: []*Usage{},
			},
			expectedError: false,
		},
		{
			name:  "lsp kind stdio empty cmd",
			input: []byte(`{"models":[], "lsps":[{"name": "1", "kind": "stdio", "settings": {"cmd": []}}], "usages":[]}`),
			expected: &Config{
				Loggers:   []Log{},
				Listeners: []Listener{},
				Models:    []*Model{},
				Lsps: []*Lsp{
					&Lsp{
						Name: "1",
						Settings: &LspSettingsStdio{
							Cmd: []string{},
						},
					},
				},
				Usages: []*Usage{},
			},
			expectedError: false,
		},
		{
			name:  "lsp kind stdio nil cmd",
			input: []byte(`{"lsps":[{"name": "1", "kind": "stdio", "settings": {}}]}`),
			expected: &Config{
				Loggers:   []Log{},
				Listeners: []Listener{},
				Models:    []*Model{},
				Lsps: []*Lsp{
					&Lsp{
						Name: "1",
						Settings: &LspSettingsStdio{
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
				Listeners: []Listener{},
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
				Listeners: []Listener{},
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
				Listeners: []Listener{},
				Models:    []*Model{},
				Lsps:      []*Lsp{},
				Usages: []*Usage{
					&Usage{
						Name: "1",
						Conditions: []*Condition{
							&Condition{
								Filetypes: []string{},
								Paths:     []string{},
								Events:    []string{},
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
			input: []byte(`{"usages":[{"name": "1", "conditions": [{"filetypes": ["s1"], "paths": ["s2", "s3"], "events": ["s4", "s5", "s6"]}]}]}`),
			expected: &Config{
				Loggers:   []Log{},
				Listeners: []Listener{},
				Models:    []*Model{},
				Lsps:      []*Lsp{},
				Usages: []*Usage{
					&Usage{
						Name: "1",
						Conditions: []*Condition{
							&Condition{
								Filetypes: []string{"s1"},
								Paths:     []string{"s2", "s3"},
								Events:    []string{"s4", "s5", "s6"},
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
				Listeners: []Listener{},
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
				Listeners: []Listener{},
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
				Listeners: []Listener{},
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
				Listeners: []Listener{},
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
				Listeners: []Listener{},
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
				Listeners: []Listener{},
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
