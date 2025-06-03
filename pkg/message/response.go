package message

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   ServerInfo         `json:"serverInfo"`
}

type InitializeResponse struct {
	Response
	Result InitializeResult `json:"result"`
}

type CompletionResponse struct {
	Response
	Result CompletionList `json:"result"`
}
