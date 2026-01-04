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
	Run(onReady func()) error
	Shutdown() error

	ScopeModelRegister(config *config.Model) error
	ScopeModelAcquire(name string, scope *scope.Scope) (string, ScopeToken, error)
	ScopeModelRelease(name string, token ScopeToken) error

	ScopeLspRegister(config *config.Lsp) error
	ScopeLspAcquire(name string, scope *scope.Scope) (string, ScopeToken, error)
	ScopeLspRelease(name string, token ScopeToken) error

	ScopeClose(scope *scope.Scope) []error

	ScopeTreeRoot() *Space
}
