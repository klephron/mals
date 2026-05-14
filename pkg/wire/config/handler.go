package config

import (
	"fmt"
	"mals/pkg/config"
)

type Handler struct {
	Name      string            `mapstructure:"name"`
	Kind      HandlerKind       `mapstructure:"kind"`
	Resources []HandlerResource `mapstructure:"resources"`
	Endpoints *HandlerEndpoints `mapstructure:"endpoints"`
}

type HandlerKind string

const (
	HandlerKindLspCompletion HandlerKind = "lsp/completion"
)

type HandlerResource struct {
	Name  string               `mapstructure:"name"`
	Model *string              `mapstructure:"model"`
	Lsp   *string              `mapstructure:"lsp"`
	Scope HandlerResourceScope `mapstructure:"scope"`
}

type HandlerResourceScope string

const (
	HandlerResourceScopeGlobal  HandlerResourceScope = "global"
	HandlerResourceScopeClient  HandlerResourceScope = "client"
	HandlerResourceScopeHandler HandlerResourceScope = "handler"
)

type HandlerEndpoints struct {
	Initialize             *HandlerEndpoint           `mapstructure:"initialize"`
	Initialized            *HandlerEndpoint           `mapstructure:"initialized"`
	Shutdown               *HandlerEndpoint           `mapstructure:"shutdown"`
	TextDocumentCompletion *HandlerEndpointCompletion `mapstructure:"textDocument/completion"`
	TextDocumentDidChange  *HandlerEndpoint           `mapstructure:"textDocument/didChange"`
	TextDocumentDidClose   *HandlerEndpoint           `mapstructure:"textDocument/didClose"`
	TextDocumentDidOpen    *HandlerEndpoint           `mapstructure:"textDocument/didOpen"`
}

type HandlerEndpoint struct {
	Default *bool `mapstructure:"default"`
}

type HandlerEndpointCompletion struct {
	HandlerEndpoint `mapstructure:",squash"`
	Execution       []Step `mapstructure:"execution"`
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
