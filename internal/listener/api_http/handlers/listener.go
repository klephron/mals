package handlers

import (
	"context"
	"mals/internal/plane"
	"mals/internal/plane/controller"
	"mals/internal/util"
	"mals/pkg/config"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type ListenerDto struct {
	Name        string                    `json:"name"`
	Status      controller.ListenerStatus `json:"status"`
	StatusFlags []string                  `json:"status_flags"`
	Config      *config.WireListener      `json:"config"`
}

type ListenerGetInput struct {
	Name string `path:"name" doc:"Listener name"`
}

var (
	listenerStatusFlags = map[controller.ListenerStatus]string{
		controller.ListenerAbsent:     "ABSENT",
		controller.ListenerRegistered: "REGISTERED",
		controller.ListenerCreated:    "CREATED",
		controller.ListenerStarted:    "STARTED",
	}
)

func listenerToDto(listener *controller.ListenerData) ListenerDto {
	statusFlags := make([]string, 0, 4)
	for bit, name := range listenerStatusFlags {
		if listener.Status&bit != 0 {
			statusFlags = append(statusFlags, name)
		}
	}

	return ListenerDto{
		Name:        listener.Name,
		Status:      listener.Status,
		StatusFlags: statusFlags,
		Config:      util.Ptr(listener.Config.Wire()),
	}
}

func ListenerGetAllOperation() huma.Operation {
	return huma.Operation{
		OperationID: "listener-get-all",
		Method:      http.MethodGet,
		Path:        "/listeners",
		Summary:     "Get all listeners",
		Tags:        []string{"listener"},
	}
}

func ListenerGetAll(plane plane.Plane) func(ctx context.Context, input *struct{}) (*struct{ Body []ListenerDto }, error) {
	return func(ctx context.Context, input *struct{}) (*struct{ Body []ListenerDto }, error) {
		listeners := plane.ListenerGetAll()

		result := make([]ListenerDto, len(listeners))
		for i, listener := range listeners {
			result[i] = listenerToDto(listener)
		}

		return &struct{ Body []ListenerDto }{Body: result}, nil
	}
}

func ListenerGetOperation() huma.Operation {
	return huma.Operation{
		OperationID: "listener-get",
		Method:      http.MethodGet,
		Path:        "/listeners/{name}",
		Summary:     "Get listener by name",
		Tags:        []string{"listener"},
	}
}

func ListenerGet(plane plane.Plane) func(ctx context.Context, input *ListenerGetInput) (*struct{ Body ListenerDto }, error) {
	return func(ctx context.Context, input *ListenerGetInput) (*struct{ Body ListenerDto }, error) {
		listener, err := plane.ListenerGet(input.Name)

		if err != nil {
			return nil, huma.Error404NotFound("Listener not found", err)
		}

		result := listenerToDto(listener)

		return &struct{ Body ListenerDto }{Body: result}, nil
	}
}
