package jsonrpc

import (
	"encoding/json"
	"fmt"
)

type wireMessage struct {
	VersionTag string           `json:"jsonrpc"`
	Id         *int32           `json:"id,omitempty"`
	Method     *string          `json:"method,omitempty"`
	Params     *json.RawMessage `json:"params,omitempty"`
	Result     *json.RawMessage `json:"result,omitempty"`
	Error      *json.RawMessage `json:"error,omitempty"`
}

func DecodeMessage(data []byte) (Message, error) {
	body, err := decode(data)
	if err != nil {
		return nil, err
	}

	msg := wireMessage{}

	if err := json.Unmarshal(body, &msg); err != nil {
		return nil, fmt.Errorf("jsonrpc unmarshalling: %w", err)
	}

	if msg.Method == nil {
		if msg.Id == nil {
			return nil, fmt.Errorf("response must have id")
		}
		response := &Response{
			Id: *msg.Id,
		}
		response.Error = msg.Error
		response.Result = msg.Result
		return response, nil
	}

	if msg.Id == nil {
		notification := &Notification{
			Method: *msg.Method,
		}
		notification.Params = msg.Params

		return notification, nil
	}

	request := &Request{
		Id:     *msg.Id,
		Method: *msg.Method,
	}
	request.Params = msg.Params

	return request, nil
}

func EncodeMessage(message Message) ([]byte, error) {
	msg := wireMessage{
		VersionTag: "2.0",
	}

	switch typed := message.(type) {
	case *Request:
		msg.Id = &typed.Id
		msg.Method = &typed.Method
		msg.Params = typed.Params
	case *Notification:
		msg.Method = &typed.Method
		msg.Params = typed.Params
	case *Response:
		msg.Id = &typed.Id
		msg.Result = typed.Result
		msg.Error = typed.Error
	default:
		return nil, fmt.Errorf("encode unhandled type %T", typed)
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("jsonrpc marshalling: %w", err)
	}

	data, err := encode(body)
	return data, err
}
