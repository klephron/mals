package middleware

import (
	"fmt"
	"mals/internal/middleware/handler"
	"mals/internal/plane"
	"mals/internal/scope"
	"mals/pkg/config"
	"mals/pkg/info"
	"mals/third_party/lsp"
)

type Middleware struct {
	listenerName string
	clientName   string

	plane plane.Plane

	textDocumentSyncKind lsp.TextDocumentSyncKind

	handlers    []*handler.Handler
	initialized bool
}

func New(listenerName string, clientName string, plane plane.Plane) *Middleware {
	return &Middleware{
		listenerName:         listenerName,
		clientName:           clientName,
		plane:                plane,
		textDocumentSyncKind: lsp.Incremental,
		handlers:             make([]*handler.Handler, 0),
		initialized:          false,
	}
}

func (s *Middleware) getScope() *scope.Scope {
	return scope.NewScopeClient(s.listenerName, s.clientName)
}

func (s *Middleware) Name() string {
	return fmt.Sprintf("%v:%v", s.listenerName, s.clientName)
}

func (s *Middleware) Initialize(params *lsp.InitializeParams) (*lsp.InitializeResult, error) {
	if s.initialized {
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

		handler := handler.New(s.listenerName, s.clientName, handlerC.Name,
			handlerCLsp.Resources, &handlerCLsp.Endpoints, s.plane)

		if err := handler.Initialize(params); err != nil {
			return nil, err
		}

		handlers = append(handlers, handler)
	}
	s.handlers = handlers

	result := &lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{
			TextDocumentSync: lsp.TextDocumentSyncOptions{
				OpenClose: true,
				Change:    s.textDocumentSyncKind,
			},
			CompletionProvider: &lsp.CompletionOptions{},
		},
		ServerInfo: &lsp.ServerInfo{
			Name:    info.MiddlewareServerName,
			Version: info.MiddlewareVersion,
		},
	}

	return result, nil
}

func (s *Middleware) Initialized(params *lsp.InitializedParams) error {
	for _, handler := range s.handlers {
		if err := handler.Initialized(params); err != nil {
			return err
		}
	}

	s.initialized = true

	return nil
}

func (s *Middleware) Shutdown() error {
	var error error

	for _, handler := range s.handlers {
		if err := handler.Shutdown(); err != nil {
			error = err
		}
	}

	s.handlers = s.handlers[:0]
	s.initialized = false

	s.plane.Scope().Close(s.getScope())

	s.plane.Infof("%T %v: Shutdown done", s, s.Name())

	return error
}

func (s *Middleware) TextDocumentCompletion(params *lsp.CompletionParams) (*lsp.CompletionList, error) {
	completionList := lsp.CompletionList{}
	var error error

	for _, handler := range s.handlers {
		list, err := handler.TextDocumentCompletion(params)
		if err != nil {
			error = err
			continue
		}
		if list.Items != nil {
			completionList.Items = append(completionList.Items, list.Items...)
		}
	}

	return &completionList, error
}

func (s *Middleware) TextDocumentDidChange(params *lsp.DidChangeTextDocumentParams) error {
	var error error

	for _, handler := range s.handlers {
		if err := handler.TextDocumentDidChange(params); err != nil {
			error = err
		}
	}

	return error
}

func (s *Middleware) TextDocumentDidClose(params *lsp.DidCloseTextDocumentParams) error {
	var error error

	for _, handler := range s.handlers {
		if err := handler.TextDocumentDidClose(params); err != nil {
			error = err
		}
	}

	return error
}

func (s *Middleware) TextDocumentDidOpen(params *lsp.DidOpenTextDocumentParams) error {
	var error error

	for _, handler := range s.handlers {
		if err := handler.TextDocumentDidOpen(params); err != nil {
			error = err
		}
	}

	return error
}
