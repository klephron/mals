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
	ListenerRegister(name string, config *config.Listener) error
	ListenerUnregister(name string) error
	ListenerCreate(name string) error
	ListenerDelete(name string) error
	ListenerStart(name string) error
	ListenerStop(name string) error
	ListenerClientAdd(name string, client string) error
	ListenerClientRemove(name string, client string) error
	ListenerGetConfig(name string) (*config.Listener, error)
	ListenerGet(name string) (*controller.ListenerData, error)
	ListenerGetAll() []*controller.ListenerData

	LogStatus(name string) controller.LogStatus
	LogRegister(name string, config *config.Log) error
	LogUnregister(name string) error
	LogCreate(name string) error
	LogDelete(name string) error
	LogStart(name string) error
	LogStop(name string) error
	LogGet(name string) (*controller.LogData, error)
	LogGetAll() []*controller.LogData

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
	LspGetCapabilities(name string) (*protocol.ServerCapabilities, error)
	LspGetInfo(name string) (*protocol.ServerInfo, error)
	LspGet(name string) (*controller.LspData, error)
	LspGetAll() []*controller.LspData
	LspEventInitialize(name string, params *protocol.InitializeParams) (*protocol.InitializeResult, error)
	LspEventInitialized(name string, params *protocol.InitializedParams) error
	LspEventTextDocumentDidOpen(name string, params *protocol.DidOpenTextDocumentParams) error
	LspEventTextDocumentDidChange(name string, params *protocol.DidChangeTextDocumentParams) error
	LspEventTextDocumentDidClose(name string, params *protocol.DidCloseTextDocumentParams) error
	LspEventTextDocumentCompletion(name string, params *protocol.CompletionParams) (*protocol.CompletionList, error)
	LspEventShutdown(name string) error

	ModelStatus(name string) controller.ModelStatus
	ModelRegister(name string, config *config.Model) error
	ModelUnregister(name string) error
	ModelCreate(name string) error
	ModelDelete(name string) error
	ModelStart(name string) error
	ModelStop(name string) error
	ModelGet(name string) (*controller.ModelData, error)
	ModelGetAll() []*controller.ModelData
	ModelTaskExecClient(modelName string, task *model.Task, client client.Client) (string, error)
	ModelTaskGetClient(modelName string, id uuid.UUID, client client.Client) (*model.Task, error)
	ModelTaskGetAllClient(modelName string, client client.Client) ([]*model.Task, error)
	ModelTaskCancelClient(modelName string, id uuid.UUID, client client.Client) (*model.Task, error)
	ModelTaskCancelAllClient(modelName string, client client.Client) ([]*model.Task, error)

	ScopeModelRegister(config *config.Model) error
	ScopeModelAcquire(name string, scope *scope.Scope) (string, controller.ScopeToken, error)
	ScopeModelRelease(name string, token controller.ScopeToken) error
	ScopeLspRegister(config *config.Lsp) error
	ScopeLspAcquire(name string, scope *scope.Scope) (string, controller.ScopeToken, error)
	ScopeLspRelease(name string, token controller.ScopeToken) error
	ScopeClose(scope *scope.Scope) []error

	UsageRegister(config *config.Usage) error
	UsageUnregister(name string) error
	UsageGet(name string) (*config.Usage, error)
	UsageGetFiltered(condition usage.ConditionFilter, event usage.EventFilter) []*config.Usage
	UsageGetFilteredClient(condition usage.ConditionFilter, event usage.EventFilter, client string) []*config.Usage
	UsageGetAll() []*config.Usage
}
