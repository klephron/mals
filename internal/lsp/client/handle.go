package client

import (
	"mals/internal/jsonrpc"
	"mals/internal/lsp/protocol"
	"mals/internal/util"
	"mals/pkg/config"
)

func (s *LspClient) handleInitialize(msg jsonrpc.Message) {
	req, ok := msg.(*jsonrpc.Request)

	if !ok {
		errorParseUnexpectedType[*jsonrpc.Request](s)
		return
	}

	params, err := util.JsonUnmarshal[protocol.InitializeParams](req.Params)
	if err != nil {
		s.plane.Warnf("%v", err)
	}

	if len(params.WorkspaceFolders) == 0 {
		resp := jsonrpc.Response{
			Id: req.Id,
			Error: &jsonrpc.Error{
				Code:    int32(protocol.InvalidRequest),
				Message: "no workspace folders",
			},
		}
		s.plane.Warnf("%v", resp.Error.Message)
		s.send(&resp)
		return
	}

	result, err := s.middleware.Initialize(&params, s.listenerName, s.clientName)
	if err != nil {
		s.plane.Errorf("%v", err)
		return
	}

	resultRaw, err := util.JsonMarshal(result)
	if err != nil {
		s.plane.Errorf("Initialize: %v", err)
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
		s.plane.Warnf("Initialized %v", err)
	}

	err = s.middleware.Initialized(&params)
	if err != nil {
		s.plane.Errorf("Initialized %v", err)
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
		s.plane.Warnf("TextDocumentDidOpen %v", err)
	}

	err = s.middleware.TextDocumentDidOpen(&params)
	if err != nil {
		s.plane.Errorf("TextDocumentDidOpen %v", err)
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
		s.plane.Warnf("TextDocumentDidChange %v", err)
	}

	err = s.middleware.TextDocumentDidChange(&params)
	if err != nil {
		s.plane.Errorf("TextDocumentDidChange %v", err)
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
		s.plane.Warnf("TextDocumentDidClose %v", err)
	}

	err = s.middleware.TextDocumentDidClose(&params)
	if err != nil {
		s.plane.Errorf("TextDocumentDidClose %v", err)
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
		s.plane.Warnf("TextDocumentCompletion %v", err)
	}

	go func() {
		list, err := s.middleware.TextDocumentCompletion(&params)
		if err != nil {
			s.plane.Errorf("TextDocumentCompletion %v", err)
			return
		}

		listRaw, err := util.JsonMarshal(&list)
		if err != nil {
			s.plane.Errorf("TextDocumentCompletion %v", err)
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
		s.plane.Errorf("TextDocumentShutdown %v", err)
		return
	}
}

func (s *LspClient) handle(bytes []byte) {
	msg, err := jsonrpc.DecodeMessage(bytes)

	if err != nil {
		s.plane.Errorf("%T %v: %v", s, s.Name(), err)
		return
	}

	var method config.Event

	switch msg := msg.(type) {
	case *jsonrpc.Request:
		method = config.Event(msg.Method)
	case *jsonrpc.Notification:
		method = config.Event(msg.Method)
	default:
		s.plane.Warnf("%T %v: unhandled type %T", s, s.Name(), msg)
		return
	}

	switch method {
	case config.EventInitialize:
		s.handleInitialize(msg)
	case config.EventInitialized:
		s.handleInitialized(msg)
	case config.EventTextDocumentDidOpen:
		s.handleTextDocumentDidOpen(msg)
	case config.EventTextDocumentDidChange:
		s.handleTextDocumentDidChange(msg)
	case config.EventTextDocumentDidClose:
		s.handleTextDocumentDidClose(msg)
	case config.EventTextDocumentCompletion:
		s.handleTextDocumentCompletion(msg)
	case config.EventShutdown:
		s.handleShutdown(msg)
	default:
		s.plane.Warnf("%T %v: unhandled method %v", s, s.Name(), method)
	}
}
