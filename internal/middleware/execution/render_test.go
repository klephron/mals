package execution

import (
	"mals/internal/util"
	"testing"
)

func TestRenderString(t *testing.T) {
	tests := []struct {
		name    string
		tmpl    string
		data    map[string]any
		want    *string
		wantErr bool
	}{
		{
			name: "simple string",
			tmpl: "{{ .name }}",
			data: map[string]any{"name": "foo"},
			want: util.Ptr("foo"),
		},
		{
			name: "nested map",
			tmpl: "{{ .file.path }}",
			data: map[string]any{"file": map[string]any{"path": "/tmp/file.txt"}},
			want: util.Ptr("/tmp/file.txt"),
		},
		{
			name: "integer value",
			tmpl: "{{ .count }}",
			data: map[string]any{"count": 42},
			want: util.Ptr("42"),
		},
		{
			name: "bool value",
			tmpl: "{{ .enabled }}",
			data: map[string]any{"enabled": true},
			want: util.Ptr("true"),
		},
		{
			name:    "invalid template",
			tmpl:    "{{ .unclosed",
			data:    map[string]any{},
			wantErr: true,
		},
		{
			name:    "missing key renders error",
			tmpl:    "{{ .missing }}",
			data:    map[string]any{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := renderString(tt.tmpl, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("wantErr=%v got err=%v", tt.wantErr, err)
				return
			}
			if tt.want == nil && got != nil {
				t.Errorf("want nil got %v", *got)
				return
			}
			if tt.want != nil && (got == nil || *got != *tt.want) {
				t.Errorf("want %v got %v", *tt.want, got)
			}
		})
	}
}

func TestRenderBool(t *testing.T) {
	tests := []struct {
		name    string
		tmpl    string
		data    map[string]any
		want    *bool
		wantErr bool
	}{
		{
			name: "true string",
			tmpl: "{{ .val }}",
			data: map[string]any{"val": "true"},
			want: util.Ptr(true),
		},
		{
			name: "false string",
			tmpl: "{{ .val }}",
			data: map[string]any{"val": "false"},
			want: util.Ptr(false),
		},
		{
			name: "bool true",
			tmpl: "{{ .val }}",
			data: map[string]any{"val": true},
			want: util.Ptr(true),
		},
		{
			name: "bool false",
			tmpl: "{{ .val }}",
			data: map[string]any{"val": false},
			want: util.Ptr(false),
		},
		{
			name:    "invalid value",
			tmpl:    "{{ .val }}",
			data:    map[string]any{"val": "yes"},
			wantErr: true,
			want:    nil,
		},
		{
			name:    "invalid template",
			tmpl:    "{{ .unclosed",
			data:    map[string]any{},
			wantErr: true,
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := renderBool(tt.tmpl, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("wantErr=%v got err=%v", tt.wantErr, err)
				return
			}
			if tt.want == nil && got != nil {
				t.Errorf("want nil got %v", *got)
				return
			}
			if tt.want != nil && (got == nil || *got != *tt.want) {
				t.Errorf("want %v got %v", *tt.want, got)
			}
		})
	}
}

func TestRenderInt(t *testing.T) {
	tests := []struct {
		name    string
		tmpl    string
		data    map[string]any
		want    *int
		wantErr bool
	}{
		{
			name: "integer value",
			tmpl: "{{ .val }}",
			data: map[string]any{"val": 42},
			want: util.Ptr(42),
		},
		{
			name: "string integer",
			tmpl: "{{ .val }}",
			data: map[string]any{"val": "7"},
			want: util.Ptr(7),
		},
		{
			name: "zero",
			tmpl: "{{ .val }}",
			data: map[string]any{"val": 0},
			want: util.Ptr(0),
		},
		{
			name: "negative",
			tmpl: "{{ .val }}",
			data: map[string]any{"val": -3},
			want: util.Ptr(-3),
		},
		{
			name:    "float string",
			tmpl:    "{{ .val }}",
			data:    map[string]any{"val": "3.14"},
			wantErr: true,
		},
		{
			name:    "non-numeric string",
			tmpl:    "{{ .val }}",
			data:    map[string]any{"val": "abc"},
			wantErr: true,
		},
		{
			name:    "invalid template",
			tmpl:    "{{ .unclosed",
			data:    map[string]any{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := renderInt(tt.tmpl, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("wantErr=%v got err=%v", tt.wantErr, err)
				return
			}
			if tt.want == nil && got != nil {
				t.Errorf("want nil got %v", *got)
				return
			}
			if tt.want != nil && (got == nil || *got != *tt.want) {
				t.Errorf("want %v got %v", *tt.want, got)
			}
		})
	}
}
