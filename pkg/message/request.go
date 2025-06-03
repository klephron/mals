package message

type Request struct {
	Message
	Id     int32  `json:"id"`
	Method string `json:"method"`
}

type InitializeRequest struct {
	Request
	Params InitializeParams `json:"params"`
}
