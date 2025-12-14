package lsp

import (
	"mals/internal/jsonrpc"
	// "mals/internal/lsp/protocol"
)

func (s *ClientLsp) handleInitialize(msg jsonrpc.Message) {

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
