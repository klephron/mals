package handlers

import (
	"context"
	"mals/internal/lsp/protocol"
	"mals/internal/plane"
	"mals/internal/plane/controller"
	"mals/pkg/wire"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type LspDto struct {
	Name         string                       `json:"name"`
	Status       controller.LspStatus         `json:"status"`
	StatusFlags  []string                     `json:"status_flags"`
	Config       *wire.Lsp                    `json:"config"`
	Capabilities *protocol.ServerCapabilities `json:"capabilities"`
	Info         *protocol.ServerInfo         `json:"info"`
}

type LspGetInput struct {
	Name string `path:"name" doc:"LSP name"`
}

var (
	lspStatusFlags = map[controller.LspStatus]string{
		controller.LspAbsent:     "ABSENT",
		controller.LspRegistered: "REGISTERED",
		controller.LspCreated:    "CREATED",
		controller.LspStarted:    "STARTED",
	}
)

func lspToDto(lsp *controller.LspData) LspDto {
	statusFlags := make([]string, 0, 4)
	for bit, name := range lspStatusFlags {
		if lsp.Status&bit != 0 {
			statusFlags = append(statusFlags, name)
		}
	}

	t := LspDto{
		Name:         lsp.Name,
		Status:       lsp.Status,
		StatusFlags:  statusFlags,
		Config:       &wire.Lsp{},
		Capabilities: lsp.Capabilities,
		Info:         lsp.Info,
	}
	t.Config.Wire(lsp.Config)

	return t
}

func LspGetAllOperation() huma.Operation {
	return huma.Operation{
		OperationID: "lsp-get-all",
		Method:      http.MethodGet,
		Path:        "/lsps",
		Summary:     "Get all LSPs",
		Tags:        []string{"lsp"},
	}
}

func LspGetAll(plane plane.Plane) func(ctx context.Context, input *struct{}) (*struct{ Body []LspDto }, error) {
	return func(ctx context.Context, input *struct{}) (*struct{ Body []LspDto }, error) {
		lsps := plane.LspGetAll()

		result := make([]LspDto, len(lsps))
		for i, lsp := range lsps {
			result[i] = lspToDto(lsp)
		}

		return &struct{ Body []LspDto }{Body: result}, nil
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

func LspGet(plane plane.Plane) func(ctx context.Context, input *LspGetInput) (*struct{ Body LspDto }, error) {
	return func(ctx context.Context, input *LspGetInput) (*struct{ Body LspDto }, error) {
		lsp, err := plane.LspGet(input.Name)

		if err != nil {
			return nil, huma.Error404NotFound("LSP not found", err)
		}

		result := lspToDto(lsp)

		return &struct{ Body LspDto }{Body: result}, nil
	}
}
