package message

type InitializeParams struct {
	ClientInfo       ClientInfo        `json:"clientInfo"`
	WorkspaceFolders []WorkspaceFolder `json:"workspaceFolders"`
}

type InitializeRequest struct {
	Request
	Params InitializeParams `json:"params"`
}

type TextDocumentPositionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"Position"`
}

type WorkDoneProgressParams struct {
}

type PartialResultParams struct {
}

type CompletionParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
	PartialResultParams
	Context CompletionContext `json:"context"`
}

type CompletionRequest struct {
	Request
	Params CompletionParams `json:"params"`
}
