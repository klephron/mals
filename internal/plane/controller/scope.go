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

	ModelRegister(config config.Model) error
	ModelAcquire(name string, scope *scope.Scope) (string, ScopeToken, error)
	ModelRelease(name string, token ScopeToken) error

	ScopeClose(scope *scope.Scope) []error
}
