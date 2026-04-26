package execution

import (
	"mals/internal/util"
	"testing"
)

func env() *ExecutionEnvironment { return New(nil) }

func assertRender(t *testing.T, got *string, err error, want string) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatal("got nil, want non-nil")
	}
	if *got != want {
		t.Fatalf("want %q, got %q", want, *got)
	}
}

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
			tmpl: "{{ .other.path }}",
			data: map[string]any{"other": map[string]any{"path": "/tmp/file.txt"}},
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
			got, err := env().renderString(tt.tmpl, tt.data)
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
			got, err := env().renderBool(tt.tmpl, tt.data)
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
			got, err := env().renderInt(tt.tmpl, tt.data)
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

func TestRenderString_SingleSegment(t *testing.T) {
	got, err := env().renderString(`{{.name}}`, map[string]any{"name": "alice"})
	assertRender(t, got, err, "alice")
}

func TestRenderString_TwoSegments(t *testing.T) {
	mem := map[string]any{"user": map[string]any{"name": "bob"}}
	got, err := env().renderString(`{{.user.name}}`, mem)
	assertRender(t, got, err, "bob")
}

func TestRenderString_ThreeSegments(t *testing.T) {
	mem := map[string]any{
		"a": map[string]any{
			"b": map[string]any{"c": "deep"},
		},
	}
	got, err := env().renderString(`{{.a.b.c}}`, mem)
	assertRender(t, got, err, "deep")
}

func TestRenderString_IndexBuiltin(t *testing.T) {
	mem := map[string]any{"items": []any{"p", "q", "r"}}
	got, err := env().renderString(`{{index .items 2}}`, mem)
	assertRender(t, got, err, "r")
}

func TestRenderString_IndexBuiltinFirstElement(t *testing.T) {
	mem := map[string]any{"items": []any{"only"}}
	got, err := env().renderString(`{{index .items 0}}`, mem)
	assertRender(t, got, err, "only")
}

func TestRenderString_IfTrue(t *testing.T) {
	mem := map[string]any{"flag": true}
	got, err := env().renderString(`{{if .flag}}yes{{end}}`, mem)
	assertRender(t, got, err, "yes")
}

func TestRenderString_IfFalseElse(t *testing.T) {
	mem := map[string]any{"flag": false}
	got, err := env().renderString(`{{if .flag}}yes{{else}}no{{end}}`, mem)
	assertRender(t, got, err, "no")
}

func TestRenderString_IfNestedField(t *testing.T) {
	mem := map[string]any{"user": map[string]any{"active": true}}
	got, err := env().renderString(`{{if .user.active}}ok{{end}}`, mem)
	assertRender(t, got, err, "ok")
}

func TestRenderString_Range(t *testing.T) {
	mem := map[string]any{"items": []any{"a", "b", "c"}}
	got, err := env().renderString(`{{range .items}}{{.}} {{end}}`, mem)
	assertRender(t, got, err, "a b c ")
}

func TestRenderString_RangeElse(t *testing.T) {
	mem := map[string]any{"items": []any{}}
	got, err := env().renderString(`{{range .items}}x{{else}}empty{{end}}`, mem)
	assertRender(t, got, err, "empty")
}

func TestRenderString_DefinedSubTemplate(t *testing.T) {
	tmpl := `{{define "sub"}}sub-content{{end}}{{template "sub" .}}`
	got, err := env().renderString(tmpl, map[string]any{})
	assertRender(t, got, err, "sub-content")
}

func TestRenderString_LiteralOnly(t *testing.T) {
	got, err := env().renderString(`hello world`, map[string]any{})
	assertRender(t, got, err, "hello world")
}

func TestRenderString_EmptyTemplate(t *testing.T) {
	got, err := env().renderString(``, map[string]any{})
	assertRender(t, got, err, "")
}

func TestRenderString_MissingTopLevelKey(t *testing.T) {
	_, err := env().renderString(`{{.missing}}`, map[string]any{})
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRenderString_MissingNestedKey(t *testing.T) {
	mem := map[string]any{"user": map[string]any{}}
	_, err := env().renderString(`{{.user.name}}`, mem)
	if err == nil {
		t.Fatal("expected error for missing nested key")
	}
}

func TestRenderString_UnclosedAction(t *testing.T) {
	_, err := env().renderString(`{{.unclosed`, map[string]any{})
	if err == nil {
		t.Fatal("expected parse error")
	}
}
