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

func (s *Plane) ClientOwn(client client.Client, listener listener.Listener) error {
	return s.client.ClientOwn(client, listener)
}

func (s *Plane) ClientServe(client client.Client) error {
	return s.client.ClientServe(client)
}

func (s *Plane) ClientShutdown(client client.Client) error {
	return s.client.ClientShutdown(client)
}

func (s *Plane) ClientShutdownSilent(client client.Client) error {
	return s.client.ClientShutdownSilent(client)
}

func (s *Plane) ListenerStatus(name string) controller.ListenerStatus {
	return s.listener.ListenerStatus(name)
}

func (s *Plane) ListenerRegister(name string, config config.Listener) error {
	return s.listener.ListenerRegister(name, config)
}

func (s *Plane) ListenerUnregister(name string) error {
	return s.listener.ListenerUnregister(name)
}

func (s *Plane) ListenerCreate(name string) error {
	return s.listener.ListenerCreate(name)
}

func (s *Plane) ListenerDelete(name string) error {
	return s.listener.ListenerDelete(name)
}

func (s *Plane) ListenerStart(name string) error {
	return s.listener.ListenerStart(name)
}

func (s *Plane) ListenerStop(name string) error {
	return s.listener.ListenerStop(name)
}

func (s *Plane) ListenerClientAdd(name string, client client.Client) error {
	return s.listener.ListenerClientAdd(name, client)
}

func (s *Plane) ListenerClientRemove(name string, client client.Client) error {
	return s.listener.ListenerClientRemove(name, client)
}

func (s *Plane) LogStatus(name string) controller.LogStatus {
	return s.log.LogStatus(name)
}

func (s *Plane) LogRegister(name string, config config.Log) error {
	return s.log.LogRegister(name, config)
}

func (s *Plane) LogUnregister(name string) error {
	return s.log.LogUnregister(name)
}

func (s *Plane) LogCreate(name string) error {
	return s.log.LogCreate(name)
}

func (s *Plane) LogDelete(name string) error {
	return s.log.LogDelete(name)
}

func (s *Plane) LogStart(name string) error {
	return s.log.LogStart(name)
}

func (s *Plane) LogStop(name string) error {
	return s.log.LogStop(name)
}

func (s *Plane) Debugf(format string, a ...any) error {
	return s.log.Debugf(format, a...)
}

func (s *Plane) Infof(format string, a ...any) error {
	return s.log.Infof(format, a...)
}

func (s *Plane) Warnf(format string, a ...any) error {
	return s.log.Warnf(format, a...)
}

func (s *Plane) Errorf(format string, a ...any) error {
	return s.log.Errorf(format, a...)
}

func (s *Plane) ModelStatus(name string) controller.ModelStatus {
	return s.model.ModelStatus(name)
}

func (s *Plane) ModelRegister(name string, config *config.Model) error {
	return s.model.ModelRegister(name, config)
}

func (s *Plane) ModelUnregister(name string) error {
	return s.model.ModelUnregister(name)
}

func (s *Plane) ModelCreate(name string) error {
	return s.model.ModelCreate(name)
}

func (s *Plane) ModelDelete(name string) error {
	return s.model.ModelDelete(name)
}

func (s *Plane) ModelStart(name string) error {
	return s.model.ModelStart(name)
}

func (s *Plane) ModelStop(name string) error {
	return s.model.ModelStop(name)
}

func (s *Plane) TaskExecClient(modelName string, task *model.Task, client client.Client) (string, error) {
	return s.model.TaskExecClient(modelName, task, client)
}

func (s *Plane) TaskGetClient(modelName string, id uuid.UUID, client client.Client) (*model.Task, error) {
	return s.model.TaskGetClient(modelName, id, client)
}

func (s *Plane) TaskGetAllClient(modelName string, client client.Client) ([]*model.Task, error) {
	return s.model.TaskGetAllClient(modelName, client)
}

func (s *Plane) TaskCancelClient(modelName string, id uuid.UUID, client client.Client) (*model.Task, error) {
	return s.model.TaskCancelClient(modelName, id, client)
}

func (s *Plane) TaskCancelAllClient(modelName string, client client.Client) ([]*model.Task, error) {
	return s.model.TaskCancelAllClient(modelName, client)
}

func (s *Plane) ScopeModelRegister(config config.Model) error {
	return s.scope.ScopeModelRegister(config)
}

func (s *Plane) ScopeModelAcquire(name string, scope *scope.Scope) (string, controller.ScopeToken, error) {
	return s.scope.ScopeModelAcquire(name, scope)
}

func (s *Plane) ScopeModelRelease(name string, token controller.ScopeToken) error {
	return s.scope.ScopeModelRelease(name, token)
}

func (s *Plane) ScopeClose(scope *scope.Scope) []error {
	return s.scope.ScopeClose(scope)
}

func (s *Plane) UsageRegister(config config.Usage) error {
	return s.usage.UsageRegister(config)
}

func (s *Plane) UsageUnregister(name string) error {
	return s.usage.UsageUnregister(name)
}

func (s *Plane) UsageGetAll() []*config.Usage {
	return s.usage.UsageGetAll()
}

func (s *Plane) UsageGet(filetype *string, path *string, event *string) []*config.Usage {
	return s.usage.UsageGet(filetype, path, event)
}
