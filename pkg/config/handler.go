package config

type Handler struct {
	Name string
	Type HandlerType
}

type HandlerType interface {
	Kind() string
}

type HandlerLsp struct {
	Resources []HandlerLspResource
	Endpoints HandlerLspEndpoints
}

type HandlerLspResource interface {
	Kind() string
}

type HandlerLspResourceScope string

const (
	HandlerLspResourceScopeGlobal HandlerLspResourceScope = "global"
	HandlerLspResourceScopeClient HandlerLspResourceScope = "client"
)

type HandlerLspResourceLsp struct {
	Lsp   string
	Scope HandlerLspResourceScope
}

func (s *HandlerLspResourceLsp) Kind() string {
	return "lsp"
}

type HandlerLspResourceModel struct {
	Model string
	Scope HandlerLspResourceScope
}

func (s *HandlerLspResourceModel) Kind() string {
	return "model"
}

type HandlerLspEndpoints struct {
	Initialize             HandlerLspEndpointInitialize
	Initialized            HandlerLspEndpointInitialized
	Shutdown               HandlerLspEndpointShutdown
	TextDocumentCompletion HandlerLspEndpointTextDocumentCompletion
	TextDocumentDidChange  HandlerLspEndpointTextDocumentDidChange
	TextDocumentDidClose   HandlerLspEndpointTextDocumentDidClose
	TextDocumentDidOpen    HandlerLspEndpointTextDocumentDidOpen
}

type HandlerLspEndpoint struct {
	Default bool
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
	Execution []Step
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
