package wire

import (
	"fmt"
	"mals/pkg/config"
)

func (o *Config) Wire(c *config.Config) error {
	if c.Logs != nil {
		o.Logs = make([]Log, len(c.Logs))
		for i, log := range c.Logs {
			wired := Log{}
			if err := wired.Wire(log); err != nil {
				return err
			}
			o.Logs[i] = wired
		}
	}

	if c.Models != nil {
		o.Models = make([]Model, len(c.Models))
		for i, model := range c.Models {
			wired := Model{}
			if err := wired.Wire(model); err != nil {
				return err
			}
			o.Models[i] = wired
		}
	}

	if c.Lsps != nil {
		o.Lsps = make([]Lsp, 0, len(c.Lsps))
		for i, lsp := range c.Lsps {
			wired := Lsp{}
			if err := wired.Wire(lsp); err != nil {
				return err
			}
			o.Lsps[i] = wired
		}
	}

	if c.Handlers != nil {
		o.Handlers = make([]Handler, 0, len(c.Handlers))
		for i, handler := range c.Handlers {
			wired := Handler{}
			if err := wired.Wire(handler); err != nil {
				return err
			}
			o.Handlers[i] = wired
		}
	}

	if c.Listeners != nil {
		o.Listeners = make([]Listener, 0, len(c.Listeners))
		for i, listener := range c.Listeners {
			wired := Listener{}
			if err := wired.Wire(listener); err != nil {
				return err
			}
			o.Listeners[i] = wired
		}
	}

	return nil
}

func (o *Config) Unwire() (*config.Config, error) {
	c := &config.Config{}

	if o.Logs != nil {
		c.Logs = make([]*config.Log, len(o.Logs))
		for i, log := range o.Logs {
			unwired, err := log.Unwire()
			if err != nil {
				return nil, err
			}
			c.Logs[i] = unwired
		}
	}

	if o.Models != nil {
		c.Models = make([]*config.Model, len(o.Models))
		for i, model := range o.Models {
			unwired, err := model.Unwire()
			if err != nil {
				return nil, err
			}
			c.Models[i] = unwired
		}
	}

	if o.Lsps != nil {
		c.Lsps = make([]*config.Lsp, len(o.Lsps))
		for i, lsp := range o.Lsps {
			unwired, err := lsp.Unwire()
			if err != nil {
				return nil, err
			}
			c.Lsps[i] = unwired
		}
	}

	if o.Handlers != nil {
		c.Handlers = make([]*config.Handler, len(o.Handlers))
		for i, handler := range o.Handlers {
			unwired, err := handler.Unwire()
			if err != nil {
				return nil, err
			}
			c.Handlers[i] = unwired
		}
	}

	if o.Listeners != nil {
		c.Listeners = make([]*config.Listener, len(o.Listeners))
		for i, listener := range o.Listeners {
			unwired, err := listener.Unwire()
			if err != nil {
				return nil, err
			}
			c.Listeners[i] = unwired
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
		o.Output = &LogOutput{
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

	switch ca := c.Api.(type) {
	case *config.ModelApiOpenai:
		o.Api = &ModelApi{
			Kind: ModelApiKindOpenai,
		}

		if ca.Url != "" {
			o.Api.Url = &ca.Url
		}
		if ca.MaxTokens != nil {
			o.Api.MaxTokens = ca.MaxTokens
		}
		if ca.Temperature != nil {
			o.Api.Temperature = ca.Temperature
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

	if o.Api != nil {
		switch o.Api.Kind {
		case ModelApiKindOpenai:
			api := &config.ModelApiOpenai{
				MaxTokens:   o.Api.MaxTokens,
				Temperature: o.Api.Temperature,
			}
			if o.Api.Url != nil {
				api.Url = *o.Api.Url
			}
			c.Api = api
		default:
			c.Api = nil
		}
	}

	return c, nil
}

func (o *Lsp) Wire(c *config.Lsp) error {
	o.Name = c.Name

	switch ca := c.Api.(type) {
	case *config.LspApiStdio:
		o.Api = &LspApi{
			Kind: LspApiKindStdio,
			Cmd:  ca.Cmd,
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

	switch cp := c.Protocol.(type) {
	case *config.ListenerProtocolApi:
		o.Protocol = &ListenerProtocol{
			Kind:     ListenerProtocolKindApi,
			Handlers: nil,
		}
	case *config.ListenerProtocolLsp:
		handlers := make([]ListenerProtocolHandler, 0, len(cp.Handlers))
		for _, handler := range cp.Handlers {
			wire := ListenerProtocolHandler{}
			if err := wire.Wire(handler); err != nil {
				return err
			}
			handlers = append(handlers, wire)
		}

		o.Protocol = &ListenerProtocol{
			Kind:     ListenerProtocolKindLsp,
			Handlers: handlers,
		}
	default:
		return fmt.Errorf("unknown listener kind")
	}

	switch i := c.Ipc.(type) {
	case *config.ListenerIpcTcp:
		o.Ipc = &ListenerIpc{
			Kind: ListenerIpcKindTcp,
			Port: i.Port,
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

	if o.Protocol != nil {
		switch o.Protocol.Kind {
		case ListenerProtocolKindApi:
			c.Protocol = &config.ListenerProtocolApi{}
		case ListenerProtocolKindLsp:
			handlers := make([]*config.ListenerProtocolLspHandler, 0, len(o.Protocol.Handlers))

			for _, handler := range o.Protocol.Handlers {
				unwired, err := handler.Unwire()
				if err != nil {
					return nil, err
				}
				handlers = append(handlers, unwired)
			}

			c.Protocol = &config.ListenerProtocolLsp{
				Handlers: handlers,
			}
		default:
			return nil, fmt.Errorf("unknown listener kind: %v", o.Protocol.Kind)
		}
	}

	if o.Ipc != nil {
		switch o.Ipc.Kind {
		case ListenerIpcKindTcp:
			tcp := &config.ListenerIpcTcp{
				Port: o.Ipc.Port,
			}
			c.Ipc = tcp
		default:
			return nil, fmt.Errorf("unknown listener ipc: %v", o.Ipc)
		}
	}

	return c, nil
}

func (o *ListenerProtocolHandler) Wire(c *config.ListenerProtocolLspHandler) error {
	o.Name = c.Name
	o.Handler = c.Handler

	return nil
}

func (o *ListenerProtocolHandler) Unwire() (*config.ListenerProtocolLspHandler, error) {
	c := config.ListenerProtocolLspHandler{
		Name:    o.Name,
		Handler: o.Handler,
	}

	return &c, nil
}

func (o *Handler) Wire(c *config.Handler) error {
	o.Name = c.Name

	switch cs := c.Spec.(type) {
	case *config.HandlerSpecLsp:
		o.Kind = HandlerKindLspCompletion

		o.Resources = make([]HandlerResource, len(cs.Resources))
		for _, resource := range cs.Resources {
			wire := HandlerResource{}
			if err := wire.Wire(resource); err != nil {
				return err
			}
			o.Resources = append(o.Resources, wire)
		}

		o.Endpoints = &HandlerEndpoints{
			Initialize:             &HandlerEndpoint{},
			Initialized:            &HandlerEndpoint{},
			Shutdown:               &HandlerEndpoint{},
			TextDocumentCompletion: &HandlerEndpointCompletion{},
			TextDocumentDidChange:  &HandlerEndpoint{},
			TextDocumentDidClose:   &HandlerEndpoint{},
			TextDocumentDidOpen:    &HandlerEndpoint{},
		}

		if err := o.Endpoints.Initialize.WireInitialize(cs.Endpoints.Initialize); err != nil {
			return err
		}

		if err := o.Endpoints.Initialized.WireInitialized(cs.Endpoints.Initialized); err != nil {
			return err
		}

		if err := o.Endpoints.Shutdown.WireShutdown(cs.Endpoints.Shutdown); err != nil {
			return err
		}

		if err := o.Endpoints.TextDocumentCompletion.Wire(cs.Endpoints.TextDocumentCompletion); err != nil {
			return err
		}

		if err := o.Endpoints.TextDocumentDidChange.WireTextDocumentDidChange(cs.Endpoints.TextDocumentDidChange); err != nil {
			return err
		}

		if err := o.Endpoints.TextDocumentDidClose.WireTextDocumentDidClose(cs.Endpoints.TextDocumentDidClose); err != nil {
			return err
		}

		if err := o.Endpoints.TextDocumentDidOpen.WireTextDocumentDidOpen(cs.Endpoints.TextDocumentDidOpen); err != nil {
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
	case HandlerKindLspCompletion:
		var resources []*config.HandlerLspResource
		if o.Resources != nil {
			resources = make([]*config.HandlerLspResource, len(o.Resources))
			for i, resource := range o.Resources {
				unwired, err := resource.Unwire()
				if err != nil {
					return nil, err
				}
				resources[i] = unwired
			}
		}

		var endpoints config.HandlerLspEndpoints
		if o.Endpoints != nil {
			if o.Endpoints.Initialize != nil {
				endpoint, err := o.Endpoints.Initialize.UnwireInitialize()
				if err != nil {
					return nil, err
				}
				endpoints.Initialize = endpoint
			}

			if o.Endpoints.Initialized != nil {
				endpoint, err := o.Endpoints.Initialized.UnwireInitialized()
				if err != nil {
					return nil, err
				}
				endpoints.Initialized = endpoint
			}

			if o.Endpoints.Shutdown != nil {
				endpoint, err := o.Endpoints.Shutdown.UnwireShutdown()
				if err != nil {
					return nil, err
				}
				endpoints.Shutdown = endpoint
			}

			if o.Endpoints.TextDocumentCompletion != nil {
				endpoint, err := o.Endpoints.TextDocumentCompletion.Unwire()
				if err != nil {
					return nil, err
				}
				endpoints.TextDocumentCompletion = endpoint
			}

			if o.Endpoints.TextDocumentDidChange != nil {
				endpoint, err := o.Endpoints.TextDocumentDidChange.UnwireTextDocumentDidChange()
				if err != nil {
					return nil, err
				}
				endpoints.TextDocumentDidChange = endpoint
			}

			if o.Endpoints.TextDocumentDidClose != nil {
				endpoint, err := o.Endpoints.TextDocumentDidClose.UnwireTextDocumentDidClose()
				if err != nil {
					return nil, err
				}
				endpoints.TextDocumentDidClose = endpoint
			}

			if o.Endpoints.TextDocumentDidOpen != nil {
				endpoint, err := o.Endpoints.TextDocumentDidOpen.UnwireTextDocumentDidOpen()
				if err != nil {
					return nil, err
				}
				endpoints.TextDocumentDidOpen = endpoint
			}
		}

		c.Spec = &config.HandlerSpecLsp{
			Resources: resources,
			Endpoints: endpoints,
		}

	default:
		return nil, fmt.Errorf("unknown handler kind: %v", o.Kind)
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
	case config.HandlerLspResourceScopeHandler:
		o.Scope = HandlerResourceScopeHandler
	default:
		return fmt.Errorf("unknown lsp resource scope: %v", c.Scope)
	}

	switch cs := c.Spec.(type) {
	case *config.HandlerLspResourceSpecLsp:
		o.Lsp = &cs.Name
	case *config.HandlerLspResourceSpecModel:
		o.Model = &cs.Name
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
	case HandlerResourceScopeHandler:
		c.Scope = config.HandlerLspResourceScopeHandler
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

func (o *HandlerEndpointCompletion) Wire(c *config.HandlerLspEndpointTextDocumentCompletion) error {
	if err := o.HandlerEndpoint.Wire(&c.HandlerLspEndpoint); err != nil {
		return err
	}

	var execution []Step
	if c.Execution != nil {
		execution = make([]Step, len(c.Execution))
		for i, step := range c.Execution {
			wired := Step{}
			if err := wired.Wire(step); err != nil {
				return err
			}
			execution[i] = wired
		}
	}
	o.Execution = execution

	return nil
}

func (o *HandlerEndpointCompletion) Unwire() (*config.HandlerLspEndpointTextDocumentCompletion, error) {
	c := config.HandlerLspEndpointTextDocumentCompletion{}

	endpoint, err := o.HandlerEndpoint.Unwire()
	if err != nil {
		return nil, err
	}
	c.HandlerLspEndpoint = *endpoint

	var execution []*config.Step
	if o.Execution != nil {
		execution = make([]*config.Step, len(o.Execution))
		for i, step := range o.Execution {
			unwired, err := step.Unwire()
			if err != nil {
				return nil, err
			}
			execution[i] = unwired
		}
	}
	c.Execution = execution

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
