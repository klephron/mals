package message

type Message struct {
	JsonRpc string `json:"jsonrpc"`
}

type RequestMessage struct {
	Message
	Method string `json:"method"`
}
