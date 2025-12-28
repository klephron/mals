package controller

import (
	"mals/internal/scope"
	"mals/pkg/config"
)

type ScopeToken interface {
	Token() string
}

type ScopeController interface {
	Shutdown() error
	Serve(onReady func()) error

	// currently can't set root of object to create scopes

	ModelRegister(config config.Model) error
	ModelAcquire(name string, scope *scope.Scope) (string, ScopeToken, error)
	ModelRelease(name string, token ScopeToken) error
}
