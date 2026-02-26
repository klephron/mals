package listener

import (
	"context"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

type Listener interface {
	Name() string
	Run(ctx context.Context) error
}

type ListenerLspClient interface {
	Name() string
	Serve(ctx context.Context) (err error)
}
