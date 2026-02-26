package client

import (
	"fmt"
	"mals/internal/jsonrpc"
	"mals/internal/lsp/protocol"
)

func errorParseUnexpectedType[T jsonrpc.Message](s *LspClient) {
	var dummy T

	resp := jsonrpc.Response{
		Error: &jsonrpc.Error{
			Code:    int32(protocol.ParseError),
			Message: fmt.Sprintf("message is not of type %T", dummy),
		},
	}

	s.plane.Warnf("%v", resp.Error.Message)
	s.send(&resp)
}
