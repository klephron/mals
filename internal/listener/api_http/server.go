package api_http

import (
	"mals/internal/listener/api_http/handlers"
	"mals/pkg/info"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
)

func (s *ListenerApiHttp) newServer() *http.Server {
	router := gin.New()

	router.GET("/metrics", handlers.Metrics())

	// api
	config := huma.DefaultConfig("MALS API", info.MiddlewareVersion)
	config.OpenAPIPath = "/openapi"
	config.DocsPath = "/docs"

	api := humagin.New(router, config)

	huma.Register(api, handlers.LspGetAllOperation(), handlers.LspGetAll(s.plane))
	huma.Register(api, handlers.LspGetOperation(), handlers.LspGet(s.plane))

	srv := &http.Server{
		Addr:    s.addr,
		Handler: router,
	}

	return srv
}
