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
				Loggers:   []*Log{},
				Models:    []*Model{},
				Lsps:      []*Lsp{},
				Usages:    []*Usage{},
				Listeners: []*Listener{},
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
				Loggers:   []*Log{},
				Models:    []*Model{},
				Lsps:      []*Lsp{},
				Usages:    []*Usage{},
				Listeners: []*Listener{},
			},
			expectedError: false,
		},
		{
			name:  "log",
			input: []byte(`{"loggers": [{"name": "l1", "kind": "file", "file": "logFile", "level": "debug"}], "models":[], "lsps":[], "usages":[]}`),
			expected: &Config{
				Loggers: []*Log{
					{
						Name:  "l1",
						Level: "debug",
						Kind:  &LogKindFile{File: "logFile"},
					},
				},
				Models:    []*Model{},
				Lsps:      []*Lsp{},
				Usages:    []*Usage{},
				Listeners: []*Listener{},
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
				Loggers:   []*Log{},
				Models:    []*Model{},
				Lsps:      []*Lsp{},
				Usages:    []*Usage{},
				Listeners: []*Listener{},
			},
			expectedError: false,
		},
		{
			name:  "listeners",
			input: []byte(`{"listeners": [{"name": "l1", "kind": "api", "ipc": "tcp", "port": 12},{"kind": "lsp", "ipc": "stdio"}], "models":[], "lsps":[], "usages":[]}`),
			expected: &Config{
				Loggers: []*Log{},
				Models:  []*Model{},
				Lsps:    []*Lsp{},
				Usages:  []*Usage{},
				Listeners: []*Listener{
					{
						Name: "l1",
						Kind: &ListenerKindApi{},
						Ipc:  &ListenerIpcTcp{Port: 12},
					},
					{
						Name: "",
						Kind: &ListenerKindLsp{
							Usages: []string{},
						},
						Ipc: &ListenerIpcStdio{},
					},
				},
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
				Loggers: []*Log{},
				Models: []*Model{
					{
						Name: "1",
						Settings: &ModelSettingsOpenAI{
							Url:         "something",
							MaxTokens:   12,
							Temperature: 1,
						},
					},
				},
				Lsps:      []*Lsp{},
				Usages:    []*Usage{},
				Listeners: []*Listener{},
			},
			expectedError: false,
		},
		{
			name:  "model kind OpenAI default",
			input: []byte(`{"models":[{"name": "1", "kind": "openai", "settings": {}}]}`),
			expected: &Config{
				Loggers: []*Log{},
				Models: []*Model{
					{
						Name:     "1",
						Settings: &ModelSettingsOpenAI{},
					},
				},
				Lsps:      []*Lsp{},
				Usages:    []*Usage{},
				Listeners: []*Listener{},
			},
			expectedError: false,
		},
		{
			name:  "lsp kind stdio",
			input: []byte(`{"models":[], "lsps":[{"name": "1", "kind": "stdio", "settings": {"cmd": ["something", "arg1", "arg2"]}}], "usages":[]}`),
			expected: &Config{
				Loggers: []*Log{},
				Models:  []*Model{},
				Lsps: []*Lsp{
					{
						Name: "1",
						Settings: &LspSettingsStdio{
							Cmd: []string{"something", "arg1", "arg2"},
						},
					},
				},
				Usages:    []*Usage{},
				Listeners: []*Listener{},
			},
			expectedError: false,
		},
		{
			name:  "lsp kind stdio empty cmd",
			input: []byte(`{"models":[], "lsps":[{"name": "1", "kind": "stdio", "settings": {"cmd": []}}], "usages":[]}`),
			expected: &Config{
				Loggers: []*Log{},
				Models:  []*Model{},
				Lsps: []*Lsp{
					{
						Name: "1",
						Settings: &LspSettingsStdio{
							Cmd: []string{},
						},
					},
				},
				Usages:    []*Usage{},
				Listeners: []*Listener{},
			},
			expectedError: false,
		},
		{
			name:  "lsp kind stdio nil cmd",
			input: []byte(`{"lsps":[{"name": "1", "kind": "stdio", "settings": {}}]}`),
			expected: &Config{
				Loggers: []*Log{},
				Models:  []*Model{},
				Lsps: []*Lsp{
					{
						Name: "1",
						Settings: &LspSettingsStdio{
							Cmd: []string{},
						},
					},
				},
				Usages:    []*Usage{},
				Listeners: []*Listener{},
			},
			expectedError: false,
		},
		{
			name:  "usage condition nil workflow nil",
			input: []byte(`{"usages":[{"name": "1"}]}`),
			expected: &Config{
				Loggers: []*Log{},
				Models:  []*Model{},
				Lsps:    []*Lsp{},
				Usages: []*Usage{
					{
						Name:       "1",
						Events:     []string{},
						Conditions: []*Condition{},
						Workflow:   nil,
					},
				},
				Listeners: []*Listener{},
			},
			expectedError: false,
		},
		{
			name:  "usage condition empty workflow nil",
			input: []byte(`{"usages":[{"name": "1", "conditions": []}]}`),
			expected: &Config{
				Loggers: []*Log{},
				Models:  []*Model{},
				Lsps:    []*Lsp{},
				Usages: []*Usage{
					{
						Name:       "1",
						Events:     []string{},
						Conditions: []*Condition{},
						Workflow:   nil,
					},
				},
				Listeners: []*Listener{},
			},
			expectedError: false,
		},
		{
			name:  "usage condition empty",
			input: []byte(`{"usages":[{"name": "1", "conditions": [{}]}]}`),
			expected: &Config{
				Loggers: []*Log{},
				Models:  []*Model{},
				Lsps:    []*Lsp{},
				Usages: []*Usage{
					{
						Name:   "1",
						Events: []string{},
						Conditions: []*Condition{
							{
								Filetypes: []string{},
								Paths:     []string{},
							},
						},
						Workflow: nil,
					},
				},
				Listeners: []*Listener{},
			},
			expectedError: false,
		},
		{
			name:  "usage condition",
			input: []byte(`{"usages":[{"name": "1", "conditions": [{"filetypes": ["s1"], "paths": ["s2", "s3"], "events": ["s4", "s5", "s6"]}]}]}`),
			expected: &Config{
				Loggers: []*Log{},
				Models:  []*Model{},
				Lsps:    []*Lsp{},
				Usages: []*Usage{
					{
						Name:   "1",
						Events: []string{},
						Conditions: []*Condition{
							{
								Filetypes: []string{"s1"},
								Paths:     []string{"s2", "s3"},
							},
						},
						Workflow: nil,
					},
				},
				Listeners: []*Listener{},
			},
			expectedError: false,
		},
		{
			name:  "usage workflow empty",
			input: []byte(`{"usages":[{"name": "1", "workflow": {}}]}`),
			expected: &Config{
				Loggers: []*Log{},
				Models:  []*Model{},
				Lsps:    []*Lsp{},
				Usages: []*Usage{
					{
						Name:       "1",
						Events:     []string{},
						Conditions: []*Condition{},
						Workflow: &Workflow{
							Steps: []*Step{},
						},
					},
				},
				Listeners: []*Listener{},
			},
			expectedError: false,
		},
		{
			name:  "usage workflow name",
			input: []byte(`{"usages":[{"workflow": {"name": "w1"}}]}`),
			expected: &Config{
				Loggers: []*Log{},
				Models:  []*Model{},
				Lsps:    []*Lsp{},
				Usages: []*Usage{
					{
						Events:     []string{},
						Conditions: []*Condition{},
						Workflow: &Workflow{
							Steps: []*Step{},
						},
					},
				},
				Listeners: []*Listener{},
			},
			expectedError: false,
		},
		{
			name:  "usage workflow steps empty",
			input: []byte(`{"usages":[{"workflow": {"steps":[]}}]}`),
			expected: &Config{
				Loggers: []*Log{},
				Models:  []*Model{},
				Lsps:    []*Lsp{},
				Usages: []*Usage{
					{
						Events:     []string{},
						Conditions: []*Condition{},
						Workflow: &Workflow{
							Steps: []*Step{},
						},
					},
				},
				Listeners: []*Listener{},
			},
			expectedError: false,
		},
		{
			name:  "usage workflow step lsp empty",
			input: []byte(`{"usages":[{"workflow": {"steps":[{"lsp": "1"}]}}]}`),
			expected: &Config{
				Loggers: []*Log{},
				Models:  []*Model{},
				Lsps:    []*Lsp{},
				Usages: []*Usage{
					{
						Events:     []string{},
						Conditions: []*Condition{},
						Workflow: &Workflow{
							Steps: []*Step{
								{
									Name:       "",
									Conditions: []*Condition{},
									Kind:       &StepKindLsp{Name: "1"},
								},
							},
						},
					},
				},
				Listeners: []*Listener{},
			},
			expectedError: false,
		},
		{
			name:  "usage workflow step lsp",
			input: []byte(`{"usages":[{"workflow": {"steps":[{"name": "1", "conditions": [], "lsp": "l1"}]}}]}`),
			expected: &Config{
				Loggers: []*Log{},
				Models:  []*Model{},
				Lsps:    []*Lsp{},
				Usages: []*Usage{
					{
						Events:     []string{},
						Conditions: []*Condition{},
						Workflow: &Workflow{
							Steps: []*Step{
								{
									Name:       "1",
									Conditions: []*Condition{},
									Kind:       &StepKindLsp{Name: "l1"},
								},
							},
						},
					},
				},
				Listeners: []*Listener{},
			},
			expectedError: false,
		},
		{
			name:  "usage workflow step model",
			input: []byte(`{"usages":[{"workflow": {"steps":[{"name": "1", "conditions": [], "model": "m1"}]}}]}`),
			expected: &Config{
				Loggers: []*Log{},
				Models:  []*Model{},
				Lsps:    []*Lsp{},
				Usages: []*Usage{
					{
						Events:     []string{},
						Conditions: []*Condition{},
						Workflow: &Workflow{
							Steps: []*Step{
								{
									Name:       "1",
									Conditions: []*Condition{},
									Kind:       &StepKindModel{Name: "m1"},
								},
							},
						},
					},
				},
				Listeners: []*Listener{},
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

			if !bytes.Equal(configJson, expectedJson) {
				t.Errorf("expected %v, got %v", string(expectedJson), string(configJson))
				return
			}
		})
	}
}
