package listener

import (
	"context"

	"github.com/gin-gonic/gin"
)

type Listener interface {
	Name() string
	Run(ctx context.Context) error
}

func init() {
	gin.SetMode(gin.ReleaseMode)
}
