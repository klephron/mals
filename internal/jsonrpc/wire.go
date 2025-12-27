package jsonrpc

import (
	"encoding/json"
	"fmt"
)

type messageWire struct {
	VersionTag string          `json:"jsonrpc"`
	Id         *int32          `json:"id,omitempty"`
	Method     *string         `json:"method,omitempty"`
	Params     json.RawMessage `json:"params,omitempty"`
	Result     json.RawMessage `json:"result,omitempty"`
	Error      *Error          `json:"error,omitempty"`
}

func DecodeMessage(data []byte) (Message, error) {
	body, err := decode(data)
	if err != nil {
		return nil, err
	}

	var wire messageWire

	if err := json.Unmarshal(body, &wire); err != nil {
		return nil, fmt.Errorf("jsonrpc unmarshalling: %w", err)
	}

	if wire.Method == nil {
		if wire.Id == nil {
			return nil, fmt.Errorf("response must have id")
		}
		response := &Response{
			Id: *wire.Id,
		}
		response.Error = wire.Error
		response.Result = wire.Result
		return response, nil
	}

	if wire.Id == nil {
		notification := &Notification{
			Method: *wire.Method,
		}
		notification.Params = wire.Params

		return notification, nil
	}

	request := &Request{
		Id:     *wire.Id,
		Method: *wire.Method,
	}
	request.Params = wire.Params

	return request, nil
}

func EncodeMessage(message Message) ([]byte, error) {
	wire := messageWire{
		VersionTag: "2.0",
	}

	switch message := message.(type) {
	case *Request:
		wire.Id = &message.Id
		wire.Method = &message.Method
		wire.Params = message.Params
	case *Notification:
		wire.Method = &message.Method
		wire.Params = message.Params
	case *Response:
		wire.Id = &message.Id
		wire.Result = message.Result
		wire.Error = message.Error
	default:
		return nil, fmt.Errorf("encode unhandled type %T", message)
	}

	body, err := json.Marshal(wire)
	if err != nil {
		return nil, fmt.Errorf("jsonrpc marshalling: %w", err)
	}

	data, err := encode(body)
	return data, err
}
