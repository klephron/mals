package config

import (
	"encoding/json"
	"fmt"
)

func (o *Config) UnmarshalJSON(data []byte) error {
	var t struct {
		Loggers   []*Log      `json:"loggers"`
		Models    []*Model    `json:"models"`
		Lsps      []*Lsp      `json:"lsps"`
		Usages    []*Usage    `json:"usages"`
		Listeners []*Listener `json:"listeners"`
	}

	t.Loggers = []*Log{}
	t.Models = []*Model{}
	t.Lsps = []*Lsp{}
	t.Usages = []*Usage{}
	t.Listeners = []*Listener{}

	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	o.Loggers = t.Loggers
	o.Models = t.Models
	o.Lsps = t.Lsps
	o.Usages = t.Usages
	o.Listeners = t.Listeners

	return nil
}

func (o *Log) UnmarshalJSON(data []byte) error {
	var kindFile LogKindFile

	var t struct {
		Name  string  `json:"name"`
		Kind  string  `json:"kind"`
		Level string  `json:"level"`
		File  *string `json:"file"`
	}

	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	o.Name = t.Name
	o.Level = t.Level

	switch t.Kind {
	case kindFile.Kind():
		file := &LogKindFile{}
		if t.File != nil {
			file.File = *t.File
		}
		o.Kind = file

	default:
		return fmt.Errorf(`in log: "kind" is unknown, got "%v"`, t.Kind)
	}

	return nil
}

func (o *Model) UnmarshalJSON(data []byte) error {
	t := &struct {
		Name     string          `json:"name"`
		Kind     string          `json:"kind"`
		Settings json.RawMessage `json:"settings"`
	}{}

	if err := json.Unmarshal(data, t); err != nil {
		return err
	}

	o.Name = t.Name

	switch t.Kind {
	case "openai":
		var ts struct {
			Url         string  `json:"url"`
			MaxTokens   int     `json:"max_tokens"`
			Temperature float32 `json:"temperature"`
		}

		if err := json.Unmarshal(t.Settings, &ts); err != nil {
			return err
		}

		o.Settings = &ModelSettingsOpenAI{
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
		Kind     string          `json:"kind"`
		Settings json.RawMessage `json:"settings"`
	}{}

	if err := json.Unmarshal(data, t); err != nil {
		return err
	}

	o.Name = t.Name

	switch t.Kind {
	case "stdio":
		var ts struct {
			Cmd []string `json:"cmd"`
		}

		ts.Cmd = []string{}

		if err := json.Unmarshal(t.Settings, &ts); err != nil {
			return err
		}

		o.Settings = &LspSettingsStdio{
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
		Events     []Event      `json:"events"`
		Conditions []*Condition `json:"conditions"`
		Workflow   *Workflow    `json:"workflow"`
	}

	t.Events = []Event{}

	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	o.Name = t.Name
	o.Events = t.Events
	o.Conditions = t.Conditions
	o.Workflow = t.Workflow

	return nil
}

func (o *Listener) UnmarshalJSON(data []byte) error {
	var kindApi ListenerKindApi
	var kindLsp ListenerKindLsp

	var ipcTcp ListenerIpcTcp
	var ipcStdio ListenerIpcStdio

	var t struct {
		Name   string   `json:"name"`
		Kind   string   `json:"kind"`
		Ipc    string   `json:"ipc"`
		Port   *int     `json:"port"`
		Usages []string `json:"usages"`
	}

	t.Usages = []string{}

	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	o.Name = t.Name

	switch t.Kind {
	case kindApi.Kind():
		api := &ListenerKindApi{}
		o.Kind = api

	case kindLsp.Kind():
		lsp := &ListenerKindLsp{Usages: t.Usages}
		o.Kind = lsp
	default:
		return fmt.Errorf(`in listener: "kind" is unknown, got "%v"`, t.Kind)
	}

	switch t.Ipc {
	case ipcStdio.Ipc():
		stdio := &ListenerIpcStdio{}
		o.Ipc = stdio

	case ipcTcp.Ipc():
		tcp := &ListenerIpcTcp{}

		if t.Port != nil {
			tcp.Port = *t.Port
		}

		o.Ipc = tcp

	default:
		return fmt.Errorf(`in listener: "ipc" is unknown, got "%v"`, t.Ipc)
	}

	return nil
}

func (o *Workflow) UnmarshalJSON(data []byte) error {
	var t struct {
		Steps []*Step `json:"steps"`
	}

	t.Steps = []*Step{}

	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	o.Steps = t.Steps

	return nil
}

func (o *Step) UnmarshalJSON(data []byte) error {
	var t struct {
		Name       string       `json:"name"`
		Conditions []*Condition `json:"conditions"`
		Model      *string      `json:"model"`
		Lsp        *string      `json:"lsp"`
		Scope      string       `json:"scope"`
	}

	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	if t.Model != nil && t.Lsp != nil {
		return fmt.Errorf(`in step %v: both "model" and "lsp" are set`, t.Name)
	}

	if t.Model == nil && t.Lsp == nil {
		return fmt.Errorf(`in step %v: both "model" and "lsp" are not set`, t.Name)
	}

	o.Name = t.Name
	o.Conditions = t.Conditions
	o.Scope = t.Scope

	if t.Model != nil {
		o.Kind = &StepKindModel{Name: *t.Model}
	}

	if t.Lsp != nil {
		o.Kind = &StepKindLsp{Name: *t.Lsp}
	}

	return nil
}

func (o *Condition) UnmarshalJSON(data []byte) error {
	var t struct {
		Filetypes []string `json:"filetypes"`
		Paths     []string `json:"paths"`
	}

	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	o.Filetypes = t.Filetypes
	o.Paths = t.Paths

	return nil
}
