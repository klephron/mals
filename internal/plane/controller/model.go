package controller

import (
	"mals/internal/client"
	"mals/internal/model"
	"mals/pkg/config"

	"github.com/google/uuid"
)

type ModelStatus int32

const (
	ModelAbsent     ModelStatus = 0
	ModelRegistered ModelStatus = (1 << 0)
	ModelCreated    ModelStatus = (1 << 1)
	ModelStarted    ModelStatus = (1 << 2)
)

type ModelController interface {
	Shutdown() error
	Serve(onReady func()) error

	Status(name string) ModelStatus
	Register(name string, config *config.Model) error
	Unregister(name string) error
	Create(name string) error
	Delete(name string) error
	Start(name string) error
	Stop(name string) error

	TaskExecClient(modelName string, task *model.Task, client client.Client) (string, error)
	TaskGetClient(modelName string, id uuid.UUID, client client.Client) (*model.Task, error)
	TaskGetAllClient(modelName string, client client.Client) ([]*model.Task, error)
	TaskCancelClient(modelName string, id uuid.UUID, client client.Client) (*model.Task, error)
	TaskCancelAllClient(modelName string, client client.Client) ([]*model.Task, error)
}
