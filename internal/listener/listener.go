package listener

import (
	"context"
)

type Listener interface {
	Type() string
	Listen(ctx context.Context) error
	Listening() bool
}
