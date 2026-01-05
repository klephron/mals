package handlers

import (
	"context"
	"mals/internal/plane"
	"mals/internal/plane/controller"
	"mals/pkg/wire"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type LogDto struct {
	Name        string               `json:"name"`
	Status      controller.LogStatus `json:"status"`
	StatusFlags []string             `json:"status_flags"`
	Config      *wire.Log            `json:"config"`
}

type LogGetInput struct {
	Name string `path:"name" doc:"Log name"`
}

var (
	logStatusFlags = map[controller.LogStatus]string{
		controller.LogAbsent:     "ABSENT",
		controller.LogRegistered: "REGISTERED",
		controller.LogCreated:    "CREATED",
		controller.LogStarted:    "STARTED",
	}
)

func logToDto(log *controller.LogData) LogDto {
	statusFlags := make([]string, 0, 4)
	for bit, name := range logStatusFlags {
		if log.Status&bit != 0 {
			statusFlags = append(statusFlags, name)
		}
	}

	t := LogDto{
		Name:        log.Name,
		Status:      log.Status,
		StatusFlags: statusFlags,
		Config:      &wire.Log{},
	}
	t.Config.Wire(log.Config)

	return t
}

func LogGetAllOperation() huma.Operation {
	return huma.Operation{
		OperationID: "log-get-all",
		Method:      http.MethodGet,
		Path:        "/logs",
		Summary:     "Get all loggers",
		Tags:        []string{"log"},
	}
}

func LogGetAll(plane plane.Plane) func(ctx context.Context, input *struct{}) (*struct{ Body []LogDto }, error) {
	return func(ctx context.Context, input *struct{}) (*struct{ Body []LogDto }, error) {
		logs := plane.LogGetAll()

		result := make([]LogDto, len(logs))
		for i, Log := range logs {
			result[i] = logToDto(Log)
		}

		return &struct{ Body []LogDto }{Body: result}, nil
	}
}

func LogGetOperation() huma.Operation {
	return huma.Operation{
		OperationID: "log-get",
		Method:      http.MethodGet,
		Path:        "/logs/{name}",
		Summary:     "Get logger by name",
		Tags:        []string{"log"},
	}
}

func LogGet(plane plane.Plane) func(ctx context.Context, input *LogGetInput) (*struct{ Body LogDto }, error) {
	return func(ctx context.Context, input *LogGetInput) (*struct{ Body LogDto }, error) {
		log, err := plane.LogGet(input.Name)

		if err != nil {
			return nil, huma.Error404NotFound("Log not found", err)
		}

		result := logToDto(log)

		return &struct{ Body LogDto }{Body: result}, nil
	}
}
