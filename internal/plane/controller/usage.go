package controller

import (
	"mals/pkg/config"
)

type UsageController interface {
	Run(onReady func()) error
	Shutdown() error

	UsageRegister(config config.Usage) error
	UsageUnregister(name string) error
	UsageGetAll() []*config.Usage
	UsageGet(filetype *string, path *string, event *string) []*config.Usage
}
