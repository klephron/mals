package message

type InitializeParams struct {
	ClientInfo       ClientInfo        `json:"clientInfo"`
	WorkspaceFolders []WorkspaceFolder `json:"workspaceFolders"`
}

type InitializeRequest struct {
	Request
	Params InitializeParams `json:"params"`
}
