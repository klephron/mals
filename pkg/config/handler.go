package config

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
	HandlerLspResourceScopeGlobal HandlerLspResourceScope = "global"
	HandlerLspResourceScopeClient HandlerLspResourceScope = "client"
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
