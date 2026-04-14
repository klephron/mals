package wire

import (
	"fmt"
	"mals/pkg/config"
)

func (o *Config) Wire(c *config.Config) error {
	if c.Logs != nil {
		o.Logs = make([]*Log, 0, len(c.Logs))
		for _, log := range c.Logs {
			wire := &Log{}
			if err := wire.Wire(log); err != nil {
				return err
			}
			o.Logs = append(o.Logs, wire)
		}
	}

	if c.Models != nil {
		o.Models = make([]*Model, 0, len(c.Models))
		for _, model := range c.Models {
			wire := &Model{}
			if err := wire.Wire(model); err != nil {
				return err
			}
			o.Models = append(o.Models, wire)
		}
	}

	if c.Lsps != nil {
		o.Lsps = make([]*Lsp, 0, len(c.Lsps))
		for _, lsp := range c.Lsps {
			wire := &Lsp{}
			if err := wire.Wire(lsp); err != nil {
				return err
			}
			o.Lsps = append(o.Lsps, wire)
		}
	}

	if c.Usages != nil {
		o.Usages = make([]*Handler, 0, len(c.Usages))
		for _, usage := range c.Usages {
			wire := &Handler{}
			if err := wire.Wire(usage); err != nil {
				return err
			}
			o.Usages = append(o.Usages, wire)
		}
	}

	if c.Listeners != nil {
		o.Listeners = make([]*Listener, 0, len(c.Listeners))
		for _, listener := range c.Listeners {
			wire := &Listener{}
			if err := wire.Wire(listener); err != nil {
				return err
			}
			o.Listeners = append(o.Listeners, wire)
		}
	}

	return nil
}

func (o *Config) Unwire() (*config.Config, error) {
	c := &config.Config{}

	if o.Logs != nil {
		c.Logs = make([]*config.Log, 0, len(o.Logs))
		for _, log := range o.Logs {
			unwired, err := log.Unwire()
			if err != nil {
				return nil, err
			}
			c.Logs = append(c.Logs, unwired)
		}
	}

	if o.Models != nil {
		c.Models = make([]*config.Model, 0, len(o.Models))
		for _, model := range o.Models {
			unwired, err := model.Unwire()
			if err != nil {
				return nil, err
			}
			c.Models = append(c.Models, unwired)
		}
	}

	if o.Lsps != nil {
		c.Lsps = make([]*config.Lsp, 0, len(o.Lsps))
		for _, lsp := range o.Lsps {
			unwired, err := lsp.Unwire()
			if err != nil {
				return nil, err
			}
			c.Lsps = append(c.Lsps, unwired)
		}
	}

	if o.Usages != nil {
		c.Usages = make([]*config.Usage, 0, len(o.Usages))
		for _, usage := range o.Usages {
			unwired, err := usage.Unwire()
			if err != nil {
				return nil, err
			}
			c.Usages = append(c.Usages, unwired)
		}
	}

	if o.Listeners != nil {
		c.Listeners = make([]*config.Listener, 0, len(o.Listeners))
		for _, listener := range o.Listeners {
			unwired, err := listener.Unwire()
			if err != nil {
				return nil, err
			}
			c.Listeners = append(c.Listeners, unwired)
		}
	}

	return c, nil
}

func (o *Log) Wire(c *config.Log) error {
	o.Name = c.Name
	o.Level = c.Level

	switch k := c.Kind.(type) {
	case *config.LogKindFile:
		o.Kind = k.Kind()
		o.File = &k.File
	default:
		return fmt.Errorf("unknown log kind")
	}

	return nil
}

func (o *Log) Unwire() (*config.Log, error) {
	c := &config.Log{
		Name:  o.Name,
		Level: o.Level,
	}

	switch o.Kind {
	case (&config.LogKindFile{}).Kind():
		file := &config.LogKindFile{}
		if o.File != nil {
			file.File = *o.File
		}
		c.Kind = file
	default:
		return nil, fmt.Errorf("unknown log kind: %v", o.Kind)
	}

	return c, nil
}

func (o *Model) Wire(c *config.Model) error {
	o.Name = c.Name

	switch s := c.Settings.(type) {
	case *config.ModelSettingsOpenAI:
		o.Kind = s.Kind()
		o.Settings = ModelSettings{
			Url:         s.Url,
			MaxTokens:   s.MaxTokens,
			Temperature: s.Temperature,
		}
	default:
		return fmt.Errorf("unknown model settings kind")
	}

	return nil
}

func (o *Model) Unwire() (*config.Model, error) {
	c := &config.Model{
		Name: o.Name,
	}

	switch o.Kind {
	case (&config.ModelSettingsOpenAI{}).Kind():
		c.Settings = &config.ModelSettingsOpenAI{
			Url:         o.Settings.Url,
			MaxTokens:   o.Settings.MaxTokens,
			Temperature: o.Settings.Temperature,
		}
	default:
		c.Settings = nil
	}

	return c, nil
}

func (o *Lsp) Wire(c *config.Lsp) error {
	o.Name = c.Name

	switch s := c.Settings.(type) {
	case *config.LspSettingsStdio:
		o.Kind = s.Kind()
		o.Settings = LspSettings{
			Cmd: s.Cmd,
		}
	default:
		return fmt.Errorf("unknown lsp settings kind")
	}

	return nil
}

func (o *Lsp) Unwire() (*config.Lsp, error) {
	c := &config.Lsp{
		Name: o.Name,
	}

	switch o.Kind {
	case (&config.LspSettingsStdio{}).Kind():
		c.Settings = &config.LspSettingsStdio{
			Cmd: o.Settings.Cmd,
		}
	default:
		c.Settings = nil
	}

	return c, nil
}

func (o *Handler) Wire(c *config.Usage) error {
	o.Name = c.Name
	o.Events = c.Events

	if c.Conditions != nil {
		o.Conditions = make([]*Condition, 0, len(c.Conditions))
		for _, cond := range c.Conditions {
			wire := &Condition{}
			if err := wire.Wire(cond); err != nil {
				return err
			}
			o.Conditions = append(o.Conditions, wire)
		}
	}

	if c.Workflow != nil {
		o.Workflow = &Workflow{}
		if err := o.Workflow.Wire(c.Workflow); err != nil {
			return err
		}
	}

	return nil
}

func (o *Handler) Unwire() (*config.Usage, error) {
	c := &config.Usage{
		Name:   o.Name,
		Events: o.Events,
	}

	if o.Conditions != nil {
		c.Conditions = make([]*config.Condition, 0, len(o.Conditions))
		for _, cond := range o.Conditions {
			unwired, err := cond.Unwire()
			if err != nil {
				return nil, err
			}
			c.Conditions = append(c.Conditions, unwired)
		}
	}

	if o.Workflow != nil {
		unwired, err := o.Workflow.Unwire()
		if err != nil {
			return nil, err
		}
		c.Workflow = unwired
	}

	return c, nil
}

func (o *Listener) Wire(c *config.Listener) error {
	o.Name = c.Name

	switch k := c.Kind.(type) {
	case *config.ListenerKindApi:
		o.Kind = k.Kind()
	case *config.ListenerKindLsp:
		o.Kind = k.Kind()
		o.Usages = k.Usages
	default:
		return fmt.Errorf("unknown listener kind")
	}

	switch i := c.Ipc.(type) {
	case *config.ListenerIpcStdio:
		o.Ipc = i.Ipc()
	case *config.ListenerIpcTcp:
		o.Ipc = i.Ipc()
		o.Port = &i.Port
	case *config.ListenerIpcHttp:
		o.Ipc = i.Ipc()
		o.Port = &i.Port
	default:
		return fmt.Errorf("unknown listener ipc")
	}

	return nil
}

func (o *Listener) Unwire() (*config.Listener, error) {
	c := &config.Listener{
		Name: o.Name,
	}

	switch o.Kind {
	case (&config.ListenerKindApi{}).Kind():
		c.Kind = &config.ListenerKindApi{}
	case (&config.ListenerKindLsp{}).Kind():
		c.Kind = &config.ListenerKindLsp{Usages: o.Usages}
	default:
		return nil, fmt.Errorf("unknown listener kind: %v", o.Kind)
	}

	switch o.Ipc {
	case (&config.ListenerIpcStdio{}).Ipc():
		c.Ipc = &config.ListenerIpcStdio{}
	case (&config.ListenerIpcTcp{}).Ipc():
		tcp := &config.ListenerIpcTcp{}
		if o.Port != nil {
			tcp.Port = *o.Port
		}
		c.Ipc = tcp
	case (&config.ListenerIpcHttp{}).Ipc():
		http := &config.ListenerIpcHttp{}
		if o.Port != nil {
			http.Port = *o.Port
		}
		c.Ipc = http
	default:
		return nil, fmt.Errorf("unknown listener ipc: %v", o.Ipc)
	}

	return c, nil
}

func (o *Workflow) Wire(c *config.Workflow) error {
	if c.Steps != nil {
		o.Steps = make([]*Step, 0, len(c.Steps))
		for _, step := range c.Steps {
			wire := &Step{}
			if err := wire.Wire(step); err != nil {
				return err
			}
			o.Steps = append(o.Steps, wire)
		}
	}
	return nil
}

func (o *Workflow) Unwire() (*config.Workflow, error) {
	c := &config.Workflow{}

	if o.Steps != nil {
		c.Steps = make([]*config.Step, 0, len(o.Steps))
		for _, step := range o.Steps {
			unwired, err := step.Unwire()
			if err != nil {
				return nil, err
			}
			c.Steps = append(c.Steps, unwired)
		}
	}

	return c, nil
}

func (o *Step) Wire(c *config.Step) error {
	o.Name = c.Name
	o.Scope = c.Scope

	if c.Conditions != nil {
		o.Conditions = make([]*Condition, 0, len(c.Conditions))
		for _, cond := range c.Conditions {
			wire := &Condition{}
			if err := wire.Wire(cond); err != nil {
				return err
			}
			o.Conditions = append(o.Conditions, wire)
		}
	}

	switch k := c.Kind.(type) {
	case *config.StepKindModel:
		o.Model = &k.Name
	case *config.StepKindLsp:
		o.Lsp = &k.Name
	default:
		return fmt.Errorf("unknown step kind")
	}

	return nil
}

func (o *Step) Unwire() (*config.Step, error) {
	c := &config.Step{
		Name:  o.Name,
		Scope: o.Scope,
	}

	if o.Conditions != nil {
		c.Conditions = make([]*config.Condition, 0, len(o.Conditions))
		for _, cond := range o.Conditions {
			unwired, err := cond.Unwire()
			if err != nil {
				return nil, err
			}
			c.Conditions = append(c.Conditions, unwired)
		}
	}

	if o.Model != nil && o.Lsp != nil {
		return nil, fmt.Errorf("both model and lsp are set")
	}

	if o.Model == nil && o.Lsp == nil {
		return nil, fmt.Errorf("none model and lsp are set")
	}

	if o.Model != nil {
		c.Kind = &config.StepKindModel{Name: *o.Model}
	}

	if o.Lsp != nil {
		c.Kind = &config.StepKindLsp{Name: *o.Lsp}
	}

	return c, nil
}

func (o *Condition) Wire(c *config.Condition) error {
	o.Filetypes = c.Filetypes
	o.Paths = c.Paths
	return nil
}

func (o *Condition) Unwire() (*config.Condition, error) {
	c := &config.Condition{
		Filetypes: o.Filetypes,
		Paths:     o.Paths,
	}
	return c, nil
}
