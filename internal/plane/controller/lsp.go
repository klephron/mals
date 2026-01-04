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
	Run(onReady func()) error
	Shutdown() error

	LspStatus(name string) LspStatus
	LspRegister(name string, config *config.Lsp) error
	LspUnregister(name string) error
	LspCreate(name string) error
	LspDelete(name string) error
	LspStart(name string) error
	LspStop(name string) error

	LspGetCapabilities(name string) (*protocol.ServerCapabilities, error)
	LspGetInfo(name string) (*protocol.ServerInfo, error)
	LspGet(name string) (*LspData, error)
	LspGetAll() []*LspData

	EventInitialize(name string, params *protocol.InitializeParams) (*protocol.InitializeResult, error)
	EventInitialized(name string, params *protocol.InitializedParams) error
	EventTextDocumentDidOpen(name string, params *protocol.DidOpenTextDocumentParams) error
	EventTextDocumentDidChange(name string, params *protocol.DidChangeTextDocumentParams) error
	EventTextDocumentDidClose(name string, params *protocol.DidCloseTextDocumentParams) error
	EventTextDocumentCompletion(name string, params *protocol.CompletionParams) (*protocol.CompletionList, error)
	EventShutdown(name string) error
}
