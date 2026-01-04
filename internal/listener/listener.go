package listener

import (
	"context"
)

type Listener interface {
	Name() string
	Run(ctx context.Context) error
}
