package handlers

import (
	"context"
	"mals/internal/lsp/protocol"
	"mals/internal/plane"
	"mals/internal/plane/controller"
	"mals/pkg/config"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type LSP struct {
	Name         string                       `json:"name"`
	Status       string                       `json:"status"`
	Config       *config.Lsp                  `json:"config"`
	Capabilities *protocol.ServerCapabilities `json:"capabilities"`
	Info         *protocol.ServerInfo         `json:"info"`
}

func toDto(lsp *controller.LspData) LSP {
	return LSP{
		Name:         lsp.Name,
		Status:       string(lsp.Status),
		Config:       lsp.Config,
		Capabilities: lsp.Capabilities,
		Info:         lsp.Info,
	}
}

func LspGetAllOperation() huma.Operation {
	return huma.Operation{
		OperationID: "lsps-get-all",
		Method:      http.MethodGet,
		Path:        "/lsps",
		Summary:     "Get all LSPs",
		Tags:        []string{"lsp"},
	}
}

func LspGetAll(plane plane.Plane) func(ctx context.Context, input *struct{}) (*struct{ Body []LSP }, error) {
	return func(ctx context.Context, input *struct{}) (*struct{ Body []LSP }, error) {
		lsps := plane.LspGetAll()

		result := make([]LSP, len(lsps))
		for i, lsp := range lsps {
			result[i] = toDto(lsp)
		}

		return &struct{ Body []LSP }{Body: result}, nil
	}
}

func LspGetOperation() huma.Operation {
	return huma.Operation{
		OperationID: "lsp-get",
		Method:      http.MethodGet,
		Path:        "/lsps/{name}",
		Summary:     "Get LSP by name",
		Tags:        []string{"lsp"},
	}
}

type LspGetInput struct {
	Name string `path:"name" doc:"LSP name"`
}

func LspGet(plane plane.Plane) func(ctx context.Context, input *LspGetInput) (*struct{ Body LSP }, error) {
	return func(ctx context.Context, input *LspGetInput) (*struct{ Body LSP }, error) {
		lsp, err := plane.LspGet(input.Name)
		result := toDto(lsp)
		return &struct{ Body LSP }{Body: result}, huma.Error404NotFound("LSP not found", err)
	}
}
