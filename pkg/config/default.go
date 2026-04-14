package config

func DefaultConfig(c *Config) {
	if c.Logs == nil {
		c.Logs = make([]*Log, 0)
	}
	if c.Models == nil {
		c.Models = make([]*Model, 0)
	}
	if c.Lsps == nil {
		c.Lsps = make([]*Lsp, 0)
	}
	for _, lsp := range c.Lsps {
		DefaultLsp(lsp)
	}
	if c.Handlers == nil {
		c.Handlers = make([]*Handler, 0)
	}
	for _, usage := range c.Handlers {
		DefaultHandler(usage)
	}
	if c.Listeners == nil {
		c.Listeners = make([]*Listener, 0)
	}
	for _, listener := range c.Listeners {
		DefaultListener(listener)
	}
}

func DefaultListener(c *Listener) {
	switch kind := c.Protocol.(type) {
	case *ListenerProtocolLsp:
		if kind.Handlers == nil {
			kind.Handlers = make([]ListenerProtocolLspHandler, 0)
		}
	case *ListenerProtocolApi:
	}
}

func DefaultLsp(c *Lsp) {
	switch settings := c.Api.(type) {
	case *LspApiStdio:
		if settings.Cmd == nil {
			settings.Cmd = make([]string, 0)
		}
	}
}

func DefaultHandler(c *Handler) {
}
