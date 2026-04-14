package handlers

import (
	"context"
	"mals/internal/plane"
	"mals/pkg/config"
	"mals/pkg/wire"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type UsageDto wire.Handler

type UsageGetInput struct {
	Name string `path:"name" doc:"usage name"`
}

func usageToDto(usage *config.Usage) UsageDto {
	t := wire.Handler{}
	t.Wire(usage)
	return UsageDto(t)
}

func UsageGetAllOperation() huma.Operation {
	return huma.Operation{
		OperationID: "usage-get-all",
		Method:      http.MethodGet,
		Path:        "/usages",
		Summary:     "Get all usages",
		Tags:        []string{"usage"},
	}
}

func UsageGetAll(plane plane.Plane) func(ctx context.Context, input *struct{}) (*struct{ Body []UsageDto }, error) {
	return func(ctx context.Context, input *struct{}) (*struct{ Body []UsageDto }, error) {
		usages := plane.Usage().GetAll()

		result := make([]UsageDto, len(usages))
		for i, usage := range usages {
			result[i] = usageToDto(usage)
		}

		return &struct{ Body []UsageDto }{Body: result}, nil
	}
}

func UsageGetOperation() huma.Operation {
	return huma.Operation{
		OperationID: "usage-get",
		Method:      http.MethodGet,
		Path:        "/usages/{name}",
		Summary:     "Get usage by name",
		Tags:        []string{"usage"},
	}
}

func UsageGet(plane plane.Plane) func(ctx context.Context, input *UsageGetInput) (*struct{ Body UsageDto }, error) {
	return func(ctx context.Context, input *UsageGetInput) (*struct{ Body UsageDto }, error) {
		usage, err := plane.Usage().Get(input.Name)

		if err != nil {
			return nil, huma.Error404NotFound("usage not found", err)
		}

		result := usageToDto(usage)

		return &struct{ Body UsageDto }{Body: result}, nil
	}
}
