package message

type Message struct {
	JsonRpc string `json:"jsonrpc"`
}

type Notification struct {
	Message
	Method string `json:"method"`
}

type Request struct {
	Notification
	Id int32 `json:"id"`
}

type Response struct {
	Message
	Id int32 `json:"id"`
}
