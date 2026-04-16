package wire

type Config struct {
	Logs      []*Log      `json:"logs"`
	Models    []*Model    `json:"models"`
	Lsps      []*Lsp      `json:"lsps"`
	Handlers  []*Handler  `json:"handlers"`
	Listeners []*Listener `json:"listeners"`
}

type Log struct {
	Name   string    `json:"name"`
	Level  LogLevel  `json:"level"`
	Output LogOutput `json:"output"`
}

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

type LogOutput struct {
	Kind LogOutputKind `json:"kind"`
	File *string       `json:"file"`
}

type LogOutputKind string

const (
	LogOutputKindFile LogOutputKind = "file"
)

type Model struct {
	Name string   `json:"name"`
	Api  ModelApi `json:"api"`
}

type ModelApi struct {
	Kind        ModelApiKind `json:"kind"`
	Url         string       `json:"url"`
	MaxTokens   int          `json:"max_tokens"`
	Temperature float32      `json:"temperature"`
}

type ModelApiKind string

const (
	ModelApiKindOpenai ModelApiKind = "openai"
)

type Lsp struct {
	Name string `json:"name"`
	Api  LspApi `json:"api"`
}

type LspApi struct {
	Kind LspApiKind `json:"kind"`
	Cmd  []string   `json:"cmd"`
}

type LspApiKind string

const (
	LspApiKindStdio LspApiKind = "stdio"
)

type Handler struct {
	Name      string             `json:"name"`
	Kind      HandlerKind        `json:"kind"`
	Resources []*HandlerResource `json:"resources"`
	Endpoints *HandlerEndpoints  `json:"endpoints"`
}

type HandlerKind string

const (
	HandlerKindLsp HandlerKind = "lsp"
)

type HandlerResource struct {
	Name  string               `json:"name"`
	Model *string              `json:"model"`
	Lsp   *string              `json:"lsp"`
	Scope HandlerResourceScope `json:"scope"`
}

type HandlerResourceScope string

const (
	HandlerResourceScopeGlobal HandlerResourceScope = "global"
	HandlerResourceScopeClient HandlerResourceScope = "client"
)

type HandlerEndpoints struct {
	Initialize             HandlerEndpoint           `json:"initialize"`
	Initialized            HandlerEndpoint           `json:"initialized"`
	Shutdown               HandlerEndpoint           `json:"shutdown"`
	TextDocumentCompletion HandlerEndpointCompletion `json:"textDocument/completion"`
	TextDocumentDidChange  HandlerEndpoint           `json:"textDocument/didChange"`
	TextDocumentDidClose   HandlerEndpoint           `json:"textDocument/didClose"`
	TextDocumentDidOpen    HandlerEndpoint           `json:"textDocument/didOpen"`
}

type HandlerEndpoint struct {
	Default bool `json:"default"`
}

type HandlerEndpointCompletion struct {
	HandlerEndpoint
	Execution []ExecutionStep `json:"execution"`
}

type ExecutionStep map[string]any

type Listener struct {
	Name     string           `json:"name"`
	Ipc      ListenerIpc      `json:"ipc"`
	Protocol ListenerProtocol `json:"protocol"`
}

type ListenerIpc struct {
	Kind ListenerIpcKind `json:"kind"`
	Port *int            `json:"port"`
}

type ListenerIpcKind string

const (
	ListenerIpcKindTcp ListenerIpcKind = "tcp"
)

type ListenerProtocol struct {
	Kind     ListenerProtocolKind       `json:"kind"`
	Handlers []*ListenerProtocolHandler `json:"handlers"`
}

type ListenerProtocolKind string

const (
	ListenerProtocolKindLsp ListenerProtocolKind = "lsp"
	ListenerProtocolKindApi ListenerProtocolKind = "api"
)

type ListenerProtocolHandler struct {
	Name      string                           `json:"name"`
	Condition ListenerProtocolHandlerCondition `json:"condition"`
	Handler   string                           `json:"handler"`
}

type ListenerProtocolHandlerCondition struct {
	Filetypes []string `json:"filetypes"`
	Paths     []string `json:"paths"`
}
