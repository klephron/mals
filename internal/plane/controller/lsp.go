package controller

import (
	"mals/internal/lsp/protocol"
	"mals/pkg/config"
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
	Capabilities *protocol.ServerCapabilities
	Info         *protocol.ServerInfo
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

	GetCapabilities(name string) (*protocol.ServerCapabilities, error)
	GetInfo(name string) (*protocol.ServerInfo, error)
	Get(name string) (*LspData, error)
	GetAll() []*LspData

	Initialize(name string, params *protocol.InitializeParams) (*protocol.InitializeResult, error)
	Initialized(name string, params *protocol.InitializedParams) error
	TextDocumentDidOpen(name string, params *protocol.DidOpenTextDocumentParams) error
	TextDocumentDidChange(name string, params *protocol.DidChangeTextDocumentParams) error
	TextDocumentDidClose(name string, params *protocol.DidCloseTextDocumentParams) error
	TextDocumentCompletion(name string, params *protocol.CompletionParams) (*protocol.CompletionList, error)
	Shutdown(name string) error
}
