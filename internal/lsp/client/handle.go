package client

import (
	"mals/internal/jsonrpc"
	"mals/internal/lsp/protocol"
	"mals/pkg/config"
)

func (s *ClientLsp) handleInitialize(msg jsonrpc.Message) {
	req, ok := msg.(*jsonrpc.Request)

	if !ok {
		errorParseUnexpectedType[*jsonrpc.Request](s)
		return
	}

	params, err := rawDecode[protocol.InitializeParams](s, req.Params)
	if err != nil {
		s.plane.Warnf("%v", err)
		return
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

	capabilities, serverInfo, err := s.middleware.Initialize(&params)
	if err != nil {
		s.plane.Errorf("%v", err)
		return
	}

	result := protocol.InitializeResult{
		Capabilities: *capabilities,
		ServerInfo:   serverInfo,
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
	req, ok := msg.(*jsonrpc.Request)

	if !ok {
		errorParseUnexpectedType[*jsonrpc.Request](s)
		return
	}

	params, err := rawDecode[protocol.InitializedParams](s, req.Params)
	if err != nil {
		s.plane.Warnf("%v", err)
		return
	}

	err = s.middleware.Initialized(&params)
	if err != nil {
		s.plane.Errorf("%v", err)
		return
	}
}

func (s *ClientLsp) handleTextDocumentDidOpen(msg jsonrpc.Message) {
	ntf, ok := msg.(*jsonrpc.Notification)
	if !ok {
		errorParseUnexpectedType[*jsonrpc.Notification](s)
		return
	}

	params, err := rawDecode[protocol.DidOpenTextDocumentParams](s, ntf.Params)
	if err != nil {
		s.plane.Warnf("%v", err)
		return
	}

	err = s.middleware.TextDocumentDidOpen(&params)
	if err != nil {
		s.plane.Errorf("%v", err)
		return
	}
}

func (s *ClientLsp) handleTextDocumentDidChange(msg jsonrpc.Message) {
	ntf, ok := msg.(*jsonrpc.Notification)
	if !ok {
		errorParseUnexpectedType[*jsonrpc.Notification](s)
		return
	}

	params, err := rawDecode[protocol.DidChangeTextDocumentParams](s, ntf.Params)
	if err != nil {
		s.plane.Warnf("%v", err)
		return
	}

	err = s.middleware.TextDocumentDidChange(&params)
	if err != nil {
		s.plane.Errorf("%v", err)
		return
	}
}

func (s *ClientLsp) handleTextDocumentDidClose(msg jsonrpc.Message) {
	ntf, ok := msg.(*jsonrpc.Notification)
	if !ok {
		errorParseUnexpectedType[*jsonrpc.Notification](s)
		return
	}

	params, err := rawDecode[protocol.DidCloseTextDocumentParams](s, ntf.Params)
	if err != nil {
		s.plane.Warnf("%v", err)
		return
	}

	err = s.middleware.TextDocumentDidClose(&params)
	if err != nil {
		s.plane.Errorf("%v", err)
		return
	}
}

func (s *ClientLsp) handleTextDocumentCompletion(msg jsonrpc.Message) {
	req, ok := msg.(*jsonrpc.Request)
	if !ok {
		errorParseUnexpectedType[*jsonrpc.Request](s)
		return
	}

	params, err := rawDecode[protocol.CompletionParams](s, req.Params)
	if err != nil {
		s.plane.Warnf("%v", err)
		return
	}

	go func() {
		list, err := s.middleware.TextDocumentCompletion(&params)
		if err != nil {
			s.plane.Errorf("%v", err)
			return
		}

		listRaw, err := rawEncode(s, &list)
		if err != nil {
			return
		}

		resp := jsonrpc.Response{
			Id:     req.Id,
			Result: listRaw,
		}
		s.send(&resp)
	}()
}

func (s *ClientLsp) handleShutdown(_ jsonrpc.Message) {
	err := s.middleware.Shutdown()
	if err != nil {
		s.plane.Errorf("%v", err)
		return
	}
}

func (s *ClientLsp) handle(bytes []byte) {
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
