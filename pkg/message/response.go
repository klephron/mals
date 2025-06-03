package message

type Response struct {
	Message
	Id  int32  `json:"id"`
}

type InitializeResponse struct {
	Response
	Result InitializeResult `json:"result"`
}
