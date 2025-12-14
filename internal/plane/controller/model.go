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

	TaskExecClient(modelName string, task *model.Task, client client.Client) model.Result
	TaskGetClient(modelName string, id uuid.UUID, client client.Client) (*model.Task, error)
	TaskGetAllClient(modelName string, client client.Client) ([]*model.Task, error)
	TaskCancelClient(modelName string, id uuid.UUID, client client.Client) (*model.Task, error)
	TaskCancelAllClient(modelName string, client client.Client) ([]*model.Task, error)
}
