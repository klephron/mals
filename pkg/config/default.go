package config

import "mals/internal/util"

func Default(c *Config) {
	DefaultConfig(c)
}

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
	if c.Handlers == nil {
		c.Handlers = make([]*Handler, 0)
	}
	if c.Listeners == nil {
		c.Listeners = make([]*Listener, 0)
	}

	for _, log := range c.Logs {
		DefaultLog(log)
	}
	for _, model := range c.Models {
		DefaultModel(model)
	}
	for _, lsp := range c.Lsps {
		DefaultLsp(lsp)
	}
	for _, handler := range c.Handlers {
		DefaultHandler(handler)
	}
	for _, listener := range c.Listeners {
		DefaultListener(listener)
	}
}

func DefaultLog(c *Log) {
	if c.Level == "" {
		c.Level = LogLevelInfo
	}
	switch co := c.Output.(type) {
	case *LogOutputFile:
		if co.File == "" {
			co.File = "/dev/stdout"
		}
	}
}

func DefaultModel(c *Model) {
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
	switch cs := c.Spec.(type) {
	case *HandlerSpecLsp:
		if cs.Resources == nil {
			cs.Resources = make([]*HandlerLspResource, 0)
		}
		for _, resource := range cs.Resources {
			if resource.Scope == "" {
				resource.Scope = HandlerLspResourceScopeClient
			}
		}

		if cs.Endpoints.Initialize == nil {
			cs.Endpoints.Initialize = &HandlerLspEndpointInitialize{
				HandlerLspEndpoint: HandlerLspEndpoint{
					Default: util.Ptr(true),
				},
			}
		}

		if cs.Endpoints.Initialized == nil {
			cs.Endpoints.Initialized = &HandlerLspEndpointInitialized{
				HandlerLspEndpoint: HandlerLspEndpoint{
					Default: util.Ptr(true),
				},
			}
		}

		if cs.Endpoints.Shutdown == nil {
			cs.Endpoints.Shutdown = &HandlerLspEndpointShutdown{
				HandlerLspEndpoint: HandlerLspEndpoint{
					Default: util.Ptr(true),
				},
			}
		}

		if cs.Endpoints.TextDocumentCompletion == nil {
			cs.Endpoints.TextDocumentCompletion = &HandlerLspEndpointTextDocumentCompletion{
				HandlerLspEndpoint: HandlerLspEndpoint{
					Default: util.Ptr(true),
				},
			}
		}
		if cs.Endpoints.TextDocumentCompletion.Execution == nil {
			cs.Endpoints.TextDocumentCompletion.Execution = make([]*Step, 0)
		}

		if cs.Endpoints.TextDocumentDidChange == nil {
			cs.Endpoints.TextDocumentDidChange = &HandlerLspEndpointTextDocumentDidChange{
				HandlerLspEndpoint: HandlerLspEndpoint{
					Default: util.Ptr(true),
				},
			}
		}

		if cs.Endpoints.TextDocumentDidClose == nil {
			cs.Endpoints.TextDocumentDidClose = &HandlerLspEndpointTextDocumentDidClose{
				HandlerLspEndpoint: HandlerLspEndpoint{
					Default: util.Ptr(true),
				},
			}
		}

		if cs.Endpoints.TextDocumentDidOpen == nil {
			cs.Endpoints.TextDocumentDidOpen = &HandlerLspEndpointTextDocumentDidOpen{
				HandlerLspEndpoint: HandlerLspEndpoint{
					Default: util.Ptr(true),
				},
			}
		}
	}
}

func DefaultListener(c *Listener) {
	switch cp := c.Protocol.(type) {
	case *ListenerProtocolLsp:
		if cp.Handlers == nil {
			cp.Handlers = make([]ListenerProtocolLspHandler, 0)
		}
	case *ListenerProtocolApi:
	}
}
