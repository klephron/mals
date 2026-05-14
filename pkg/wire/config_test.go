package wire

import (
	"bytes"
	"encoding/json"
	"mals/internal/util"
	"testing"

	. "mals/pkg/wire/config"

	"github.com/spf13/viper"
)

func TestUnmarshal(t *testing.T) {
	type testExpected struct {
		output *Config
		error  bool
	}
	type test struct {
		name     string
		input    []byte
		expected testExpected
	}

	tests := []test{
		{
			name:  "empty",
			input: []byte("{}"),
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
		{
			name:  "empty error",
			input: []byte("{"),
			expected: testExpected{
				output: nil,
				error:  true,
			},
		},
		{
			name: "empty explicit",
			input: []byte(`
---
logs: []
models: []
lsps: []
handlers: []
listeners: []
`),
			expected: testExpected{
				output: &Config{
					Logs:      []Log{},
					Models:    []Model{},
					Lsps:      []Lsp{},
					Handlers:  []Handler{},
					Listeners: []Listener{},
				},
				error: false,
			},
		},
		{
			name: "log",
			input: []byte(`
---
logs:
  - name: first
    level: debug
    output:
      kind: file
      file: /dev/null
  - name: second
    level: error
    output:
      kind: file
      file: "something"
`),
			expected: testExpected{
				output: &Config{
					Logs: []Log{
						{
							Name:  "first",
							Level: LogLevelDebug,
							Output: &LogOutput{
								Kind: LogOutputKindFile,
								File: util.Ptr("/dev/null"),
							},
						},
						{
							Name:  "second",
							Level: LogLevelError,
							Output: &LogOutput{
								Kind: LogOutputKindFile,
								File: util.Ptr("something"),
							},
						},
					},
				},
				error: false,
			},
		},
		{
			name: "model 1",
			input: []byte(`
---
models:
  - name: first
    api:
      kind: openai
      url: "https://localhost:9091"
      max_tokens: 250
      temperature: 0.87
`),
			expected: testExpected{
				output: &Config{
					Models: []Model{
						{
							Name: "first",
							Api: &ModelApi{
								Kind:        ModelApiKindOpenai,
								Url:         util.Ptr("https://localhost:9091"),
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
			name: "model 2",
			input: []byte(`
---
models:
  - name: second
    api:
      kind: openai
      url: "https://localhost:9091"
`),
			expected: testExpected{
				output: &Config{
					Models: []Model{
						{
							Name: "second",
							Api: &ModelApi{
								Kind:        ModelApiKindOpenai,
								Url:         util.Ptr("https://localhost:9091"),
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
			name: "model 3",
			input: []byte(`
---
models:
  - name: third
`),
			expected: testExpected{
				output: &Config{
					Models: []Model{
						{
							Name: "third",
							Api:  nil,
						},
					},
				},
				error: false,
			},
		},
		{
			name: "lsp 1",
			input: []byte(`
---
lsps:
  - name: first
    api:
      kind: stdio
      cmd:
        - arg1
        - arg2
        - arg3
`),
			expected: testExpected{
				output: &Config{
					Lsps: []Lsp{
						{
							Name: "first",
							Api: &LspApi{
								Kind: LspApiKindStdio,
								Cmd:  []string{"arg1", "arg2", "arg3"},
							},
						},
					},
				},
				error: false,
			},
		},
		{
			name: "lsp 2",
			input: []byte(`
---
lsps:
  - name: second
    api:
      kind: stdio
`),
			expected: testExpected{
				output: &Config{
					Lsps: []Lsp{
						{
							Name: "second",
							Api: &LspApi{
								Kind: LspApiKindStdio,
								Cmd:  nil,
							},
						},
					},
				},
				error: false,
			},
		},
		{
			name: "lsp 3",
			input: []byte(`
---
lsps:
  - name: third
`),
			expected: testExpected{
				output: &Config{
					Lsps: []Lsp{
						{
							Name: "third",
							Api:  nil,
						},
					},
				},
				error: false,
			},
		},
		{
			name: "listener api",
			input: []byte(`
---
listeners:
  - name: first
    ipc:
      kind: tcp
      port: 8091
    protocol:
      kind: api
`),
			expected: testExpected{
				output: &Config{
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
				error: false,
			},
		},
		{
			name: "listener lsp 1",
			input: []byte(`
---
listeners:
  - name: first
    protocol:
      kind: lsp
      handlers:
        - name: first_h
`),
			expected: testExpected{
				output: &Config{
					Listeners: []Listener{
						{
							Name: "first",
							Ipc:  nil,
							Protocol: &ListenerProtocol{
								Kind: ListenerProtocolKindLsp,
								Handlers: []ListenerProtocolHandler{
									{
										Name: "first_h",
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
			name: "listener lsp 2",
			input: []byte(`
---
listeners:
  - name: first
    protocol:
      kind: lsp
      handlers:
        - name: first_1
        - name: first_2
`),
			expected: testExpected{
				output: &Config{
					Listeners: []Listener{
						{
							Name: "first",
							Ipc:  nil,
							Protocol: &ListenerProtocol{
								Kind: ListenerProtocolKindLsp,
								Handlers: []ListenerProtocolHandler{
									{
										Name: "first_1",
									},
									{
										Name: "first_2",
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
			name: "listener lsp 3",
			input: []byte(`
---
listeners:
  - name: first
    protocol:
      kind: lsp
      handlers:
        - name: first_1
        - name: first_2
`),
			expected: testExpected{
				output: &Config{
					Listeners: []Listener{
						{
							Name: "first",
							Ipc:  nil,
							Protocol: &ListenerProtocol{
								Kind: ListenerProtocolKindLsp,
								Handlers: []ListenerProtocolHandler{
									{
										Name: "first_1",
									},
									{
										Name: "first_2",
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
			name: "handler",
			input: []byte(`
---
handlers:
  - name: first
    kind: lsp/completion
  - name: second
    kind: lsp/completion
`),
			expected: testExpected{
				output: &Config{
					Handlers: []Handler{
						{
							Name: "first",
							Kind: HandlerKindLspCompletion,
						},
						{
							Name: "second",
							Kind: HandlerKindLspCompletion,
						},
					},
				},
				error: false,
			},
		},
		{
			name: "handler resource",
			input: []byte(`
---
handlers:
  - name: first
    kind: lsp/completion
    resources:
      - name: first_1
        model: model_1
        scope: client
      - name: first_2
        lsp: lsp_2
        scope: global
`),
			expected: testExpected{
				output: &Config{
					Handlers: []Handler{
						{
							Name: "first",
							Kind: HandlerKindLspCompletion,
							Resources: []HandlerResource{
								{
									Name:  "first_1",
									Model: util.Ptr("model_1"),
									Scope: HandlerResourceScopeClient,
								},
								{
									Name:  "first_2",
									Lsp:   util.Ptr("lsp_2"),
									Scope: HandlerResourceScopeGlobal,
								},
							},
						},
					},
				},
				error: false,
			},
		},
		{
			name: "handler resource model/lsp not present",
			input: []byte(`
---
handlers:
  - name: first
    kind: lsp/completion
    resources:
      - name: first_1
        scope: client
`),
			expected: testExpected{
				output: &Config{
					Handlers: []Handler{
						{
							Name: "first",
							Kind: HandlerKindLspCompletion,
							Resources: []HandlerResource{
								{
									Name:  "first_1",
									Scope: HandlerResourceScopeClient,
								},
							},
						},
					},
				},
				error: false,
			},
		},
		{
			name: "handler resource model/lsp both present",
			input: []byte(`
---
handlers:
  - name: first
    kind: lsp/completion
    resources:
      - name: first_1
        model: model_1
        lsp: lsp_1
`),
			expected: testExpected{
				output: &Config{
					Handlers: []Handler{
						{
							Name: "first",
							Kind: HandlerKindLspCompletion,
							Resources: []HandlerResource{
								{
									Name:  "first_1",
									Model: util.Ptr("model_1"),
									Lsp:   util.Ptr("lsp_1"),
								},
							},
						},
					},
				},
				error: false,
			},
		},
		{
			name: "handler endpoints empty",
			input: []byte(`
---
handlers:
  - name: first
    kind: lsp/completion
    endpoints: {}
`),
			expected: testExpected{
				output: &Config{
					Handlers: []Handler{
						{
							Name:      "first",
							Kind:      HandlerKindLspCompletion,
							Endpoints: &HandlerEndpoints{},
						},
					},
				},
				error: false,
			},
		},
		{
			name: "handler endpoints all",
			input: []byte(`
---
handlers:
  - name: first
    kind: lsp/completion
    endpoints:
      initialize: {}
      initialized: {}
      shutdown: {}
      textDocument/completion: {}
      textDocument/didChange: {}
      textDocument/didClose: {}
      textDocument/didOpen: {}
`),
			expected: testExpected{
				output: &Config{
					Handlers: []Handler{
						{
							Name: "first",
							Kind: HandlerKindLspCompletion,
							Endpoints: &HandlerEndpoints{
								Initialize:             &HandlerEndpoint{},
								Initialized:            &HandlerEndpoint{},
								Shutdown:               &HandlerEndpoint{},
								TextDocumentCompletion: &HandlerEndpointCompletion{},
								TextDocumentDidChange:  &HandlerEndpoint{},
								TextDocumentDidClose:   &HandlerEndpoint{},
								TextDocumentDidOpen:    &HandlerEndpoint{},
							},
						},
					},
				},
				error: false,
			},
		},
		{
			name: "handler endpoints default",
			input: []byte(`
---
handlers:
  - name: first
    kind: lsp/completion
    endpoints:
      initialize:
        default: true
      initialized:
        default: false
`),
			expected: testExpected{
				output: &Config{
					Handlers: []Handler{
						{
							Name: "first",
							Kind: HandlerKindLspCompletion,
							Endpoints: &HandlerEndpoints{
								Initialize: &HandlerEndpoint{
									Default: util.Ptr(true),
								},
								Initialized: &HandlerEndpoint{
									Default: util.Ptr(false),
								},
							},
						},
					},
				},
				error: false,
			},
		},
		{
			name: "handler endpoints completion 1",
			input: []byte(`
---
handlers:
  - name: first
    kind: lsp/completion
    endpoints:
      textDocument/completion:
        default: false
`),
			expected: testExpected{
				output: &Config{
					Handlers: []Handler{
						{
							Name: "first",
							Kind: HandlerKindLspCompletion,
							Endpoints: &HandlerEndpoints{
								TextDocumentCompletion: &HandlerEndpointCompletion{
									HandlerEndpoint: HandlerEndpoint{
										Default: util.Ptr(false),
									},
									Execution: nil,
								},
							},
						},
					},
				},
				error: false,
			},
		},
		{
			name: "handler endpoints completion 2",
			input: []byte(`
---
handlers:
  - name: first
    kind: lsp/completion
    endpoints:
      textDocument/completion:
        execution: []
`),
			expected: testExpected{
				output: &Config{
					Handlers: []Handler{
						{
							Name: "first",
							Kind: HandlerKindLspCompletion,
							Endpoints: &HandlerEndpoints{
								TextDocumentCompletion: &HandlerEndpointCompletion{
									Execution: []Step{},
								},
							},
						},
					},
				},
				error: false,
			},
		},
		{
			name: "handler endpoints completion 3",
			input: []byte(`
---
handlers:
  - name: first
    kind: lsp/completion
    endpoints:
      textDocument/completion:
        execution:
          - lsp/completion:
              resource: clangd
            assign: clangd
          - json/dumps:
              input: "{{ clangd.out }}"
            assign: clangd_str
`),
			expected: testExpected{
				output: &Config{
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
				error: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := viper.New()
			v.SetConfigType("yaml")

			if err := v.ReadConfig(bytes.NewBuffer(tt.input)); err != nil {
				if !tt.expected.error {
					t.Fatalf("ReadConfig() unexpected error: %v", err)
				}
				return
			}

			var output Config
			err := v.Unmarshal(&output)

			if tt.expected.error {
				if err == nil {
					t.Fatalf("Unmarshal() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unmarshal() unexpected error: %v", err)
			}

			outputJSON, err := json.Marshal(output)
			if err != nil {
				t.Fatalf("marshal output: %v", err)
			}

			expectedJSON, err := json.Marshal(tt.expected.output)
			if err != nil {
				t.Fatalf("marshal expected: %v", err)
			}

			if !bytes.Equal(outputJSON, expectedJSON) {
				t.Fatalf("mismatch\nexpected: %s\n output:  %s", expectedJSON, outputJSON)
			}
		})
	}
}
