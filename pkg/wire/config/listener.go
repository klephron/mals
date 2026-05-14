package config

import (
	"fmt"
	"mals/pkg/config"
)

type Listener struct {
	Name     string            `mapstructure:"name"`
	Ipc      *ListenerIpc      `mapstructure:"ipc"`
	Protocol *ListenerProtocol `mapstructure:"protocol"`
}

type ListenerIpc struct {
	Kind ListenerIpcKind `mapstructure:"kind"`
	Port *int32          `mapstructure:"port"`
}

type ListenerIpcKind string

const (
	ListenerIpcKindTcp ListenerIpcKind = "tcp"
)

type ListenerProtocol struct {
	Kind     ListenerProtocolKind      `mapstructure:"kind"`
	Handlers []ListenerProtocolHandler `mapstructure:"handlers"`
}

type ListenerProtocolKind string

const (
	ListenerProtocolKindLsp ListenerProtocolKind = "lsp"
	ListenerProtocolKindApi ListenerProtocolKind = "api"
)

type ListenerProtocolHandler struct {
	Name    string `mapstructure:"name"`
	Handler string `mapstructure:"handler"`
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
