package listener

import (
	"context"
)

type Listener interface {
	Type() string
	Name() string
	Listen(ctx context.Context) error
	Listening() bool
}
