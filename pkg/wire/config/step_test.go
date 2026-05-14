package config

import (
	"bytes"
	"encoding/json"
	"mals/internal/util"
	"mals/pkg/config"
	"mals/pkg/core"
	"testing"
)

func TestConfigStepUnwire(t *testing.T) {
	type testExpected struct {
		output *config.Step
		error  bool
	}
	type test struct {
		name     string
		input    *Step
		expected testExpected
	}

	tests := []test{
		{
			name: "name",
			input: &Step{
				"name": "your name",
			},
			expected: testExpected{
				output: &config.Step{
					Name: util.Ptr("your name"),
				},
				error: false,
			},
		},
		{
			name: "assign",
			input: &Step{
				"assign": "clangd",
			},
			expected: testExpected{
				output: &config.Step{
					Assign: util.Ptr("clangd"),
				},
				error: false,
			},
		},
		{
			name: "lsp/completion",
			input: &Step{
				"lsp/completion": map[string]any{
					"resource": "clangd",
				},
			},
			expected: testExpected{
				output: &config.Step{
					Definition: &config.StepLspCompletion{
						Resource: util.Ptr("clangd"),
					},
				},
				error: false,
			},
		},
		{
			name: "json/dumps",
			input: &Step{
				"json/dumps": map[string]any{
					"input": "clangd",
				},
			},
			expected: testExpected{
				output: &config.Step{
					Definition: &config.StepJsonDumps{
						Input: util.Ptr("clangd"),
					},
				},
				error: false,
			},
		},
		{
			name: "json/parse",
			input: &Step{
				"json/parse": map[string]any{
					"input": "clangd",
				},
			},
			expected: testExpected{
				output: &config.Step{
					Definition: &config.StepJsonParse{
						Input: util.Ptr("clangd"),
					},
				},
				error: false,
			},
		},
		{
			name: "json/parse/completion",
			input: &Step{
				"json/parse/completion": map[string]any{
					"input": "clangd",
				},
			},
			expected: testExpected{
				output: &config.Step{
					Definition: &config.StepJsonParseCompletion{
						Input: util.Ptr("clangd"),
					},
				},
				error: false,
			},
		},
		{
			name: "model/raw",
			input: &Step{
				"model/raw": map[string]any{
					"resource": "r1",
					"prompt":   "p1",
				},
			},
			expected: testExpected{
				output: &config.Step{
					Definition: &config.StepModelRaw{
						Resource: util.Ptr("r1"),
						Prompt:   util.Ptr("p1"),
					},
				},
				error: false,
			},
		},
		{
			name: "model",
			input: &Step{
				"model": map[string]any{
					"resource": "r1",
					"prompt":   "p1",
					"schema":   "json/completionItems",
				},
			},
			expected: testExpected{
				output: &config.Step{
					Definition: &config.StepModel{
						Resource: util.Ptr("r1"),
						Prompt:   util.Ptr("p1"),
						Schema:   core.ModelSchemaJsonCompletionItems,
					},
				},
				error: false,
			},
		},
		{
			name: "model with messages",
			input: &Step{
				"model": map[string]any{
					"resource": "r1",
					"messages": []any{
						map[string]any{"role": "system", "content": "You are a helpful assistant."},
						map[string]any{"role": "user", "content": "Hello"},
					},
					"schema": "json/completionItems",
				},
			},
			expected: testExpected{
				output: &config.Step{
					Definition: &config.StepModel{
						Resource: util.Ptr("r1"),
						Messages: []*config.StepModelMessage{
							{Role: core.ModelRoleSystem, Content: util.Ptr("You are a helpful assistant.")},
							{Role: core.ModelRoleUser, Content: util.Ptr("Hello")},
						},
						Schema: core.ModelSchemaJsonCompletionItems,
					},
				},
				error: false,
			},
		},
		{
			name: "return",
			input: &Step{
				"return": "result",
			},
			expected: testExpected{
				output: &config.Step{
					Definition: &config.StepReturn{
						Input: util.Ptr("result"),
					},
				},
				error: false,
			},
		},
		{
			name: "if",
			input: &Step{
				"if": map[string]any{
					"condition": "something",
					"then": []any{
						map[string]any{
							"name": "step1",
						},
						map[string]any{
							"name": "step2",
						},
					},
					"else": []any{
						map[string]any{
							"name": "step3",
						},
					},
				},
			},
			expected: testExpected{
				output: &config.Step{
					Definition: &config.StepIf{
						Condition: util.Ptr("something"),
						Then: []*config.Step{
							{
								Name: util.Ptr("step1"),
							},
							{
								Name: util.Ptr("step2"),
							},
						},
						Else: []*config.Step{
							{
								Name: util.Ptr("step3"),
							},
						},
					},
				},
				error: false,
			},
		},
		{
			name: "if inner",
			input: &Step{
				"if": map[string]any{
					"condition": "something",
					"then": []any{
						map[string]any{
							"if": map[string]any{
								"condition": "something",
								"then": []any{
									map[string]any{
										"name": "step1",
									},
									map[string]any{
										"name": "step2",
									},
								},
								"else": []any{
									map[string]any{
										"name": "step3",
									},
								},
							},
						},
					},
					"else": []any{
						map[string]any{
							"name": "step3",
						},
					},
				},
			},
			expected: testExpected{
				output: &config.Step{
					Definition: &config.StepIf{
						Condition: util.Ptr("something"),
						Then: []*config.Step{
							{
								Definition: &config.StepIf{
									Condition: util.Ptr("something"),
									Then: []*config.Step{
										{
											Name: util.Ptr("step1"),
										},
										{
											Name: util.Ptr("step2"),
										},
									},
									Else: []*config.Step{
										{
											Name: util.Ptr("step3"),
										},
									},
								},
							},
						},
						Else: []*config.Step{
							{
								Name: util.Ptr("step3"),
							},
						},
					},
				},
				error: false,
			},
		},
		{
			name: "for",
			input: &Step{
				"for": map[string]any{
					"condition": "something",
					"do": []any{
						map[string]any{
							"name": "step3",
						},
					},
				},
			},
			expected: testExpected{
				output: &config.Step{
					Definition: &config.StepFor{
						Condition: util.Ptr("something"),
						Do: []*config.Step{
							{
								Name: util.Ptr("step3"),
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
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if output.Definition != nil && tt.expected.output.Definition != nil {
				if output.Definition.StepDefinitionKind() != tt.expected.output.Definition.StepDefinitionKind() {
					t.Fatalf("kind mismatch\nexpected: %v\n output: %v",
						tt.expected.output.Definition.StepDefinitionKind(),
						output.Definition.StepDefinitionKind(),
					)
				}
			}

			outputJSON, err := json.Marshal(output)
			expectedJSON, err := json.Marshal(tt.expected.output)

			if !bytes.Equal(outputJSON, expectedJSON) {
				t.Fatalf("value mismatch\nexpected: %s\n output:  %s", expectedJSON, outputJSON)
			}
		})
	}
}
