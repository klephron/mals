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

	if c.Handlers != nil {
		o.Handlers = make([]*Handler, 0, len(c.Handlers))
		for _, handler := range c.Handlers {
			wire := &Handler{}
			if err := wire.Wire(handler); err != nil {
				return err
			}
			o.Handlers = append(o.Handlers, wire)
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

	if o.Handlers != nil {
		c.Handlers = make([]*config.Handler, 0, len(o.Handlers))
		for _, handler := range o.Handlers {
			unwired, err := handler.Unwire()
			if err != nil {
				return nil, err
			}
			c.Handlers = append(c.Handlers, unwired)
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

	switch c.Level {
	case config.LogLevelError:
		o.Level = LogLevelError
	case config.LogLevelWarn:
		o.Level = LogLevelWarn
	case config.LogLevelInfo:
		o.Level = LogLevelInfo
	case config.LogLevelDebug:
		o.Level = LogLevelDebug
	default:
		return fmt.Errorf("unknown log level")
	}

	switch k := c.Output.(type) {
	case *config.LogOutputFile:
		o.Output = LogOutput{
			Kind: LogOutputKindFile,
			File: &k.File,
		}
	default:
		return fmt.Errorf("unknown log kind")
	}

	return nil
}

func (o *Log) Unwire() (*config.Log, error) {
	c := &config.Log{
		Name: o.Name,
	}

	switch o.Level {
	case LogLevelError:
		c.Level = config.LogLevelError
	case LogLevelWarn:
		c.Level = config.LogLevelWarn
	case LogLevelInfo:
		c.Level = config.LogLevelInfo
	case LogLevelDebug:
		c.Level = config.LogLevelDebug
	default:
		return nil, fmt.Errorf("unknown log level")
	}

	switch o.Output.Kind {
	case LogOutputKindFile:
		output := &config.LogOutputFile{}
		if o.Output.File != nil {
			output.File = *o.Output.File
		}
		c.Output = output
	default:
		return nil, fmt.Errorf("unknown log kind: %v", o.Output.Kind)
	}

	return c, nil
}

func (o *Model) Wire(c *config.Model) error {
	o.Name = c.Name

	switch s := c.Api.(type) {
	case *config.ModelApiOpenai:
		o.Api = ModelApi{
			Kind:        ModelApiKindOpenai,
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

	switch o.Api.Kind {
	case ModelApiKindOpenai:
		c.Api = &config.ModelApiOpenai{
			Url:         o.Api.Url,
			MaxTokens:   o.Api.MaxTokens,
			Temperature: o.Api.Temperature,
		}
	default:
		c.Api = nil
	}

	return c, nil
}

func (o *Lsp) Wire(c *config.Lsp) error {
	o.Name = c.Name

	switch s := c.Api.(type) {
	case *config.LspApiStdio:
		o.Api = LspApi{
			Kind: LspApiKindStdio,
			Cmd:  s.Cmd,
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

	switch o.Api.Kind {
	case LspApiKindStdio:
		c.Api = &config.LspApiStdio{
			Cmd: o.Api.Cmd,
		}
	default:
		c.Api = nil
	}

	return c, nil
}

func (o *Listener) Wire(c *config.Listener) error {
	o.Name = c.Name

	switch k := c.Protocol.(type) {
	case *config.ListenerProtocolApi:
		o.Protocol = ListenerProtocol{
			Kind:     ListenerProtocolKindApi,
			Handlers: nil,
		}
	case *config.ListenerProtocolLsp:
		handlers := make([]*ListenerProtocolHandler, 0, len(k.Handlers))
		for _, handler := range k.Handlers {
			wire := &ListenerProtocolHandler{}
			if err := wire.Wire(&handler); err != nil {
				return err
			}
			handlers = append(handlers, wire)
		}

		o.Protocol = ListenerProtocol{
			Kind:     ListenerProtocolKindLsp,
			Handlers: handlers,
		}
	default:
		return fmt.Errorf("unknown listener kind")
	}

	switch i := c.Ipc.(type) {
	case *config.ListenerIpcTcp:
		o.Ipc = ListenerIpc{
			Kind: ListenerIpcKindTcp,
			Port: &i.Port,
		}
	default:
		return fmt.Errorf("unknown listener ipc")
	}

	return nil
}

func (o *Listener) Unwire() (*config.Listener, error) {
	c := &config.Listener{
		Name: o.Name,
	}

	switch o.Protocol.Kind {
	case ListenerProtocolKindApi:
		c.Protocol = &config.ListenerProtocolApi{}
	case ListenerProtocolKindLsp:
		handlers := make([]config.ListenerProtocolLspHandler, 0, len(o.Protocol.Handlers))

		for _, handler := range o.Protocol.Handlers {
			unwired, err := handler.Unwire()
			if err != nil {
				return nil, err
			}
			handlers = append(handlers, *unwired)
		}

		c.Protocol = &config.ListenerProtocolLsp{
			Handlers: handlers,
		}
	default:
		return nil, fmt.Errorf("unknown listener kind: %v", o.Protocol.Kind)
	}

	switch o.Ipc.Kind {
	case ListenerIpcKindTcp:
		tcp := &config.ListenerIpcTcp{}
		if o.Ipc.Port != nil {
			tcp.Port = *o.Ipc.Port
		}
		c.Ipc = tcp
	default:
		return nil, fmt.Errorf("unknown listener ipc: %v", o.Ipc)
	}

	return c, nil
}

func (o *ListenerProtocolHandler) Wire(c *config.ListenerProtocolLspHandler) error {
	o.Name = c.Name

	o.Condition = ListenerProtocolHandlerCondition{
		Filetypes: c.Condition.Filetypes,
		Paths:     c.Condition.Paths,
	}

	o.Handler = c.Handler

	return nil
}

func (o *ListenerProtocolHandler) Unwire() (*config.ListenerProtocolLspHandler, error) {
	c := config.ListenerProtocolLspHandler{
		Name:    o.Name,
		Handler: o.Handler,
	}

	c.Condition = config.ListenerProtocolLspHandlerCondition{
		Filetypes: o.Condition.Filetypes,
		Paths:     o.Condition.Paths,
	}

	return &c, nil
}

func (o *Handler) Wire(c *config.Handler) error {
	// TODO
	return nil
}

func (o *Handler) Unwire() (*config.Handler, error) {
	// TODO
	return nil, nil
}
