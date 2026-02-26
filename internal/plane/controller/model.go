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
	ControllerRun(onReady func()) error
	ControllerShutdown() error

	Status(name string) ModelStatus
	Register(name string, config *config.Model) error
	Unregister(name string) error
	Create(name string) error
	Delete(name string) error
	Start(name string) error
	Stop(name string) error

	Get(name string) (*ModelData, error)
	GetAll() []*ModelData

	TaskExecClient(name string, task *model.Task, clientName string) (string, error)
	TaskGet(name string, id uuid.UUID) (*model.Task, error)
	TaskGetClient(name string, id uuid.UUID, clientName string) (*model.Task, error)
	TaskGetAll(name string) ([]*model.Task, error)
	TaskGetAllClient(name string, clientName string) ([]*model.Task, error)
	TaskCancelClient(name string, id uuid.UUID, clientName string) (*model.Task, error)
	TaskCancelAllClient(name string, clientName string) ([]*model.Task, error)
}
