package plane

import (
	"mals/internal/client"
	"mals/internal/listener"
	"mals/internal/model"
	"mals/internal/plane/controller"
	"mals/internal/scope"
	"mals/internal/usage"
	"mals/pkg/config"

	"github.com/google/uuid"
)

type Plane interface {
	Run(onReady func())
	Shutdown() error

	ClientStatus(name string) controller.ClientStatus
	ClientOwn(name string, client client.Client, listener listener.Listener) error
	ClientServe(name string) error
	ClientShutdown(name string) error
	ClientShutdownSilent(name string) error
	ClientGetListener(name string) (string, error)

	ListenerStatus(name string) controller.ListenerStatus
	ListenerRegister(name string, config config.Listener) error
	ListenerUnregister(name string) error
	ListenerCreate(name string) error
	ListenerDelete(name string) error
	ListenerStart(name string) error
	ListenerStop(name string) error
	ListenerClientAdd(name string, client string) error
	ListenerClientRemove(name string, client string) error
	ListenerGetConfig(name string) (config.Listener, error)

	LogStatus(name string) controller.LogStatus
	LogRegister(name string, config config.Log) error
	LogUnregister(name string) error
	LogCreate(name string) error
	LogDelete(name string) error
	LogStart(name string) error
	LogStop(name string) error

	Debugf(format string, a ...any) error
	Infof(format string, a ...any) error
	Warnf(format string, a ...any) error
	Errorf(format string, a ...any) error

	LspStatus(name string) controller.LspStatus
	LspRegister(name string, config *config.Lsp) error
	LspUnregister(name string) error
	LspCreate(name string) error
	LspDelete(name string) error
	LspStart(name string) error
	LspStop(name string) error

	ModelStatus(name string) controller.ModelStatus
	ModelRegister(name string, config *config.Model) error
	ModelUnregister(name string) error
	ModelCreate(name string) error
	ModelDelete(name string) error
	ModelStart(name string) error
	ModelStop(name string) error

	TaskExecClient(modelName string, task *model.Task, client client.Client) (string, error)
	TaskGetClient(modelName string, id uuid.UUID, client client.Client) (*model.Task, error)
	TaskGetAllClient(modelName string, client client.Client) ([]*model.Task, error)
	TaskCancelClient(modelName string, id uuid.UUID, client client.Client) (*model.Task, error)
	TaskCancelAllClient(modelName string, client client.Client) ([]*model.Task, error)

	ScopeModelRegister(config config.Model) error
	ScopeModelAcquire(name string, scope *scope.Scope) (string, controller.ScopeToken, error)
	ScopeModelRelease(name string, token controller.ScopeToken) error
	ScopeClose(scope *scope.Scope) []error

	UsageRegister(config config.Usage) error
	UsageUnregister(name string) error
	UsageGetAll() []*config.Usage
	UsageGetFiltered(condition usage.ConditionFilter, event usage.EventFilter) []*config.Usage
	UsageGetFilteredClient(condition usage.ConditionFilter, event usage.EventFilter, client string) []*config.Usage
}
