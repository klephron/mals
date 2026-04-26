package execution

import (
	"bytes"
	"fmt"
	"mals/internal/util"
	"strconv"
	"strings"
	"text/template"
)

func renderString(tmpl string, data map[string]any) (*string, error) {
	t, err := template.New("").Option("missingkey=error").Parse(tmpl)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return nil, err
	}
	return util.Ptr(buf.String()), nil
}

func renderBool(tmpl string, data map[string]any) (*bool, error) {
	str, err := renderString(tmpl, data)
	if err != nil {
		return nil, err
	}
	if str == nil {
		return nil, nil
	}
	switch *str {
	case "true":
		return util.Ptr(true), nil
	case "false":
		return util.Ptr(false), nil
	default:
		return nil, fmt.Errorf("cannot convert '%v' to bool", *str)
	}
}

func renderInt(tmpl string, data map[string]any) (*int, error) {
	str, err := renderString(tmpl, data)
	if err != nil {
		return nil, err
	}
	if str == nil {
		return nil, nil
	}
	i, err := strconv.Atoi(strings.TrimSpace(*str))
	if err != nil {
		return nil, fmt.Errorf("cannot convert '%v' to int", *str)
	}
	return util.Ptr(i), nil
}
