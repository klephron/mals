package lsp

import (
	"mals/internal/info"
	"mals/internal/jsonrpc"
	"mals/internal/lsp/protocol"
)

func (s *ClientLsp) handleInitialize(msg jsonrpc.Message) {
	req, ok := msg.(*jsonrpc.Request)

	if !ok {
		resp := jsonrpc.Response{
			Error: &jsonrpc.Error{
				Code:    int32(protocol.ParseError),
				Message: "message is not of type Request",
			},
		}
		s.plane.Log().Warnf("%v", resp.Error.Message)
		s.send(resp)
		return
	}

	params, _ := rawDecode[protocol.InitializeParams](s, req.Params)

	if len(params.WorkspaceFolders) == 0 {
		resp := jsonrpc.Response{
			Id: req.Id,
			Error: &jsonrpc.Error{
				Code:    int32(protocol.InvalidRequest),
				Message: "no workspace folders",
			},
		}
		s.plane.Log().Warnf("%v", resp.Error.Message)
		s.send(&resp)
		return
	}

	for _, workspace := range params.WorkspaceFolders {
		s.plane.Log().Infof("workspace %s", workspace)
	}

	result := protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			TextDocumentSync: protocol.TextDocumentSyncOptions{
				OpenClose: true,
				Change:    protocol.Full,
			},
			CompletionProvider: &protocol.CompletionOptions{},
		},
		ServerInfo: &protocol.ServerInfo{
			Name:    info.LspServerName,
			Version: info.Version,
		},
	}

	resultRaw, err := rawEncode(s, &result)
	if err != nil {
		return
	}
	resp := jsonrpc.Response{
		Id:     req.Id,
		Result: resultRaw,
	}
	s.send(&resp)
}

func (s *ClientLsp) handleInitialized(msg jsonrpc.Message) {

}

func (s *ClientLsp) handleTextDocumentDidOpen(msg jsonrpc.Message) {

}

func (s *ClientLsp) handleTextDocumentDidChange(msg jsonrpc.Message) {

}

func (s *ClientLsp) handleTextDocumentDidClose(msg jsonrpc.Message) {

}

func (s *ClientLsp) handleTextDocumentCompletion(msg jsonrpc.Message) {

}

func (s *ClientLsp) handle(bytes []byte) {
	msg, err := jsonrpc.DecodeMessage(bytes)

	if err != nil {
		s.plane.Log().Errorf("%T %v: %v", s, s.Name(), err)
		return
	}

	var method string

	switch msg := msg.(type) {
	case *jsonrpc.Request:
		method = msg.Method
	case *jsonrpc.Notification:
		method = msg.Method
	default:
		s.plane.Log().Warnf("%T %v: unhandled type %T", s, s.Name(), msg)
		return
	}

	switch method {
	case "initialize":
		s.handleInitialize(msg)
	case "initialized":
		s.handleInitialized(msg)
	case "textDocument/didOpen":
		s.handleTextDocumentDidOpen(msg)
	case "textDocument/didChange":
		s.handleTextDocumentDidChange(msg)
	case "textDocument/didClose":
		s.handleTextDocumentDidClose(msg)
	case "textDocument/completion":
		s.handleTextDocumentCompletion(msg)
	default:
		s.plane.Log().Warnf("%T %v: unhandled method %v", s, s.Name(), method)
	}
}
