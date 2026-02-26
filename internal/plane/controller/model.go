package controller

import (
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

type ModelData struct {
	Name   string
	Status ModelStatus
	Config *config.Model
	Tasks  []*model.Task
}

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

	ModelGet(name string) (*ModelData, error)
	ModelGetAll() []*ModelData

	TaskExecClient(modelName string, task *model.Task, clientName string) (string, error)
	TaskGet(modelName string, id uuid.UUID) (*model.Task, error)
	TaskGetClient(modelName string, id uuid.UUID, clientName string) (*model.Task, error)
	TaskGetAll(modelName string) ([]*model.Task, error)
	TaskGetAllClient(modelName string, clientName string) ([]*model.Task, error)
	TaskCancelClient(modelName string, id uuid.UUID, clientName string) (*model.Task, error)
	TaskCancelAllClient(modelName string, clientName string) ([]*model.Task, error)
}
