package jsonrpc

import (
	"encoding/json"
)

type wireMessage struct {
	VersionTag string           `json:"jsonrpc"`
	Id         *int32           `json:"id,omitempty"`
	Method     *string          `json:"method,omitempty"`
	Params     *json.RawMessage `json:"params,omitempty"`
	Result     *json.RawMessage `json:"result,omitempty"`
	Error      *json.RawMessage `json:"error,omitempty"`
}
