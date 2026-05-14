package config

import (
	"bytes"
	"encoding/json"
	"mals/internal/util"
	"mals/pkg/config"
	"testing"
)

func TestConfigUnwire(t *testing.T) {
	type testExpected struct {
		output *config.Config
		error  bool
	}
	type test struct {
		name     string
		input    *Config
		expected testExpected
	}

	tests := []test{
		{
			name:  "empty",
			input: &Config{},
			expected: testExpected{
				output: &config.Config{
					Logs:      nil,
					Models:    nil,
					Lsps:      nil,
					Handlers:  nil,
					Listeners: nil,
				},
				error: false,
			},
		},
		{
			name: "log level unknown",
			input: &Config{
				Logs: []Log{
					{
						Name:  "first",
						Level: "unknown",
						Output: &LogOutput{
							Kind: LogOutputKindFile,
						},
					},
				},
			},
			expected: testExpected{
				output: nil,
				error:  true,
			},
		},
		{
			name: "log output nil",
			input: &Config{
				Logs: []Log{
					{
						Name:   "first",
						Level:  "unknown",
						Output: nil,
					},
				},
			},
			expected: testExpected{
				output: nil,
				error:  true,
			},
		},
		{
			name: "log output unknown",
			input: &Config{
				Logs: []Log{
					{
						Name:  "first",
						Level: LogLevelInfo,
						Output: &LogOutput{
							Kind: "unknown",
						},
					},
				},
			},
			expected: testExpected{
				output: nil,
				error:  true,
			},
		},
		{
			name: "log",
			input: &Config{
				Logs: []Log{
					{
						Name:  "first",
						Level: LogLevelInfo,
						Output: &LogOutput{
							Kind: LogOutputKindFile,
							File: nil,
						},
					},
					{
						Name:  "second",
						Level: LogLevelError,
						Output: &LogOutput{
							Kind: LogOutputKindFile,
							File: util.Ptr("/dev/null"),
						},
					},
				},
			},
			expected: testExpected{
				output: &config.Config{
					Logs: []*config.Log{
						{
							Name:  "first",
							Level: config.LogLevelInfo,
							Output: &config.LogOutputFile{
								File: "",
							},
						},
						{
							Name:  "second",
							Level: config.LogLevelError,
							Output: &config.LogOutputFile{
								File: "/dev/null",
							},
						},
					},
				},
				error: false,
			},
		},
		{
			name: "model 1",
			input: &Config{
				Models: []Model{
					{
						Name: "first",
						Api: &ModelApi{
							Kind: ModelApiKindOpenai,
						},
					},
				},
			},
			expected: testExpected{
				output: &config.Config{
					Models: []*config.Model{
						{
							Name: "first",
							Api: &config.ModelApiOpenai{
								Url:         "",
								MaxTokens:   nil,
								Temperature: nil,
							},
						},
					},
				},
				error: false,
			},
		},
		{
			name: "model 2",
			input: &Config{
				Models: []Model{
					{
						Name: "first",
						Api: &ModelApi{
							Kind:        ModelApiKindOpenai,
							Url:         util.Ptr("http://localhost:8091"),
							MaxTokens:   util.Ptr(int32(250)),
							Temperature: util.Ptr(float32(0.87)),
						},
					},
				},
			},
			expected: testExpected{
				output: &config.Config{
					Models: []*config.Model{
						{
							Name: "first",
							Api: &config.ModelApiOpenai{
								Url:         "http://localhost:8091",
								MaxTokens:   util.Ptr(int32(250)),
								Temperature: util.Ptr(float32(0.87)),
							},
						},
					},
				},
				error: false,
			},
		},
		{
			name: "lsp 1",
			input: &Config{
				Lsps: []Lsp{
					{
						Name: "first",
						Api: &LspApi{
							Kind: LspApiKindStdio,
							Cmd:  []string{"a", "b", "c"},
						},
					},
					{
						Name: "second",
						Api: &LspApi{
							Kind: LspApiKindStdio,
							Cmd:  nil,
						},
					},
				},
			},
			expected: testExpected{
				output: &config.Config{
					Lsps: []*config.Lsp{
						{
							Name: "first",
							Api: &config.LspApiStdio{
								Cmd: []string{"a", "b", "c"},
							},
						},
						{
							Name: "second",
							Api: &config.LspApiStdio{
								Cmd: nil,
							},
						},
					},
				},
				error: false,
			},
		},
		{
			name: "lsp 2",
			input: &Config{
				Lsps: []Lsp{
					{
						Name: "first",
						Api: &LspApi{
							Kind: LspApiKindStdio,
							Cmd:  []string{"a", "b", "c"},
						},
					},
					{
						Name: "second",
						Api: &LspApi{
							Kind: LspApiKindStdio,
							Cmd:  nil,
						},
					},
				},
			},
			expected: testExpected{
				output: &config.Config{
					Lsps: []*config.Lsp{
						{
							Name: "first",
							Api: &config.LspApiStdio{
								Cmd: []string{"a", "b", "c"},
							},
						},
						{
							Name: "second",
							Api: &config.LspApiStdio{
								Cmd: nil,
							},
						},
					},
				},
				error: false,
			},
		},
		{
			name: "listener ipc nil",
			input: &Config{
				Listeners: []Listener{
					{
						Name: "first",
						Ipc:  nil,
						Protocol: &ListenerProtocol{
							Kind: ListenerProtocolKindApi,
						},
					},
				},
			},
			expected: testExpected{
				output: &config.Config{
					Listeners: []*config.Listener{
						{
							Name:     "first",
							Ipc:      nil,
							Protocol: &config.ListenerProtocolApi{},
						},
					},
				},
				error: false,
			},
		},
		{
			name: "listener protocol nil",
			input: &Config{
				Listeners: []Listener{
					{
						Name: "first",
						Ipc: &ListenerIpc{
							Kind: ListenerIpcKindTcp,
						},
						Protocol: nil,
					},
				},
			},
			expected: testExpected{
				output: &config.Config{
					Listeners: []*config.Listener{
						{
							Name: "first",
							Ipc: &config.ListenerIpcTcp{
								Port: nil,
							},
							Protocol: nil,
						},
					},
				},
				error: false,
			},
		},
		{
			name: "listener api",
			input: &Config{
				Listeners: []Listener{
					{
						Name: "first",
						Ipc: &ListenerIpc{
							Kind: ListenerIpcKindTcp,
							Port: util.Ptr(int32(8091)),
						},
						Protocol: &ListenerProtocol{
							Kind: ListenerProtocolKindApi,
						},
					},
				},
			},
			expected: testExpected{
				output: &config.Config{
					Listeners: []*config.Listener{
						{
							Name: "first",
							Ipc: &config.ListenerIpcTcp{
								Port: util.Ptr(int32(8091)),
							},
							Protocol: &config.ListenerProtocolApi{},
						},
					},
				},
				error: false,
			},
		},
		{
			name: "listener lsp",
			input: &Config{
				Listeners: []Listener{
					{
						Name: "first",
						Ipc: &ListenerIpc{
							Kind: ListenerIpcKindTcp,
							Port: util.Ptr(int32(8091)),
						},
						Protocol: &ListenerProtocol{
							Kind: ListenerProtocolKindLsp,
							Handlers: []ListenerProtocolHandler{
								{
									Name:    "1",
									Handler: "1_h",
								},
								{
									Name:    "2",
									Handler: "2_h",
								},
							},
						},
					},
				},
			},
			expected: testExpected{
				output: &config.Config{
					Listeners: []*config.Listener{
						{
							Name: "first",
							Ipc: &config.ListenerIpcTcp{
								Port: util.Ptr(int32(8091)),
							},
							Protocol: &config.ListenerProtocolLsp{
								Handlers: []*config.ListenerProtocolLspHandler{
									{
										Name:    "1",
										Handler: "1_h",
									},
									{
										Name:    "2",
										Handler: "2_h",
									},
								},
							},
						},
					},
				},
				error: false,
			},
		},
		{
			name: "handler no kind",
			input: &Config{
				Handlers: []Handler{
					{
						Name: "first",
					},
				},
			},
			expected: testExpected{
				output: nil,
				error:  true,
			},
		},
		{
			name: "handler 1",
			input: &Config{
				Handlers: []Handler{
					{
						Name: "first",
						Kind: HandlerKindLspCompletion,
					},
				},
			},
			expected: testExpected{
				output: &config.Config{
					Handlers: []*config.Handler{
						{
							Name: "first",
							Spec: &config.HandlerSpecLsp{},
						},
					},
				},
				error: false,
			},
		},
		{
			name: "handler resources",
			input: &Config{
				Handlers: []Handler{
					{
						Name: "first",
						Kind: HandlerKindLspCompletion,
						Resources: []HandlerResource{
							{
								Name:  "r1",
								Model: util.Ptr("m1"),
								Scope: HandlerResourceScopeClient,
							},
							{
								Name:  "r2",
								Lsp:   util.Ptr("l1"),
								Scope: HandlerResourceScopeGlobal,
							},
						},
					},
				},
			},
			expected: testExpected{
				output: &config.Config{
					Handlers: []*config.Handler{
						{
							Name: "first",
							Spec: &config.HandlerSpecLsp{
								Resources: []*config.HandlerLspResource{
									{
										Name:  "r1",
										Scope: config.HandlerLspResourceScopeClient,
										Spec: &config.HandlerLspResourceSpecModel{
											Name: "m1",
										},
									},
									{
										Name:  "r2",
										Scope: config.HandlerLspResourceScopeGlobal,
										Spec: &config.HandlerLspResourceSpecLsp{
											Name: "l1",
										},
									},
								},
							},
						},
					},
				},
				error: false,
			},
		},
		{
			name: "handler endpoints",
			input: &Config{
				Handlers: []Handler{
					{
						Name: "first",
						Kind: HandlerKindLspCompletion,
						Endpoints: &HandlerEndpoints{
							Initialize: &HandlerEndpoint{},
							Initialized: &HandlerEndpoint{
								Default: util.Ptr(false),
							},
							Shutdown: &HandlerEndpoint{
								Default: util.Ptr(true),
							},
							TextDocumentCompletion: &HandlerEndpointCompletion{},
							TextDocumentDidChange:  &HandlerEndpoint{},
							TextDocumentDidClose:   &HandlerEndpoint{},
							TextDocumentDidOpen:    &HandlerEndpoint{},
						},
					},
				},
			},
			expected: testExpected{
				output: &config.Config{
					Handlers: []*config.Handler{
						{
							Name: "first",
							Spec: &config.HandlerSpecLsp{
								Endpoints: config.HandlerLspEndpoints{
									Initialize: &config.HandlerLspEndpointInitialize{},
									Initialized: &config.HandlerLspEndpointInitialized{
										HandlerLspEndpoint: config.HandlerLspEndpoint{
											Default: util.Ptr(false),
										},
									},
									Shutdown: &config.HandlerLspEndpointShutdown{
										HandlerLspEndpoint: config.HandlerLspEndpoint{
											Default: util.Ptr(true),
										},
									},
									TextDocumentCompletion: &config.HandlerLspEndpointTextDocumentCompletion{},
									TextDocumentDidChange:  &config.HandlerLspEndpointTextDocumentDidChange{},
									TextDocumentDidClose:   &config.HandlerLspEndpointTextDocumentDidClose{},
									TextDocumentDidOpen:    &config.HandlerLspEndpointTextDocumentDidOpen{},
								},
							},
						},
					},
				},
				error: false,
			},
		},
		{
			name: "handler completion step",
			input: &Config{
				Handlers: []Handler{
					{
						Name: "first",
						Kind: HandlerKindLspCompletion,
						Endpoints: &HandlerEndpoints{
							TextDocumentCompletion: &HandlerEndpointCompletion{
								Execution: []Step{
									{
										"lsp/completion": map[string]any{
											"resource": "clangd",
										},
										"assign": "clangd",
									},
									{
										"json/dumps": map[string]any{
											"input": "{{ clangd.out }}",
										},
										"assign": "clangd_str",
									},
								},
							},
						},
					},
				},
			},
			expected: testExpected{
				output: &config.Config{
					Handlers: []*config.Handler{
						{
							Name: "first",
							Spec: &config.HandlerSpecLsp{
								Endpoints: config.HandlerLspEndpoints{
									TextDocumentCompletion: &config.HandlerLspEndpointTextDocumentCompletion{
										Execution: []*config.Step{
											{
												Name:   nil,
												Assign: util.Ptr("clangd"),
												Definition: &config.StepLspCompletion{
													Resource: util.Ptr("clangd"),
												},
											},
											{
												Name:   nil,
												Assign: util.Ptr("clangd_str"),
												Definition: &config.StepJsonDumps{
													Input: util.Ptr("{{ clangd.out }}"),
												},
											},
										},
									},
								},
							},
						},
					},
				},
				error: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			output, err := tt.input.Unwire()

			if tt.expected.error {
				if err == nil {
					t.Errorf("Unmarshal() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unmarshal() unexpected error: %v", err)
				return
			}

			outputJSON, err := json.Marshal(output)
			expectedJSON, err := json.Marshal(tt.expected.output)

			if !bytes.Equal(outputJSON, expectedJSON) {
				t.Fatalf("mismatch\nexpected: %s\n output:  %s", expectedJSON, outputJSON)
			}
		})
	}
}

func TestConfigWire(t *testing.T) {
	type testExpected struct {
		output *Config
		error  bool
	}
	type test struct {
		name     string
		input    *config.Config
		expected testExpected
	}

	tests := []test{
		{
			name:  "empty",
			input: &config.Config{},
			expected: testExpected{
				output: &Config{
					Logs:      nil,
					Models:    nil,
					Lsps:      nil,
					Handlers:  nil,
					Listeners: nil,
				},
				error: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var output Config
			err := output.Wire(tt.input)

			if tt.expected.error {
				if err == nil {
					t.Errorf("Unmarshal() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unmarshal() unexpected error: %v", err)
				return
			}

			configJson, err := json.Marshal(output)
			expectedJson, err := json.Marshal(tt.expected.output)

			if !bytes.Equal(configJson, expectedJson) {
				t.Errorf("expected %v, got %v", string(expectedJson), string(configJson))
				return
			}
		})
	}
}
