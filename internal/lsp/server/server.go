package server

import (
	"context"
	"mals/third_party/lsp"
)

type LspServer interface {
	Name() string
	Kind() string
	Run(ctx context.Context, onReady func()) error

	Capabilities() (*lsp.ServerCapabilities, error)
	Info() (*lsp.ServerInfo, error)

	Initialize(params *lsp.InitializeParams) (*lsp.InitializeResult, error)
	Initialized(params *lsp.InitializedParams) error
	TextDocumentDidOpen(params *lsp.DidOpenTextDocumentParams) error
	TextDocumentDidChange(params *lsp.DidChangeTextDocumentParams) error
	TextDocumentDidClose(params *lsp.DidCloseTextDocumentParams) error
	TextDocumentCompletion(params *lsp.CompletionParams) (*lsp.CompletionList, error)
	Shutdown() error
}
