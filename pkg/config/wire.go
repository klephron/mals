package config

import (
	"encoding/json"
	"fmt"
)

func (o *Config) UnmarshalJSON(data []byte) error {
	var t struct {
		Loggers   []*json.RawMessage `json:"loggers"`
		Listeners []*Listener        `json:"listeners"`
		Models    []*Model           `json:"models"`
		Lsps      []*Lsp             `json:"lsps"`
		Usages    []*Usage           `json:"usages"`
	}

	t.Loggers = []*json.RawMessage{}
	t.Listeners = []*Listener{}
	t.Models = []*Model{}
	t.Lsps = []*Lsp{}
	t.Usages = []*Usage{}

	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	loggers := []Log{}
	for _, tlogger := range t.Loggers {
		logger, err := logUnmarshalJSON(tlogger)
		if err != nil {
			return err
		}
		loggers = append(loggers, logger)
	}
	o.Loggers = loggers

	o.Listeners = t.Listeners
	o.Models = t.Models
	o.Lsps = t.Lsps
	o.Usages = t.Usages

	return nil
}

func logUnmarshalJSON(data *json.RawMessage) (Log, error) {
	var t struct {
		Name  string  `json:"name"`
		Type  string  `json:"type"`
		Level string  `json:"level"`
		File  *string `json:"file"`
	}

	if err := json.Unmarshal(*data, &t); err != nil {
		return nil, err
	}

	switch t.Type {
	case "file":
		file := &LogFile{
			LogGeneric: LogGeneric{
				Name: t.Name,
			},
			Level: t.Level,
		}
		if t.File != nil {
			file.File = *t.File
		}
		return file, nil
	}

	return nil, fmt.Errorf(`in log: "type" is not or not "file", got "%v"`, t.Type)
}

func (o *Listener) UnmarshalJSON(data []byte) error {
	var t struct {
		Name string `json:"name"`
		Type string `json:"type"`
		Port int    `json:"port"`
	}

	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	o.Name = t.Name
	o.Type = t.Type
	o.Port = t.Port

	return nil
}

func (o *Model) UnmarshalJSON(data []byte) error {
	t := &struct {
		Name     string          `json:"name"`
		Spec     string          `json:"spec"`
		Settings json.RawMessage `json:"settings"`
	}{}

	if err := json.Unmarshal(data, t); err != nil {
		return err
	}

	o.Name = t.Name

	switch t.Spec {
	case "openai":
		var ts struct {
			Url         string  `json:"url"`
			MaxTokens   int     `json:"max_tokens"`
			Temperature float32 `json:"temperature"`
		}

		if err := json.Unmarshal(t.Settings, &ts); err != nil {
			return err
		}

		o.Settings = &ModelSpecOpenAI{
			Url:         ts.Url,
			MaxTokens:   ts.MaxTokens,
			Temperature: ts.Temperature,
		}

	default:
		o.Settings = nil
	}

	return nil
}

func (o *Lsp) UnmarshalJSON(data []byte) error {
	t := &struct {
		Name     string          `json:"name"`
		Spec     string          `json:"spec"`
		Settings json.RawMessage `json:"settings"`
	}{}

	if err := json.Unmarshal(data, t); err != nil {
		return err
	}

	o.Name = t.Name

	switch t.Spec {
	case "stdio":
		var ts struct {
			Cmd []string `json:"cmd"`
		}

		ts.Cmd = []string{}

		if err := json.Unmarshal(t.Settings, &ts); err != nil {
			return err
		}

		o.Settings = LspSpecStdio{
			Cmd: ts.Cmd,
		}

	default:
		o.Settings = nil
	}

	return nil
}

func (o *Usage) UnmarshalJSON(data []byte) error {
	var t struct {
		Name       string       `json:"name"`
		Conditions []*Condition `json:"conditions"`
		Workflow   *Workflow    `json:"workflow"`
	}

	t.Conditions = []*Condition{}

	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	o.Name = t.Name
	o.Conditions = t.Conditions
	o.Workflow = t.Workflow

	return nil
}

func (o *Condition) UnmarshalJSON(data []byte) error {
	var t struct {
		Filetypes []string `json:"filetypes"`
		Paths     []string `json:"paths"`
		Types     []string `json:"types"`
	}

	t.Filetypes = []string{}
	t.Paths = []string{}
	t.Types = []string{}

	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	o.Filetypes = t.Filetypes
	o.Paths = t.Paths
	o.Types = t.Types

	return nil
}

func stepUnmarshalJSON(data *json.RawMessage) (Step, error) {
	var t struct {
		Name       string       `json:"name"`
		Conditions []*Condition `json:"conditions"`
		Scope      string       `json:"scope"`
		Model      *string      `json:"model"`
		Template   *string      `json:"template"`
		Lsp        *string      `json:"lsp"`
	}

	t.Conditions = []*Condition{}

	if err := json.Unmarshal(*data, &t); err != nil {
		return nil, err
	}

	if t.Model != nil && t.Lsp != nil {
		return nil, fmt.Errorf(`in step %v: both "model" and "lsp" are set`, t.Name)
	}

	generic := StepGeneric{
		Name:       t.Name,
		Conditions: t.Conditions,
		Scope:      t.Scope,
	}

	if t.Model != nil {
		model := &StepModel{
			StepGeneric: generic,
			Model:       *t.Model,
		}

		if t.Template != nil {
			model.Template = *t.Template
		}

		return model, nil
	}

	if t.Lsp != nil {
		lsp := &StepLsp{
			StepGeneric: generic,
			Lsp:         *t.Lsp,
		}

		if t.Template != nil {
			lsp.Template = *t.Template
		}

		return lsp, nil
	}

	return nil, fmt.Errorf(`in step %v: both "model" and "lsp" are not set`, t.Name)
}

func (o *Workflow) UnmarshalJSON(data []byte) error {
	var t struct {
		Name  string             `json:"name"`
		Steps []*json.RawMessage `json:"steps"`
	}

	t.Steps = []*json.RawMessage{}

	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	o.Name = t.Name

	steps := []Step{}
	for _, tstep := range t.Steps {
		step, err := stepUnmarshalJSON(tstep)
		if err != nil {
			return err
		}
		steps = append(steps, step)
	}
	o.Steps = steps

	return nil
}
