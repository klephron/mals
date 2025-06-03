package message

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type WorkspaceFolder struct {
	URI  string `json:"uri"`
	Name string `json:"name"`
}

type WorkDoneProgressOptions struct {
}

type CompletionItem struct {
	LabelDetailsSupport bool `json:"labelDetailsSupport"`
}

type CompletionOptions struct {
	WorkDoneProgressOptions
	CompletionItem CompletionItem `json:"completionItem"`
}

type HoverOptions struct {
	WorkDoneProgressOptions
}

type CodeActionOptions struct {
	WorkDoneProgressOptions
	CodeActionKinds []CodeActionKind `json:"codeActionKinds"`
}

type CodeLensOptions struct {
	WorkDoneProgressOptions
}

type TextDocumentSyncKind int
type CodeActionKind string

const (
	NONE        TextDocumentSyncKind = 0
	FULL        TextDocumentSyncKind = 1
	INCREMENTAL TextDocumentSyncKind = 2

	REFACTOR CodeActionKind = "refactor"
	QUICKFIX CodeActionKind = "quickfix"
)

type TextDocumentSyncOptions struct {
	OpenClose bool                 `json:"openClose"`
	Change    TextDocumentSyncKind `json:"change"`
}

type ServerCapabilities struct {
	TextDocumentSync   TextDocumentSyncOptions `json:"textDocumentSync"`
	CompletionProvider CompletionOptions       `json:"completionProvider"`
	HoverProvider      HoverOptions            `json:"hoverProvider"`
	CodeActionProvider CodeActionOptions       `json:"codeActionProvider"`
	CodeLensProvider   CodeLensOptions         `json:"codeLensProvider"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeParams struct {
	ClientInfo       ClientInfo        `json:"clientInfo"`
	WorkspaceFolders []WorkspaceFolder `json:"workspaceFolders"`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   ServerInfo         `json:"serverInfo"`
}
