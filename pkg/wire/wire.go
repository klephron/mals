package wire

type Config struct {
	Logs      []*Log      `mapstructure:"logs"`
	Models    []*Model    `mapstructure:"models"`
	Lsps      []*Lsp      `mapstructure:"lsps"`
	Handlers  []*Handler  `mapstructure:"handlers"`
	Listeners []*Listener `mapstructure:"listeners"`
}

type Log struct {
	Name   string     `mapstructure:"name"`
	Level  LogLevel   `mapstructure:"level"`
	Output *LogOutput `mapstructure:"output"`
}

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

type LogOutput struct {
	Kind LogOutputKind `mapstructure:"kind"`
	File *string       `mapstructure:"file"`
}

type LogOutputKind string

const (
	LogOutputKindFile LogOutputKind = "file"
)

type Model struct {
	Name string    `mapstructure:"name"`
	Api  *ModelApi `mapstructure:"api"`
}

type ModelApi struct {
	Kind        ModelApiKind `mapstructure:"kind"`
	Url         *string      `mapstructure:"url"`
	MaxTokens   *int32       `mapstructure:"max_tokens"`
	Temperature *float32     `mapstructure:"temperature"`
}

type ModelApiKind string

const (
	ModelApiKindOpenai ModelApiKind = "openai"
)

type Lsp struct {
	Name string  `mapstructure:"name"`
	Api  *LspApi `mapstructure:"api"`
}

type LspApi struct {
	Kind LspApiKind `mapstructure:"kind"`
	Cmd  []string   `mapstructure:"cmd"`
}

type LspApiKind string

const (
	LspApiKindStdio LspApiKind = "stdio"
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
	HandlerResourceScopeGlobal HandlerResourceScope = "global"
	HandlerResourceScopeClient HandlerResourceScope = "client"
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

type Step map[string]any

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
	Name      string                            `mapstructure:"name"`
	Condition *ListenerProtocolHandlerCondition `mapstructure:"condition"`
	Handler   string                            `mapstructure:"handler"`
}

type ListenerProtocolHandlerCondition struct {
	Filetypes []string `mapstructure:"filetypes"`
	Paths     []string `mapstructure:"paths"`
}
