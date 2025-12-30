package plane

import (
	"mals/internal/client"
	"mals/internal/listener"
	"mals/internal/model"
	"mals/internal/plane/controller"
	"mals/internal/scope"
	"mals/pkg/config"

	"github.com/google/uuid"
)

type Plane interface {
	Run(onReady func())
	Shutdown() error

	ClientOwn(client client.Client, listener listener.Listener) error
	ClientServe(client client.Client) error
	ClientShutdown(client client.Client) error
	ClientShutdownSilent(client client.Client) error

	ListenerStatus(name string) controller.ListenerStatus
	ListenerRegister(name string, config config.Listener) error
	ListenerUnregister(name string) error
	ListenerCreate(name string) error
	ListenerDelete(name string) error
	ListenerStart(name string) error
	ListenerStop(name string) error
	ListenerClientAdd(name string, client client.Client) error
	ListenerClientRemove(name string, client client.Client) error

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
	UsageGet(filetype *string, path *string, event *string) []*config.Usage
}
