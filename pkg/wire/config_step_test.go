package wire

import (
	"bytes"
	"encoding/json"
	"mals/internal/util"
	"mals/pkg/config"
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
			name: "lsp/completion",
			input: &Step{
				"lsp/completion": map[string]any{
					"resource": "clangd",
				},
				"assign": "clangd",
			},
			expected: testExpected{
				output: &config.Step{
					Name:   nil,
					Assign: util.Ptr("clangd"),
					Definition: &config.StepLspCompetion{
						Resource: util.Ptr("clangd"),
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
