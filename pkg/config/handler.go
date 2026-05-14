package config

import "mals/internal/util"

type Handler struct {
	Name string
	Spec HandlerSpec
}

type HandlerSpec interface {
	HandlerSpecKind() string
}

type HandlerSpecLsp struct {
	Resources []*HandlerLspResource
	Endpoints HandlerLspEndpoints
}

func (s *HandlerSpecLsp) HandlerSpecKind() string {
	return "lsp"
}

type HandlerLspResource struct {
	Name  string
	Scope HandlerLspResourceScope
	Spec  HandlerLspResourceSpec
}

type HandlerLspResourceScope string

const (
	HandlerLspResourceScopeGlobal  HandlerLspResourceScope = "global"
	HandlerLspResourceScopeClient  HandlerLspResourceScope = "client"
	HandlerLspResourceScopeHandler HandlerLspResourceScope = "handler"
)

type HandlerLspResourceSpec interface {
	HandlerLspResourceSpecKind() string
}

type HandlerLspResourceSpecLsp struct {
	Name string
}

func (s *HandlerLspResourceSpecLsp) HandlerLspResourceSpecKind() string {
	return "lsp"
}

type HandlerLspResourceSpecModel struct {
	Name string
}

func (s *HandlerLspResourceSpecModel) HandlerLspResourceSpecKind() string {
	return "model"
}

type HandlerLspEndpoints struct {
	Initialize             *HandlerLspEndpointInitialize
	Initialized            *HandlerLspEndpointInitialized
	Shutdown               *HandlerLspEndpointShutdown
	TextDocumentCompletion *HandlerLspEndpointTextDocumentCompletion
	TextDocumentDidChange  *HandlerLspEndpointTextDocumentDidChange
	TextDocumentDidClose   *HandlerLspEndpointTextDocumentDidClose
	TextDocumentDidOpen    *HandlerLspEndpointTextDocumentDidOpen
}

type HandlerLspEndpoint struct {
	Default *bool
}

type HandlerLspEndpointInitialize struct {
	HandlerLspEndpoint
}

type HandlerLspEndpointInitialized struct {
	HandlerLspEndpoint
}

type HandlerLspEndpointShutdown struct {
	HandlerLspEndpoint
}

type HandlerLspEndpointTextDocumentCompletion struct {
	HandlerLspEndpoint
	Execution []*Step
}

type HandlerLspEndpointTextDocumentDidChange struct {
	HandlerLspEndpoint
}

type HandlerLspEndpointTextDocumentDidClose struct {
	HandlerLspEndpoint
}

type HandlerLspEndpointTextDocumentDidOpen struct {
	HandlerLspEndpoint
}

func (s *Handler) Default() {
	switch cs := s.Spec.(type) {
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
			cs.Endpoints.Initialize = &HandlerLspEndpointInitialize{}
		}
		if cs.Endpoints.Initialize.Default == nil {
			cs.Endpoints.Initialize.Default = util.Ptr(true)
		}

		if cs.Endpoints.Initialized == nil {
			cs.Endpoints.Initialized = &HandlerLspEndpointInitialized{}
		}
		if cs.Endpoints.Initialized.Default == nil {
			cs.Endpoints.Initialized.Default = util.Ptr(true)
		}

		if cs.Endpoints.Shutdown == nil {
			cs.Endpoints.Shutdown = &HandlerLspEndpointShutdown{}
		}
		if cs.Endpoints.Shutdown.Default == nil {
			cs.Endpoints.Shutdown.Default = util.Ptr(true)
		}

		if cs.Endpoints.TextDocumentCompletion == nil {
			cs.Endpoints.TextDocumentCompletion = &HandlerLspEndpointTextDocumentCompletion{}
		}
		if cs.Endpoints.TextDocumentCompletion.Execution == nil {
			cs.Endpoints.TextDocumentCompletion.Execution = make([]*Step, 0)
		}
		if cs.Endpoints.TextDocumentCompletion.Default == nil {
			cs.Endpoints.TextDocumentCompletion.Default = util.Ptr(len(cs.Endpoints.TextDocumentCompletion.Execution) == 0)
		}

		if cs.Endpoints.TextDocumentDidChange == nil {
			cs.Endpoints.TextDocumentDidChange = &HandlerLspEndpointTextDocumentDidChange{}
		}
		if cs.Endpoints.TextDocumentDidChange.Default == nil {
			cs.Endpoints.TextDocumentDidChange.Default = util.Ptr(true)
		}

		if cs.Endpoints.TextDocumentDidClose == nil {
			cs.Endpoints.TextDocumentDidClose = &HandlerLspEndpointTextDocumentDidClose{}
		}
		if cs.Endpoints.TextDocumentDidClose.Default == nil {
			cs.Endpoints.TextDocumentDidClose.Default = util.Ptr(true)
		}

		if cs.Endpoints.TextDocumentDidOpen == nil {
			cs.Endpoints.TextDocumentDidOpen = &HandlerLspEndpointTextDocumentDidOpen{}
		}
		if cs.Endpoints.TextDocumentDidOpen.Default == nil {
			cs.Endpoints.TextDocumentDidOpen.Default = util.Ptr(true)
		}
	}
}
