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
	if c.Usages == nil {
		c.Usages = make([]*Usage, 0)
	}
	for _, usage := range c.Usages {
		DefaultUsage(usage)
	}
	if c.Listeners == nil {
		c.Listeners = make([]*Listener, 0)
	}
	for _, listener := range c.Listeners {
		DefaultListener(listener)
	}
}

func DefaultListener(c *Listener) {
	switch kind := c.Kind.(type) {
	case *ListenerKindLsp:
		if kind.Usages == nil {
			kind.Usages = make([]string, 0)
		}
	case *ListenerKindApi:
	}
}

func DefaultLsp(c *Lsp) {
	switch settings := c.Settings.(type) {
	case *LspSettingsStdio:
		if settings.Cmd == nil {
			settings.Cmd = make([]string, 0)
		}
	}
}

func DefaultUsage(c *Usage) {
	if c.Events == nil {
		c.Events = make([]Event, 0)
	}
	if c.Workflow != nil {
		DefaultWorkflow(c.Workflow)
	}
}

func DefaultWorkflow(c *Workflow) {
	if c.Steps == nil {
		c.Steps = make([]*Step, 0)
	}
}
