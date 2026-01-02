package controller

import (
	"mals/internal/scope"
	"mals/pkg/config"
)

type ScopeToken interface {
	Token() string
}

type ScopeController interface {
	Run(onReady func()) error
	Shutdown() error

	ScopeModelRegister(config config.Model) error
	ScopeModelAcquire(name string, scope *scope.Scope) (string, ScopeToken, error)
	ScopeModelRelease(name string, token ScopeToken) error

	ScopeLspRegister(config config.Lsp) error
	ScopeLspAcquire(name string, scope *scope.Scope) (string, ScopeToken, error)
	ScopeLspRelease(name string, token ScopeToken) error

	ScopeClose(scope *scope.Scope) []error
}
