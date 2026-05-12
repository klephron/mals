package controller

import (
	"mals/pkg/config"
	"mals/third_party/lsp"
)

type LspStatus int32

const (
	LspAbsent     LspStatus = 0
	LspRegistered LspStatus = (1 << 0)
	LspCreated    LspStatus = (1 << 1)
	LspStarted    LspStatus = (1 << 2)
)

type LspData struct {
	Name         string
	Status       LspStatus
	Config       *config.Lsp
	Capabilities *lsp.ServerCapabilities
	Info         *lsp.ServerInfo
}

type LspController interface {
	ControllerRun(onReady func()) error
	ControllerShutdown() error

	Status(name string) LspStatus
	Register(name string, config *config.Lsp) error
	Unregister(name string) error
	Create(name string) error
	Delete(name string) error
	Start(name string) error
	Stop(name string) error

	GetCapabilities(name string) (*lsp.ServerCapabilities, error)
	GetInfo(name string) (*lsp.ServerInfo, error)
	Get(name string) (*LspData, error)
	GetAll() []*LspData

	Initialize(name string, params *lsp.InitializeParams) (*lsp.InitializeResult, error)
	Initialized(name string, params *lsp.InitializedParams) error
	TextDocumentDidOpen(name string, params *lsp.DidOpenTextDocumentParams) error
	TextDocumentDidChange(name string, params *lsp.DidChangeTextDocumentParams) error
	TextDocumentDidClose(name string, params *lsp.DidCloseTextDocumentParams) error
	TextDocumentCompletion(name string, params *lsp.CompletionParams) (*lsp.CompletionList, error)
	Shutdown(name string) error
}
