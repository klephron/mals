package controller

import (
	"mals/internal/client"
	"mals/internal/model"
	"mals/pkg/config"

	"github.com/google/uuid"
)

type ModelController interface {
	Shutdown() error
	Serve(onReady func()) error

	Register(config config.Model) error
	Unregister(name string) error
	Create(name string) error
	Delete(name string) error
	Start(name string) error
	Stop(name string) error

	TaskExecClient(model string, task model.Task, client client.Client) model.Result
	TaskGetClient(model string, id uuid.UUID, client client.Client) (model.Task, error)
	TaskGetAllClient(model string, client client.Client) ([]model.Task, error)
	TaskGetAllClientName(model string, client string) ([]model.Task, error)
	TaskDeleteClient(model string, id uuid.UUID, client client.Client) (model.Task, error)
	TaskDeleteAllClient(model string, id uuid.UUID, client client.Client) ([]model.Task, error)
}
