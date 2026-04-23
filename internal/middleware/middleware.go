package middleware

import (
	"fmt"
	"mals/internal/lsp/protocol"
	"mals/internal/middleware/handler"
	"mals/internal/plane"
	"mals/pkg/config"
	"mals/pkg/info"
)

type Middleware struct {
	listenerName string
	clientName   string

	plane plane.Plane

	textDocumentSyncKind protocol.TextDocumentSyncKind

	handlers []*handler.Handler
}

func New(listenerName string, clientName string, plane plane.Plane) *Middleware {
	return &Middleware{
		listenerName:         listenerName,
		clientName:           clientName,
		plane:                plane,
		textDocumentSyncKind: protocol.Incremental,
		handlers:             make([]*handler.Handler, 0),
	}
}

func (s *Middleware) Name() string {
	return fmt.Sprintf("%v:%v", s.listenerName, s.clientName)
}

func (s *Middleware) Initialize(params *protocol.InitializeParams) (*protocol.InitializeResult, error) {
	if len(s.handlers) > 0 {
		err := fmt.Errorf("middleware is already initialized")
		s.plane.Errorf("%T: %v", s, err)
		return nil, err
	}

	listenerHandlers, err := s.plane.Listener().LspHandlerGetAll(s.listenerName)
	if err != nil {
		return nil, err
	}

	handlers := make([]*handler.Handler, 0)
	for _, listenerHandler := range listenerHandlers {
		handlerC, err := s.plane.Handler().Get(listenerHandler.Handler)
		if err != nil {
			return nil, err
		}

		handlerCLsp, ok := handlerC.Spec.(*config.HandlerSpecLsp)
		if !ok {
			err := fmt.Errorf("handler %v is not of type lsp", handlerC.Name)
			s.plane.Errorf("%T : %v", s, err)
			return nil, err
		}

		handler := handler.New(s.listenerName, s.clientName,
			handlerCLsp.Resources, &handlerCLsp.Endpoints, s.plane)

		if err := handler.Initialize(params); err != nil {
			return nil, err
		}

		handlers = append(handlers, handler)
	}
	s.handlers = handlers

	result := &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			TextDocumentSync: protocol.TextDocumentSyncOptions{
				OpenClose: true,
				Change:    s.textDocumentSyncKind,
			},
			CompletionProvider: &protocol.CompletionOptions{},
		},
		ServerInfo: &protocol.ServerInfo{
			Name:    info.MiddlewareServerName,
			Version: info.MiddlewareVersion,
		},
	}

	return result, nil
}

func (s *Middleware) Initialized(params *protocol.InitializedParams) error {
	// TODO
	return nil
}

func (s *Middleware) Shutdown() error {
	// TODO
	return nil
}

func (s *Middleware) TextDocumentCompletion(params *protocol.CompletionParams) (*protocol.CompletionList, error) {
	// TODO
	return nil, nil
}

func (s *Middleware) TextDocumentDidChange(params *protocol.DidChangeTextDocumentParams) error {
	// TODO
	return nil
}

func (s *Middleware) TextDocumentDidClose(params *protocol.DidCloseTextDocumentParams) error {
	// TODO
	return nil
}

func (s *Middleware) TextDocumentDidOpen(params *protocol.DidOpenTextDocumentParams) error {
	// TODO
	return nil
}
