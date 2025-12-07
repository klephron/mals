package config

import (
	"encoding/json"
	"fmt"
)

func (o *Config) UnmarshalJSON(data []byte) error {
	var t struct {
		Loggers   []*json.RawMessage `json:"loggers"`
		Listeners []*json.RawMessage `json:"listeners"`
		Models    []*Model           `json:"models"`
		Lsps      []*Lsp             `json:"lsps"`
		Usages    []*Usage           `json:"usages"`
	}

	t.Loggers = []*json.RawMessage{}
	t.Listeners = []*json.RawMessage{}
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

	listeners := []Listener{}
	for _, tlistener := range t.Listeners {
		listener, err := listenerUnmarshalJSON(tlistener)
		if err != nil {
			return err
		}
		listeners = append(listeners, listener)
	}
	o.Listeners = listeners

	o.Models = t.Models
	o.Lsps = t.Lsps
	o.Usages = t.Usages

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
		Events    []string `json:"events"`
	}

	t.Filetypes = []string{}
	t.Paths = []string{}
	t.Events = []string{}

	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	o.Filetypes = t.Filetypes
	o.Paths = t.Paths
	o.Events = t.Events

	return nil
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

func logUnmarshalJSON(data *json.RawMessage) (Log, error) {
	var t struct {
		Name  string  `json:"name"`
		Kind  string  `json:"kind"`
		Level string  `json:"level"`
		File  *string `json:"file"`
	}

	if err := json.Unmarshal(*data, &t); err != nil {
		return nil, err
	}

	switch t.Kind {
	case "file":
		file := &LogFile{
			LogGeneric: NewLogGeneric(t.Name),
			Level:      t.Level,
		}
		if t.File != nil {
			file.File = *t.File
		}
		return file, nil

	default:
		return nil, fmt.Errorf(`in log: "kind" is not or not "file", got "%v"`, t.Kind)
	}
}

func listenerUnmarshalJSON(data *json.RawMessage) (Listener, error) {
	var t struct {
		Name string `json:"name"`
		Kind string `json:"kind"`
		Ipc  string `json:"ipc"`
		Port *int   `json:"port"`
	}

	if err := json.Unmarshal(*data, &t); err != nil {
		return nil, err
	}

	switch t.Ipc {
	case "stdio":
		stdio := &ListenerStdio{
			ListenerGeneric: NewListenerGeneric(t.Name, t.Kind),
		}
		return stdio, nil

	case "tcp":
		socket := &ListenerTcp{
			ListenerGeneric: NewListenerGeneric(t.Name, t.Kind),
		}
		if t.Port != nil {
			socket.Port = *t.Port
		}
		return socket, nil

	default:
		return nil, fmt.Errorf(`in listener: "kind" is not or not "stdio"|"socket", got "%v"`, t.Kind)
	}
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
