package controller

import (
	"mals/internal/usage"
	"mals/pkg/config"
)

type UsageController interface {
	Run(onReady func()) error
	Shutdown() error

	UsageRegister(config *config.Usage) error
	UsageUnregister(name string) error
	UsageGet(name string) (*config.Usage, error)
	UsageGetFiltered(condition usage.ConditionFilter, event usage.EventFilter) []*config.Usage
	UsageGetFilteredClient(condition usage.ConditionFilter, event usage.EventFilter, listener string, client string) []*config.Usage
	UsageGetAll() []*config.Usage
}
