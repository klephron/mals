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

type CompletionOptions struct {
	WorkDoneProgressOptions
}

type HoverOptions struct {
	WorkDoneProgressOptions
}

type CodeActionKind string

const (
	REFACTOR CodeActionKind = "refactor"
	QUICKFIX CodeActionKind = "quickfix"
)

type CodeActionOptions struct {
	WorkDoneProgressOptions
	CodeActionKinds []CodeActionKind `json:"codeActionKinds"`
}

type CodeLensOptions struct {
	WorkDoneProgressOptions
}

type TextDocumentSyncKind int

const (
	NONE        TextDocumentSyncKind = 0
	FULL        TextDocumentSyncKind = 1
	INCREMENTAL TextDocumentSyncKind = 2
)

type TextDocumentSyncOptions struct {
	OpenClose bool                 `json:"openClose"`
	Change    TextDocumentSyncKind `json:"change"`
}

type ServerCapabilities struct {
	TextDocumentSync   TextDocumentSyncOptions `json:"textDocumentSync"`
	CompletionProvider CompletionOptions       `json:"completionProvider"`
	// HoverProvider      HoverOptions            `json:"hoverProvider"`
	// CodeActionProvider CodeActionOptions       `json:"codeActionProvider"`
	// CodeLensProvider   CodeLensOptions         `json:"codeLensProvider"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type TextDocumentItem struct {
	Uri        string `json:"uri"`
	LanguageId string `json:"languageId"`
	Version    int32  `json:"version"`
	Text       string `json:"text"`
}

type TextDocumentIdentifier struct {
	Uri string `json:"uri"`
}

type VersionedTextDocumentIdentifier struct {
	TextDocumentIdentifier
	Version int32 `json:"version"`
}

type Position struct {
	Line      uint32 `json:"line"`
	Character uint32 `json:"character"`
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type TextDocumentContentChangeEvent struct {
	Range Range  `json:"range"`
	Text  string `json:"text"`
}

type CompletionTriggerKind int32

const (
	INVOKED                CompletionTriggerKind = 1
	CHARACTER              CompletionTriggerKind = 2
	INCOMPLETE_COMPLETIONS CompletionTriggerKind = 3
)

type CompletionContext struct {
	TriggerKind      CompletionTriggerKind `json:"triggerKind"`
	TriggerCharacter string                `json:"triggerCharacter"`
}

type CompletionItem struct {
	Label         string `json:"label"`
	Detail        string `json:"detail"`
	Documentation string `json:"documentation"`
}

type CompletionList struct {
	IsIncomplete bool             `json:"isIncomplete"`
	Items        []CompletionItem `json:"items"`
}
