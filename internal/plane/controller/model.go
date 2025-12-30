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
	Run(onReady func()) error
	Shutdown() error

	ModelStatus(name string) ModelStatus
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
}
