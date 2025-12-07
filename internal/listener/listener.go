package listener

import (
	"context"
)

type Listener interface {
	Name() string
	Kind() string
	Ipc() string
	Listen(ctx context.Context) error
}
