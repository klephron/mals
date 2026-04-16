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
	o.Name = c.Name

	switch s := c.Spec.(type) {
	case *config.HandlerSpecLsp:
		o.Kind = HandlerKindLsp

		o.Resources = make([]*HandlerResource, len(s.Resources))
		for _, resource := range s.Resources {
			wire := &HandlerResource{}
			if err := wire.Wire(&resource); err != nil {
				return err
			}
			o.Resources = append(o.Resources, wire)
		}

		o.Endpoints = &HandlerEndpoints{
			Initialize:             HandlerEndpoint{},
			Initialized:            HandlerEndpoint{},
			Shutdown:               HandlerEndpoint{},
			TextDocumentCompletion: HandlerEndpointCompletion{},
			TextDocumentDidChange:  HandlerEndpoint{},
			TextDocumentDidClose:   HandlerEndpoint{},
			TextDocumentDidOpen:    HandlerEndpoint{},
		}

		if err := o.Endpoints.Initialize.WireInitialize(&s.Endpoints.Initialize); err != nil {
			return err
		}

		if err := o.Endpoints.Initialized.WireInitialized(&s.Endpoints.Initialized); err != nil {
			return err
		}

		if err := o.Endpoints.Shutdown.WireShutdown(&s.Endpoints.Shutdown); err != nil {
			return err
		}

		if err := o.Endpoints.TextDocumentCompletion.WireTextDocumentCompletion(&s.Endpoints.TextDocumentCompletion); err != nil {
			return err
		}

		if err := o.Endpoints.TextDocumentDidChange.WireTextDocumentDidChange(&s.Endpoints.TextDocumentDidChange); err != nil {
			return err
		}

		if err := o.Endpoints.TextDocumentDidClose.WireTextDocumentDidClose(&s.Endpoints.TextDocumentDidClose); err != nil {
			return err
		}

		if err := o.Endpoints.TextDocumentDidOpen.WireTextDocumentDidOpen(&s.Endpoints.TextDocumentDidOpen); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown handler kind")
	}

	return nil
}

func (o *Handler) Unwire() (*config.Handler, error) {
	c := config.Handler{
		Name: o.Name,
	}

	switch o.Kind {
	case HandlerKindLsp:
		resources := make([]config.HandlerLspResource, len(o.Resources))
		for _, resource := range o.Resources {
			unwired, err := resource.Unwire()
			if err != nil {
				return nil, err
			}
			resources = append(resources, *unwired)
		}

		initialize, err := o.Endpoints.Initialize.UnwireInitialize()
		if err != nil {
			return nil, err
		}

		initialized, err := o.Endpoints.Initialized.UnwireInitialized()
		if err != nil {
			return nil, err
		}

		shutdown, err := o.Endpoints.Shutdown.UnwireShutdown()
		if err != nil {
			return nil, err
		}

		textDocumentCompletion, err := o.Endpoints.TextDocumentCompletion.UnwireTextDocumentCompletion()
		if err != nil {
			return nil, err
		}

		textDocumentDidChange, err := o.Endpoints.TextDocumentDidChange.UnwireTextDocumentDidChange()
		if err != nil {
			return nil, err
		}

		textDocumentDidClose, err := o.Endpoints.TextDocumentDidClose.UnwireTextDocumentDidClose()
		if err != nil {
			return nil, err
		}

		textDocumentDidOpen, err := o.Endpoints.TextDocumentDidOpen.UnwireTextDocumentDidOpen()
		if err != nil {
			return nil, err
		}

		c.Spec = &config.HandlerSpecLsp{
			Resources: resources,
			Endpoints: config.HandlerLspEndpoints{
				Initialize:             *initialize,
				Initialized:            *initialized,
				Shutdown:               *shutdown,
				TextDocumentCompletion: *textDocumentCompletion,
				TextDocumentDidChange:  *textDocumentDidChange,
				TextDocumentDidClose:   *textDocumentDidClose,
				TextDocumentDidOpen:    *textDocumentDidOpen,
			},
		}

	default:
		return nil, fmt.Errorf("unknown handler kind")
	}

	return &c, nil
}

func (o *HandlerResource) Wire(c *config.HandlerLspResource) error {
	o.Name = c.Name

	switch c.Scope {
	case config.HandlerLspResourceScopeGlobal:
		o.Scope = HandlerResourceScopeGlobal
	case config.HandlerLspResourceScopeClient:
		o.Scope = HandlerResourceScopeClient
	default:
		return fmt.Errorf("unknown lsp resource scope")
	}

	switch s := c.Spec.(type) {
	case *config.HandlerLspResourceSpecLsp:
		o.Lsp = &s.Name
	case *config.HandlerLspResourceSpecModel:
		o.Model = &s.Name
	default:
		return fmt.Errorf("unknown lsp resource spec")
	}
	return nil
}

func (o *HandlerResource) Unwire() (*config.HandlerLspResource, error) {
	c := config.HandlerLspResource{
		Name: o.Name,
	}

	switch o.Scope {
	case HandlerResourceScopeGlobal:
		c.Scope = config.HandlerLspResourceScopeGlobal
	case HandlerResourceScopeClient:
		c.Scope = config.HandlerLspResourceScopeClient
	default:
		return nil, fmt.Errorf("unknown lsp resource scope")
	}

	if o.Lsp == nil && o.Model == nil {
		return nil, fmt.Errorf("both model and lsp are not set")
	}

	if o.Lsp != nil && o.Model != nil {
		return nil, fmt.Errorf("both lsp and model are set")
	}

	if o.Lsp != nil {
		c.Spec = &config.HandlerLspResourceSpecLsp{
			Name: *o.Lsp,
		}
	}

	if o.Model != nil {
		c.Spec = &config.HandlerLspResourceSpecModel{
			Name: *o.Model,
		}
	}

	return &c, nil
}

func (o *HandlerEndpoint) Wire(c *config.HandlerLspEndpoint) error {
	o.Default = c.Default
	return nil
}

func (o *HandlerEndpoint) Unwire() (*config.HandlerLspEndpoint, error) {
	c := config.HandlerLspEndpoint{
		Default: o.Default,
	}
	return &c, nil
}

func (o *HandlerEndpoint) WireInitialize(c *config.HandlerLspEndpointInitialize) error {
	return o.Wire(&c.HandlerLspEndpoint)
}

func (o *HandlerEndpoint) UnwireInitialize() (*config.HandlerLspEndpointInitialize, error) {
	c := config.HandlerLspEndpointInitialize{}
	endpoint, err := o.Unwire()
	if err != nil {
		return nil, err
	}
	c.HandlerLspEndpoint = *endpoint
	return &c, nil
}

func (o *HandlerEndpoint) WireInitialized(c *config.HandlerLspEndpointInitialized) error {
	return o.Wire(&c.HandlerLspEndpoint)
}

func (o *HandlerEndpoint) UnwireInitialized() (*config.HandlerLspEndpointInitialized, error) {
	c := config.HandlerLspEndpointInitialized{}
	endpoint, err := o.Unwire()
	if err != nil {
		return nil, err
	}
	c.HandlerLspEndpoint = *endpoint
	return &c, nil
}

func (o *HandlerEndpoint) WireShutdown(c *config.HandlerLspEndpointShutdown) error {
	return o.Wire(&c.HandlerLspEndpoint)
}

func (o *HandlerEndpoint) UnwireShutdown() (*config.HandlerLspEndpointShutdown, error) {
	c := config.HandlerLspEndpointShutdown{}
	endpoint, err := o.Unwire()
	if err != nil {
		return nil, err
	}
	c.HandlerLspEndpoint = *endpoint
	return &c, nil
}

func (o *HandlerEndpointCompletion) WireTextDocumentCompletion(c *config.HandlerLspEndpointTextDocumentCompletion) error {
	if err := o.HandlerEndpoint.Wire(&c.HandlerLspEndpoint); err != nil {
		return err
	}

	return nil
}

func (o *HandlerEndpoint) UnwireTextDocumentCompletion() (*config.HandlerLspEndpointTextDocumentCompletion, error) {
	c := config.HandlerLspEndpointTextDocumentCompletion{}
	endpoint, err := o.Unwire()
	if err != nil {
		return nil, err
	}
	c.HandlerLspEndpoint = *endpoint

	return &c, nil
}

func (o *HandlerEndpoint) WireTextDocumentDidChange(c *config.HandlerLspEndpointTextDocumentDidChange) error {
	return o.Wire(&c.HandlerLspEndpoint)
}

func (o *HandlerEndpoint) UnwireTextDocumentDidChange() (*config.HandlerLspEndpointTextDocumentDidChange, error) {
	c := config.HandlerLspEndpointTextDocumentDidChange{}
	endpoint, err := o.Unwire()
	if err != nil {
		return nil, err
	}
	c.HandlerLspEndpoint = *endpoint
	return &c, nil
}

func (o *HandlerEndpoint) WireTextDocumentDidClose(c *config.HandlerLspEndpointTextDocumentDidClose) error {
	return o.Wire(&c.HandlerLspEndpoint)
}

func (o *HandlerEndpoint) UnwireTextDocumentDidClose() (*config.HandlerLspEndpointTextDocumentDidClose, error) {
	c := config.HandlerLspEndpointTextDocumentDidClose{}
	endpoint, err := o.Unwire()
	if err != nil {
		return nil, err
	}
	c.HandlerLspEndpoint = *endpoint
	return &c, nil
}

func (o *HandlerEndpoint) WireTextDocumentDidOpen(c *config.HandlerLspEndpointTextDocumentDidOpen) error {
	return o.Wire(&c.HandlerLspEndpoint)
}

func (o *HandlerEndpoint) UnwireTextDocumentDidOpen() (*config.HandlerLspEndpointTextDocumentDidOpen, error) {
	c := config.HandlerLspEndpointTextDocumentDidOpen{}
	endpoint, err := o.Unwire()
	if err != nil {
		return nil, err
	}
	c.HandlerLspEndpoint = *endpoint
	return &c, nil
}
