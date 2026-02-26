package handlers

import (
	"context"
	"mals/internal/plane"
	"mals/internal/plane/controller"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type SpaceDto struct {
	Name     string                      `json:"name"`
	Children map[string]SpaceDto         `json:"children"`
	Lsps     map[string]ResourceLspDto   `json:"lsps"`
	Models   map[string]ResourceModelDto `json:"models"`
}

type ResourceLspDto struct {
	Fullname     string                   `json:"fullname"`
	Dependencies map[string]ScopeTokenDto `json:"dependencies"`
}

type ResourceModelDto struct {
	Fullname     string                   `json:"fullname"`
	Dependencies map[string]ScopeTokenDto `json:"dependencies"`
}

type ScopeTokenDto struct {
	From string `json:"from"`
}

func spaceToDto(s *controller.Space) SpaceDto {
	if s == nil {
		return SpaceDto{}
	}

	dto := SpaceDto{
		Name:     s.Space.Name(),
		Children: make(map[string]SpaceDto),
		Lsps:     make(map[string]ResourceLspDto),
		Models:   make(map[string]ResourceModelDto),
	}

	for key, child := range s.Children {
		dto.Children[key.Name()] = spaceToDto(child)
	}

	for key, lsp := range s.Lsps {
		deps := make(map[string]ScopeTokenDto)
		for k, v := range lsp.Dependencies {
			deps[k] = scopeTokenToDto(v)
		}
		dto.Lsps[key] = ResourceLspDto{
			Fullname:     lsp.Fullname,
			Dependencies: deps,
		}
	}

	for key, model := range s.Models {
		deps := make(map[string]ScopeTokenDto)
		for k, v := range model.Dependencies {
			deps[k] = scopeTokenToDto(v)
		}
		dto.Models[key] = ResourceModelDto{
			Fullname:     model.Fullname,
			Dependencies: deps,
		}
	}

	return dto
}

func scopeTokenToDto(t controller.ScopeToken) ScopeTokenDto {
	return ScopeTokenDto{
		From: t.From().Format("2006-01-02T15:04:05.000"),
	}
}

func ScopeTreeRootOperation() huma.Operation {
	return huma.Operation{
		OperationID: "scope-tree-root",
		Method:      http.MethodGet,
		Path:        "/scopes",
		Summary:     "Get scope tree from root",
		Tags:        []string{"scope"},
	}
}

func ScopeTreeRoot(plane plane.Plane) func(ctx context.Context, input *struct{}) (*struct{ Body SpaceDto }, error) {
	return func(ctx context.Context, input *struct{}) (*struct{ Body SpaceDto }, error) {
		root := plane.ScopeTreeRoot()

		result := spaceToDto(root)

		return &struct{ Body SpaceDto }{Body: result}, nil
	}
}
