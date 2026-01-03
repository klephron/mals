package server

import (
	"context"
	"mals/internal/lsp/protocol"
)

type LspServer interface {
	Name() string
	Kind() string
	Run(ctx context.Context, onReady func()) error

	Capabilities() (*protocol.ServerCapabilities, error)
	Info() (*protocol.ServerInfo, error)

	Initialize(params *protocol.InitializeParams) (*protocol.InitializeResult, error)
	Initialized(params *protocol.InitializedParams) error
	TextDocumentDidOpen(params *protocol.DidOpenTextDocumentParams) error
	TextDocumentDidChange(params *protocol.DidChangeTextDocumentParams) error
	TextDocumentDidClose(params *protocol.DidCloseTextDocumentParams) error
	TextDocumentCompletion(params *protocol.CompletionParams) (*protocol.CompletionList, error)
	Shutdown() error
}
