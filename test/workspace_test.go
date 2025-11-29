package test

import (
	. "mals/internal/lsp/workspace"
	"testing"
)

func TestGetLastNChars(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		line     uint32
		char     uint32
		n        uint64
		expected string
	}{
		{
			name:     "Simple single line - get 3 chars",
			content:  "Hello World",
			line:     0,
			char:     5,
			n:        3,
			expected: "llo",
		},
		{
			name:     "Single line - get more chars than available",
			content:  "Hi",
			line:     0,
			char:     2,
			n:        5,
			expected: "Hi",
		},
		{
			name:     "Multi-line - position at start of second line",
			content:  "First line\nSecond line",
			line:     1,
			char:     0,
			n:        5,
			expected: "line\n",
		},
		{
			name:     "Multi-line - position in middle of second line",
			content:  "First line\nSecond line",
			line:     1,
			char:     6,
			n:        8,
			expected: "e\nSecond",
		},
		{
			name:     "Multi-line - get last chars crossing lines",
			content:  "Line 1\nLine 2\nLine 3",
			line:     2,
			char:     3,
			n:        10,
			expected: "Line 2\nLin",
		},
		{
			name:     "Empty content",
			content:  "",
			line:     0,
			char:     0,
			n:        5,
			expected: "",
		},
		{
			name:     "Position at start of content",
			content:  "Test content",
			line:     0,
			char:     0,
			n:        3,
			expected: "",
		},
		{
			name:     "Line beyond content length",
			content:  "Only one line",
			line:     5,
			char:     0,
			n:        5,
			expected: " line",
		},
		{
			name:     "Char beyond line length",
			content:  "Short\nLonger line here",
			line:     0,
			char:     100,
			n:        3,
			expected: "ort",
		},
		{
			name:     "Zero characters requested",
			content:  "Some content",
			line:     0,
			char:     4,
			n:        0,
			expected: "",
		},
		{
			name:     "Multi-line with empty lines",
			content:  "First\n\nThird",
			line:     2,
			char:     2,
			n:        4,
			expected: "\n\nTh",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetLastNChars(tt.content, tt.line, tt.char, tt.n)
			if result != tt.expected {
				t.Errorf("GetLastNChars() = %q, want %q", result, tt.expected)
				t.Errorf("Content length: %d", len(tt.content))
			}
		})
	}
}
