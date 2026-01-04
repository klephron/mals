package config

import (
	"encoding/json"
	"fmt"
)

type WireConfig struct {
	Logs      []*Log      `json:"logs"`
	Models    []*Model    `json:"models"`
	Lsps      []*Lsp      `json:"lsps"`
	Usages    []*Usage    `json:"usages"`
	Listeners []*Listener `json:"listeners"`
}

func (o *Config) UnmarshalJSON(data []byte) error {
	t := WireConfig{}

	t.Logs = []*Log{}
	t.Models = []*Model{}
	t.Lsps = []*Lsp{}
	t.Usages = []*Usage{}
	t.Listeners = []*Listener{}

	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	o.Logs = t.Logs
	o.Models = t.Models
	o.Lsps = t.Lsps
	o.Usages = t.Usages
	o.Listeners = t.Listeners

	return nil
}

func (o *Config) Wire() WireConfig {
	t := WireConfig{
		Logs:      o.Logs,
		Models:    o.Models,
		Lsps:      o.Lsps,
		Usages:    o.Usages,
		Listeners: o.Listeners,
	}
	return t
}

func (o *Config) MarshalJSON() ([]byte, error) {
	t := o.Wire()
	return json.Marshal(t)
}

type WireLog struct {
	Name  string  `json:"name"`
	Kind  string  `json:"kind"`
	Level string  `json:"level"`
	File  *string `json:"file"`
}

func (o *Log) UnmarshalJSON(data []byte) error {
	var kindFile LogKindFile

	t := WireLog{}

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

func (o *Log) Wire() WireLog {
	t := WireLog{
		Name:  o.Name,
		Level: o.Level,
	}

	switch k := o.Kind.(type) {
	case *LogKindFile:
		t.Kind = k.Kind()
		t.File = &k.File
	}

	return t
}

func (o *Log) MarshalJSON() ([]byte, error) {
	t := o.Wire()
	return json.Marshal(t)
}

type WireModel struct {
	Name     string            `json:"name"`
	Kind     string            `json:"kind"`
	Settings WireModelSettings `json:"settings"`
}

type WireModelSettings struct {
	Url         string  `json:"url"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float32 `json:"temperature"`
}

func (o *Model) UnmarshalJSON(data []byte) error {
	var settingsOpenAI ModelSettingsOpenAI

	t := WireModel{}

	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	o.Name = t.Name

	switch t.Kind {
	case settingsOpenAI.Kind():
		o.Settings = &ModelSettingsOpenAI{
			Url:         t.Settings.Url,
			MaxTokens:   t.Settings.MaxTokens,
			Temperature: t.Settings.Temperature,
		}

	default:
		o.Settings = nil
	}

	return nil
}

func (o *Model) Wire() WireModel {
	t := WireModel{
		Name: o.Name,
	}

	switch s := o.Settings.(type) {
	case *ModelSettingsOpenAI:
		t.Kind = s.Kind()
		t.Settings = WireModelSettings{
			Url:         s.Url,
			MaxTokens:   s.MaxTokens,
			Temperature: s.Temperature,
		}
	}

	return t
}

func (o *Model) MarshalJSON() ([]byte, error) {
	t := o.Wire()
	return json.Marshal(t)
}

type WireLsp struct {
	Name     string          `json:"name"`
	Kind     string          `json:"kind"`
	Settings WireLspSettings `json:"settings"`
}

type WireLspSettings struct {
	Cmd []string `json:"cmd"`
}

func (o *Lsp) UnmarshalJSON(data []byte) error {
	var settingsStdio LspSettingsStdio

	t := WireLsp{}

	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	o.Name = t.Name

	switch t.Kind {
	case settingsStdio.Kind():
		o.Settings = &LspSettingsStdio{
			Cmd: t.Settings.Cmd,
		}

	default:
		o.Settings = nil
	}

	return nil
}

func (o *Lsp) Wire() WireLsp {
	t := WireLsp{
		Name: o.Name,
	}

	switch s := o.Settings.(type) {
	case *LspSettingsStdio:
		t.Kind = s.Kind()
		t.Settings = WireLspSettings{
			Cmd: s.Cmd,
		}
	}

	return t
}

func (o *Lsp) MarshalJSON() ([]byte, error) {
	t := o.Wire()
	return json.Marshal(t)
}

type WireUsage struct {
	Name       string       `json:"name"`
	Events     []Event      `json:"events"`
	Conditions []*Condition `json:"conditions"`
	Workflow   *Workflow    `json:"workflow"`
}

func (o *Usage) UnmarshalJSON(data []byte) error {
	t := WireUsage{}

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

func (o *Usage) Wire() WireUsage {
	t := WireUsage{
		Name:       o.Name,
		Events:     o.Events,
		Conditions: o.Conditions,
		Workflow:   o.Workflow,
	}

	return t
}

func (o *Usage) MarshalJSON() ([]byte, error) {
	t := o.Wire()
	return json.Marshal(t)
}

type WireListener struct {
	Name   string   `json:"name"`
	Kind   string   `json:"kind"`
	Ipc    string   `json:"ipc"`
	Port   *int     `json:"port"`
	Usages []string `json:"usages"`
}

func (o *Listener) UnmarshalJSON(data []byte) error {
	var kindApi ListenerKindApi
	var kindLsp ListenerKindLsp

	var ipcStdio ListenerIpcStdio
	var ipcTcp ListenerIpcTcp
	var ipcHttp ListenerIpcHttp

	t := WireListener{}

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

	case ipcHttp.Ipc():
		http := &ListenerIpcHttp{}

		if t.Port != nil {
			http.Port = *t.Port
		}

		o.Ipc = http

	default:
		return fmt.Errorf(`in listener: "ipc" is unknown, got "%v"`, t.Ipc)
	}

	return nil
}

func (o *Listener) Wire() WireListener {
	t := WireListener{
		Name: o.Name,
	}

	switch k := o.Kind.(type) {
	case *ListenerKindApi:
		t.Kind = k.Kind()
	case *ListenerKindLsp:
		t.Kind = k.Kind()
		t.Usages = k.Usages
	}

	switch i := o.Ipc.(type) {
	case *ListenerIpcStdio:
		t.Ipc = i.Ipc()
	case *ListenerIpcTcp:
		t.Ipc = i.Ipc()
		t.Port = &i.Port
	case *ListenerIpcHttp:
		t.Ipc = i.Ipc()
		t.Port = &i.Port
	}

	return t
}

func (o *Listener) MarshalJSON() ([]byte, error) {
	t := o.Wire()
	return json.Marshal(t)
}

type WireWorkflow struct {
	Steps []*Step `json:"steps"`
}

func (o *Workflow) UnmarshalJSON(data []byte) error {
	t := WireWorkflow{}

	t.Steps = []*Step{}

	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	o.Steps = t.Steps

	return nil
}

func (o *Workflow) Wire() WireWorkflow {
	t := WireWorkflow{
		Steps: o.Steps,
	}
	return t
}

func (o *Workflow) MarshalJSON() ([]byte, error) {
	t := o.Wire()
	return json.Marshal(t)
}

type WireStep struct {
	Name       string       `json:"name"`
	Conditions []*Condition `json:"conditions"`
	Model      *string      `json:"model"`
	Lsp        *string      `json:"lsp"`
	Scope      ScopeKind    `json:"scope"`
}

func (o *Step) UnmarshalJSON(data []byte) error {
	t := WireStep{}

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

func (o *Step) Wire() WireStep {
	t := WireStep{
		Name:       o.Name,
		Conditions: o.Conditions,
		Scope:      o.Scope,
	}

	switch k := o.Kind.(type) {
	case *StepKindModel:
		t.Model = &k.Name
	case *StepKindLsp:
		t.Lsp = &k.Name
	}

	return t
}

func (o *Step) MarshalJSON() ([]byte, error) {
	t := o.Wire()
	return json.Marshal(t)
}

type WireCondition struct {
	Filetypes []string `json:"filetypes"`
	Paths     []string `json:"paths"`
}

func (o *Condition) UnmarshalJSON(data []byte) error {
	t := WireCondition{}

	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	o.Filetypes = t.Filetypes
	o.Paths = t.Paths

	return nil
}

func (o *Condition) Wire() WireCondition {
	t := WireCondition{
		Filetypes: o.Filetypes,
		Paths:     o.Paths,
	}
	return t
}

func (o *Condition) MarshalJSON() ([]byte, error) {
	t := o.Wire()
	return json.Marshal(t)
}
