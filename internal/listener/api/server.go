package api

import (
	"mals/internal/listener/api/handler"
	"mals/pkg/core"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
)

func (s *ListenerApi) newServer() *http.Server {
	router := gin.New()

	router.GET("/metrics", handler.Metrics())

	// api
	config := huma.DefaultConfig("MALS API", core.MiddlewareVersion)
	config.OpenAPIPath = "/openapi"
	config.DocsPath = "/docs"

	api := humagin.New(router, config)

	huma.Register(api, handler.LspGetAllOperation(), handler.LspGetAll(s.plane))
	huma.Register(api, handler.LspGetOperation(), handler.LspGet(s.plane))

	huma.Register(api, handler.ListenerGetAllOperation(), handler.ListenerGetAll(s.plane))
	huma.Register(api, handler.ListenerGetOperation(), handler.ListenerGet(s.plane))

	huma.Register(api, handler.LogGetAllOperation(), handler.LogGetAll(s.plane))
	huma.Register(api, handler.LogGetOperation(), handler.LogGet(s.plane))

	huma.Register(api, handler.ModelGetAllOperation(), handler.ModelGetAll(s.plane))
	huma.Register(api, handler.ModelGetOperation(), handler.ModelGet(s.plane))

	huma.Register(api, handler.HandlerGetAllOperation(), handler.HandlerGetAll(s.plane))
	huma.Register(api, handler.HandlerGetOperation(), handler.HandlerGet(s.plane))

	huma.Register(api, handler.ScopeTreeRootOperation(), handler.ScopeTreeRoot(s.plane))

	srv := &http.Server{
		Addr:    s.addr,
		Handler: router,
	}

	return srv
}
