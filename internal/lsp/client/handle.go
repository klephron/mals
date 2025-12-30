package client

import (
	"fmt"
	"mals/internal/info"
	"mals/internal/jsonrpc"
	"mals/internal/lsp/protocol"
	"mals/pkg/config"
)

func errorParseUnexpectedType[T jsonrpc.Message](s *ClientLsp) {
	var dummy T

	resp := jsonrpc.Response{
		Error: &jsonrpc.Error{
			Code:    int32(protocol.ParseError),
			Message: fmt.Sprintf("message is not of type %T", dummy),
		},
	}

	s.plane.Log().Warnf("%v", resp.Error.Message)
	s.send(&resp)
}

func (s *ClientLsp) handleInitialize(msg jsonrpc.Message) {
	req, ok := msg.(*jsonrpc.Request)

	if !ok {
		errorParseUnexpectedType[*jsonrpc.Request](s)
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
		s.workspaceAdd(workspace.URI, workspace.Name)
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

func (s *ClientLsp) handleInitialized(_ jsonrpc.Message) {
	s.initialized = true
}

func (s *ClientLsp) handleTextDocumentDidOpen(msg jsonrpc.Message) {
	ntf, ok := msg.(*jsonrpc.Notification)
	if !ok {
		errorParseUnexpectedType[*jsonrpc.Notification](s)
		return
	}

	params, _ := rawDecode[protocol.DidOpenTextDocumentParams](s, ntf.Params)

	workspaces := s.workspaceFindAll(params.TextDocument.URI)
	for _, workspace := range workspaces {
		s.documentAdd(workspace, params.TextDocument.URI, &params.TextDocument.Text)
	}
}

func (s *ClientLsp) handleTextDocumentDidChange(msg jsonrpc.Message) {
	ntf, ok := msg.(*jsonrpc.Notification)
	if !ok {
		errorParseUnexpectedType[*jsonrpc.Notification](s)
		return
	}

	params, _ := rawDecode[protocol.DidChangeTextDocumentParams](s, ntf.Params)

	workspaces := s.workspaceFindAll(params.TextDocument.URI)
	for _, workspace := range workspaces {
		document := s.documentGet(workspace, params.TextDocument.URI)
		if document == nil {
			continue
		}

		var text *string
		for _, change := range params.ContentChanges {
			if change.Range != nil {
				s.plane.Log().Warnf("%s: unsupported ContentChange Range is unsupported", s.Name())
				continue
			}

			text = &change.Text
		}

		if text != nil {
			s.documentUpdateFull(workspace, document, text, params.TextDocument.Version)
		}
	}
}

func (s *ClientLsp) handleTextDocumentDidClose(msg jsonrpc.Message) {
	ntf, ok := msg.(*jsonrpc.Notification)
	if !ok {
		errorParseUnexpectedType[*jsonrpc.Notification](s)
		return
	}

	params, _ := rawDecode[protocol.DidCloseTextDocumentParams](s, ntf.Params)

	workspaces := s.workspaceFindAll(params.TextDocument.URI)
	for _, workspace := range workspaces {
		document := s.documentGet(workspace, params.TextDocument.URI)
		if document == nil {
			continue
		}

		s.documentDelete(workspace, document.uri)
	}
}

func (s *ClientLsp) handleTextDocumentCompletion(msg jsonrpc.Message) {
	req, ok := msg.(*jsonrpc.Request)
	if !ok {
		errorParseUnexpectedType[*jsonrpc.Request](s)
		return
	}

	params, _ := rawDecode[protocol.CompletionParams](s, req.Params)

	workspaces := s.workspaceFindAll(params.TextDocument.URI)

	for _, workspace := range workspaces {
		document := s.documentGet(workspace, params.TextDocument.URI)
		if document == nil {
			continue
		}

		s.plane.Log().Infof("%T %v: workspace %v document %v: completion", s, s.Name(), workspace.name, document.uri)

		usages := s.plane.Usage().Get(nil, &document.uri, &EventTextDocumentCompletion)

		for _, usage := range usages {
			go func(params protocol.CompletionParams, document *Document, workflow *config.Workflow) {
				list, err := s.completionWorkflow(&params, document, workflow)

				if err != nil {
					s.plane.Log().Warnf("%T %v: workspace %v %v", s, s.Name(), workspace.name, err)
					return
				}

				listRaw, err := rawEncode(s, &list)
				if err != nil {
					return
				}

				s.plane.Log().Infof("%T %v: workspace %v document %v: completion result %v", s, s.Name(), workspace.name, document.uri, string(listRaw))

				resp := jsonrpc.Response{
					Id:     req.Id,
					Result: listRaw,
				}
				s.send(&resp)

			}(params, document, usage.Workflow)
		}
	}
}

func (s *ClientLsp) handleShutdown(_ jsonrpc.Message) {
	s.workspaces.Range(func(key string, value *Workspace) bool {
		s.workspaceDelete(key)
		return true
	})
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
	case EventInitialize:
		s.handleInitialize(msg)
	case EventInitialized:
		s.handleInitialized(msg)
	case EventTextDocumentDidOpen:
		s.handleTextDocumentDidOpen(msg)
	case EventTextDocumentDidChange:
		s.handleTextDocumentDidChange(msg)
	case EventTextDocumentDidClose:
		s.handleTextDocumentDidClose(msg)
	case EventTextDocumentCompletion:
		s.handleTextDocumentCompletion(msg)
	case EventShutdown:
		s.handleShutdown(msg)
	default:
		s.plane.Log().Warnf("%T %v: unhandled method %v", s, s.Name(), method)
	}
}
