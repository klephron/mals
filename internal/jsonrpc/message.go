package jsonrpc

import (
	"encoding/json"
)

type Message interface {
	message()
}

type Request struct {
	Message
	Id     int32
	Method string
	Params json.RawMessage
}

type Notification struct {
	Message
	Method string
	Params json.RawMessage
}

type Response struct {
	Message
	Id     int32
	Result json.RawMessage
	Error  *Error
}

type Error struct {
	Code    int32
	Message string
	Data    json.RawMessage
}
