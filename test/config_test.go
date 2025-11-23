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
				Models: []*Model{},
				Lsps:   []*Lsp{},
				Usages: []*Usage{},
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
				Models: []*Model{},
				Lsps:   []*Lsp{},
				Usages: []*Usage{},
			},
			expectedError: false,
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
