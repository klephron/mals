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

	EventInitialize(name string, params *protocol.InitializeParams) (*protocol.InitializeResult, error)
	EventInitialized(name string, params *protocol.InitializedParams) error
	EventTextDocumentDidOpen(name string, params *protocol.DidOpenTextDocumentParams) error
	EventTextDocumentDidChange(name string, params *protocol.DidChangeTextDocumentParams) error
	EventTextDocumentDidClose(name string, params *protocol.DidCloseTextDocumentParams) error
	EventTextDocumentCompletion(name string, params *protocol.CompletionParams) (*protocol.CompletionList, error)
	EventShutdown(name string) error
}
