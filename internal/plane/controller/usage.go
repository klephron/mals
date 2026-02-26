package controller

import (
	"mals/internal/usage"
	"mals/pkg/config"
)

type UsageController interface {
	ControllerRun(onReady func()) error
	ControllerShutdown() error

	Register(config *config.Usage) error
	Unregister(name string) error
	Get(name string) (*config.Usage, error)
	GetFiltered(condition usage.ConditionFilter, event usage.EventFilter) []*config.Usage
	GetFilteredClient(condition usage.ConditionFilter, event usage.EventFilter, listener string, client string) []*config.Usage
	GetAll() []*config.Usage
}
