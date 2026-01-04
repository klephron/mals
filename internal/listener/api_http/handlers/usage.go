package handlers

import (
	"context"
	"mals/internal/plane"
	"mals/pkg/config"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type UsageDto config.WireUsage

type UsageGetInput struct {
	Name string `path:"name" doc:"usage name"`
}

func usageToDto(usage *config.Usage) UsageDto {
	return UsageDto(usage.Wire())
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
		usages := plane.UsageGetAll()

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
		usage, err := plane.UsageGet(input.Name)

		if err != nil {
			return nil, huma.Error404NotFound("usage not found", err)
		}

		result := usageToDto(usage)

		return &struct{ Body UsageDto }{Body: result}, nil
	}
}
