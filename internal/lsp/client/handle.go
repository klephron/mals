package client

import (
	"mals/internal/jsonrpc"
	"mals/internal/lsp/protocol"
	"mals/internal/util"
)

func (s *LspClient) handleInitialize(msg jsonrpc.Message) {
	req, ok := msg.(*jsonrpc.Request)

	if !ok {
		errorParseUnexpectedType[*jsonrpc.Request](s)
		return
	}

	params, err := util.JsonUnmarshal[protocol.InitializeParams](req.Params)
	if err != nil {
		s.plane.Warnf("%T %v: %v", s, s.Name(), err)
	}

	if len(params.WorkspaceFolders) == 0 {
		resp := jsonrpc.Response{
			Id: req.Id,
			Error: &jsonrpc.Error{
				Code:    int32(protocol.InvalidRequest),
				Message: "no workspace folders",
			},
		}
		s.plane.Warnf("%T %v: %v", s, s.Name(), resp.Error.Message)
		s.send(&resp)
		return
	}

	result, err := s.middleware.Initialize(&params, s.listenerName, s.clientName)
	if err != nil {
		s.plane.Errorf("%T %v: %v", s, s.Name(), err)
		return
	}

	resultRaw, err := util.JsonMarshal(result)
	if err != nil {
		s.plane.Errorf("%T %v: Initialize %v", s, s.Name(), err)
		return
	}
	resp := jsonrpc.Response{
		Id:     req.Id,
		Result: resultRaw,
	}
	s.send(&resp)
}

func (s *LspClient) handleInitialized(msg jsonrpc.Message) {
	ntf, ok := msg.(*jsonrpc.Notification)

	if !ok {
		errorParseUnexpectedType[*jsonrpc.Notification](s)
		return
	}

	params, err := util.JsonUnmarshal[protocol.InitializedParams](ntf.Params)
	if err != nil {
		s.plane.Warnf("%T %v: Initialized %v", s, s.Name(), err)
	}

	err = s.middleware.Initialized(&params)
	if err != nil {
		s.plane.Errorf("%T %v: Initialized %v", s, s.Name(), err)
		return
	}
}

func (s *LspClient) handleTextDocumentDidOpen(msg jsonrpc.Message) {
	ntf, ok := msg.(*jsonrpc.Notification)
	if !ok {
		errorParseUnexpectedType[*jsonrpc.Notification](s)
		return
	}

	params, err := util.JsonUnmarshal[protocol.DidOpenTextDocumentParams](ntf.Params)
	if err != nil {
		s.plane.Warnf("%T %v: TextDocumentDidOpen %v", s, s.Name(), err)
	}

	err = s.middleware.TextDocumentDidOpen(&params)
	if err != nil {
		s.plane.Errorf("%T %v: TextDocumentDidOpen %v", s, s.Name(), err)
		return
	}
}

func (s *LspClient) handleTextDocumentDidChange(msg jsonrpc.Message) {
	ntf, ok := msg.(*jsonrpc.Notification)
	if !ok {
		errorParseUnexpectedType[*jsonrpc.Notification](s)
		return
	}

	params, err := util.JsonUnmarshal[protocol.DidChangeTextDocumentParams](ntf.Params)
	if err != nil {
		s.plane.Warnf("%T %v: TextDocumentDidChange %v", s, s.Name(), err)
	}

	err = s.middleware.TextDocumentDidChange(&params)
	if err != nil {
		s.plane.Errorf("%T %v: TextDocumentDidChange %v", s, s.Name(), err)
		return
	}
}

func (s *LspClient) handleTextDocumentDidClose(msg jsonrpc.Message) {
	ntf, ok := msg.(*jsonrpc.Notification)
	if !ok {
		errorParseUnexpectedType[*jsonrpc.Notification](s)
		return
	}

	params, err := util.JsonUnmarshal[protocol.DidCloseTextDocumentParams](ntf.Params)
	if err != nil {
		s.plane.Warnf("%T %v: TextDocumentDidClose %v", s, s.Name(), err)
	}

	err = s.middleware.TextDocumentDidClose(&params)
	if err != nil {
		s.plane.Errorf("%T %v: TextDocumentDidClose %v", s, s.Name(), err)
		return
	}
}

func (s *LspClient) handleTextDocumentCompletion(msg jsonrpc.Message) {
	req, ok := msg.(*jsonrpc.Request)
	if !ok {
		errorParseUnexpectedType[*jsonrpc.Request](s)
		return
	}

	params, err := util.JsonUnmarshal[protocol.CompletionParams](req.Params)
	if err != nil {
		s.plane.Warnf("%T %v: TextDocumentCompletion %v", s, s.Name(), err)
	}

	go func() {
		list, err := s.middleware.TextDocumentCompletion(&params)
		if err != nil {
			s.plane.Errorf("%T %v: TextDocumentCompletion %v", s, s.Name(), err)
			return
		}

		listRaw, err := util.JsonMarshal(&list)
		if err != nil {
			s.plane.Errorf("%T %v: TextDocumentCompletion %v", s, s.Name(), err)
			return
		}

		resp := jsonrpc.Response{
			Id:     req.Id,
			Result: listRaw,
		}
		s.send(&resp)
	}()
}

func (s *LspClient) handleShutdown(_ jsonrpc.Message) {
	err := s.middleware.Shutdown()
	if err != nil {
		s.plane.Errorf("%T %v: TextDocumentShutdown %v", s, s.Name(), err)
		return
	}
}

func (s *LspClient) handle(bytes []byte) {
	msg, err := jsonrpc.DecodeMessage(bytes)

	if err != nil {
		s.plane.Errorf("%T %v: %v", s, s.Name(), err)
		return
	}

	var method string

	switch msg := msg.(type) {
	case *jsonrpc.Request:
	case *jsonrpc.Notification:
		method = msg.Method
	default:
		s.plane.Warnf("%T %v: unhandled type %T", s, s.Name(), msg)
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
	case "textDocument/shutdown":
		s.handleShutdown(msg)
	default:
		s.plane.Warnf("%T %v: unhandled method %v", s, s.Name(), method)
	}
}
