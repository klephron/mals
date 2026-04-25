package handler

import (
	"context"
	"mals/internal/model"
	"mals/internal/plane"
	"mals/internal/plane/controller"
	"mals/pkg/wire"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type ModelDto struct {
	Name        string                 `json:"name"`
	Status      controller.ModelStatus `json:"status"`
	StatusFlags []string               `json:"status_flags"`
	Config      *wire.Model            `json:"config"`
	Tasks       []*model.Task          `json:"tasks"`
}

type ModelGetInput struct {
	Name string `path:"name" doc:"model name"`
}

var (
	modelStatusFlags = map[controller.ModelStatus]string{
		controller.ModelAbsent:     "ABSENT",
		controller.ModelRegistered: "REGISTERED",
		controller.ModelCreated:    "CREATED",
		controller.ModelStarted:    "STARTED",
	}
)

func modelToDto(model *controller.ModelData) ModelDto {
	statusFlags := make([]string, 0, 4)
	for bit, name := range modelStatusFlags {
		if model.Status&bit != 0 {
			statusFlags = append(statusFlags, name)
		}
	}

	t := ModelDto{
		Name:        model.Name,
		Status:      model.Status,
		StatusFlags: statusFlags,
		Config:      &wire.Model{},
		Tasks:       model.Tasks,
	}
	t.Config.Wire(model.Config)

	return t
}

func ModelGetAllOperation() huma.Operation {
	return huma.Operation{
		OperationID: "model-get-all",
		Method:      http.MethodGet,
		Path:        "/models",
		Summary:     "Get all models",
		Tags:        []string{"model"},
	}
}

func ModelGetAll(plane plane.Plane) func(ctx context.Context, input *struct{}) (*struct{ Body []ModelDto }, error) {
	return func(ctx context.Context, input *struct{}) (*struct{ Body []ModelDto }, error) {
		models := plane.Model().GetAll()

		result := make([]ModelDto, len(models))
		for i, model := range models {
			result[i] = modelToDto(model)
		}

		return &struct{ Body []ModelDto }{Body: result}, nil
	}
}

func ModelGetOperation() huma.Operation {
	return huma.Operation{
		OperationID: "model-get",
		Method:      http.MethodGet,
		Path:        "/models/{name}",
		Summary:     "Get model by name",
		Tags:        []string{"model"},
	}
}

func ModelGet(plane plane.Plane) func(ctx context.Context, input *ModelGetInput) (*struct{ Body ModelDto }, error) {
	return func(ctx context.Context, input *ModelGetInput) (*struct{ Body ModelDto }, error) {
		model, err := plane.Model().Get(input.Name)

		if err != nil {
			return nil, huma.Error404NotFound("model not found", err)
		}

		result := modelToDto(model)

		return &struct{ Body ModelDto }{Body: result}, nil
	}
}
