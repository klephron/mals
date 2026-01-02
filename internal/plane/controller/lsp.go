package controller

import "mals/pkg/config"

type LspStatus int32

const (
	LspAbsent     LspStatus = 0
	LspRegistered LspStatus = (1 << 0)
	LspCreated    LspStatus = (1 << 1)
	LspStarted    LspStatus = (1 << 2)
)

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
}
