package handlers

import (
	"context"
	"mals/internal/plane"
	"mals/pkg/config"
	"mals/pkg/wire"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type HandlerDto wire.Handler

type HandlerGetInput struct {
	Name string `path:"name" doc:"handler name"`
}

func handlerToDto(handler *config.Handler) HandlerDto {
	t := wire.Handler{}
	t.Wire(handler)
	return HandlerDto(t)
}

func HandlerGetAllOperation() huma.Operation {
	return huma.Operation{
		OperationID: "handler-get-all",
		Method:      http.MethodGet,
		Path:        "/handlers",
		Summary:     "Get all handlers",
		Tags:        []string{"handler"},
	}
}

func HandlerGetAll(plane plane.Plane) func(ctx context.Context, input *struct{}) (*struct{ Body []HandlerDto }, error) {
	return func(ctx context.Context, input *struct{}) (*struct{ Body []HandlerDto }, error) {
		handlers := plane.Handler().GetAll()

		result := make([]HandlerDto, len(handlers))
		for i, handler := range handlers {
			result[i] = handlerToDto(handler)
		}

		return &struct{ Body []HandlerDto }{Body: result}, nil
	}
}

func HandlerGetOperation() huma.Operation {
	return huma.Operation{
		OperationID: "handler-get",
		Method:      http.MethodGet,
		Path:        "/handler/{name}",
		Summary:     "Get handler by name",
		Tags:        []string{"handler"},
	}
}

func HandlerGet(plane plane.Plane) func(ctx context.Context, input *HandlerGetInput) (*struct{ Body HandlerDto }, error) {
	return func(ctx context.Context, input *HandlerGetInput) (*struct{ Body HandlerDto }, error) {
		handler, err := plane.Handler().Get(input.Name)

		if err != nil {
			return nil, huma.Error404NotFound("handler not found", err)
		}

		result := handlerToDto(handler)

		return &struct{ Body HandlerDto }{Body: result}, nil
	}
}
