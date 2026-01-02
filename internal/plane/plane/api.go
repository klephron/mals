package plane

import (
	"mals/internal/client"
	"mals/internal/listener"
	"mals/internal/lsp/protocol"
	"mals/internal/model"
	"mals/internal/plane/controller"
	"mals/internal/scope"
	"mals/internal/usage"
	"mals/pkg/config"

	"github.com/google/uuid"
)

func (s *Plane) ClientStatus(name string) controller.ClientStatus {
	return s.client.ClientStatus(name)
}

func (s *Plane) ClientOwn(name string, client client.Client, listener listener.Listener) error {
	return s.client.ClientOwn(name, client, listener)
}

func (s *Plane) ClientServe(name string) error {
	return s.client.ClientServe(name)
}

func (s *Plane) ClientShutdown(name string) error {
	return s.client.ClientShutdown(name)
}

func (s *Plane) ClientShutdownSilent(name string) error {
	return s.client.ClientShutdownSilent(name)
}

func (s *Plane) ClientGetListener(name string) (string, error) {
	return s.client.ClientGetListener(name)
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

func (s *Plane) ListenerClientAdd(name string, client string) error {
	return s.listener.ListenerClientAdd(name, client)
}

func (s *Plane) ListenerClientRemove(name string, client string) error {
	return s.listener.ListenerClientRemove(name, client)
}

func (s *Plane) ListenerGetConfig(name string) (config.Listener, error) {
	return s.listener.ListenerGetConfig(name)
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

func (s *Plane) LspStatus(name string) controller.LspStatus {
	return s.lsp.LspStatus(name)
}

func (s *Plane) LspRegister(name string, config *config.Lsp) error {
	return s.lsp.LspRegister(name, config)
}

func (s *Plane) LspUnregister(name string) error {
	return s.lsp.LspUnregister(name)
}

func (s *Plane) LspCreate(name string) error {
	return s.lsp.LspCreate(name)
}

func (s *Plane) LspDelete(name string) error {
	return s.lsp.LspDelete(name)
}

func (s *Plane) LspStart(name string) error {
	return s.lsp.LspStart(name)
}

func (s *Plane) LspStop(name string) error {
	return s.lsp.LspStop(name)
}

func (s *Plane) LspEventInitialize(name string, params *protocol.InitializeParams) (*protocol.InitializeResult, error) {
	return s.lsp.EventInitialize(name, params)
}

func (s *Plane) LspEventInitialized(name string, params *protocol.InitializedParams) error {
	return s.lsp.EventInitialized(name, params)
}

func (s *Plane) LspEventTextDocumentDidOpen(name string, params *protocol.DidOpenTextDocumentParams) error {
	return s.lsp.EventTextDocumentDidOpen(name, params)
}

func (s *Plane) LspEventTextDocumentDidChange(name string, params *protocol.DidChangeTextDocumentParams) error {
	return s.lsp.EventTextDocumentDidChange(name, params)
}

func (s *Plane) LspEventTextDocumentDidClose(name string, params *protocol.DidCloseTextDocumentParams) error {
	return s.lsp.EventTextDocumentDidClose(name, params)
}

func (s *Plane) LspEventTextDocumentCompletion(name string, params *protocol.CompletionParams) (*protocol.CompletionList, error) {
	return s.lsp.EventTextDocumentCompletion(name, params)
}

func (s *Plane) LspEventShutdown(name string) error {
	return s.lsp.EventShutdown(name)
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

func (s *Plane) ModelTaskExecClient(modelName string, task *model.Task, client client.Client) (string, error) {
	return s.model.TaskExecClient(modelName, task, client)
}

func (s *Plane) ModelTaskGetClient(modelName string, id uuid.UUID, client client.Client) (*model.Task, error) {
	return s.model.TaskGetClient(modelName, id, client)
}

func (s *Plane) ModelTaskGetAllClient(modelName string, client client.Client) ([]*model.Task, error) {
	return s.model.TaskGetAllClient(modelName, client)
}

func (s *Plane) ModelTaskCancelClient(modelName string, id uuid.UUID, client client.Client) (*model.Task, error) {
	return s.model.TaskCancelClient(modelName, id, client)
}

func (s *Plane) ModelTaskCancelAllClient(modelName string, client client.Client) ([]*model.Task, error) {
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

func (s *Plane) ScopeLspRegister(config config.Lsp) error {
	return s.scope.ScopeLspRegister(config)
}

func (s *Plane) ScopeLspAcquire(name string, scope *scope.Scope) (string, controller.ScopeToken, error) {
	return s.scope.ScopeLspAcquire(name, scope)
}

func (s *Plane) ScopeLspRelease(name string, token controller.ScopeToken) error {
	return s.scope.ScopeLspRelease(name, token)
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

func (s *Plane) UsageGetFiltered(condition usage.ConditionFilter, event usage.EventFilter) []*config.Usage {
	return s.usage.UsageGetFiltered(condition, event)
}

func (s *Plane) UsageGetFilteredClient(condition usage.ConditionFilter, event usage.EventFilter, client string) []*config.Usage {
	return s.usage.UsageGetFilteredClient(condition, event, client)
}
