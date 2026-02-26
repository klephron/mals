package controller

import (
	"mals/internal/scope"
	"mals/pkg/config"
	"time"
)

type ScopeToken interface {
	Token() string
	From() time.Time
}

type ResourceLsp struct {
	Fullname     string
	Config       *config.Lsp
	Dependencies map[string]ScopeToken
}

type ResourceModel struct {
	Fullname     string
	Config       *config.Model
	Dependencies map[string]ScopeToken
}

type Space struct {
	Space    scope.Space
	Children map[scope.Space]*Space
	Lsps     map[string]ResourceLsp
	Models   map[string]ResourceModel
}

type ScopeController interface {
	ControllerRun(onReady func()) error
	ControllerShutdown() error

	ModelRegister(config *config.Model) error
	ModelAcquire(name string, scope *scope.Scope) (string, ScopeToken, error)
	ModelRelease(name string, token ScopeToken) error

	LspRegister(config *config.Lsp) error
	LspAcquire(name string, scope *scope.Scope) (string, ScopeToken, error)
	LspRelease(name string, token ScopeToken) error

	Close(scope *scope.Scope) []error

	TreeRoot() *Space
}
