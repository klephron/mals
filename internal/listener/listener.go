package listener

import (
	"context"
)

type Listener interface {
	Name() string
	Kind() string
	Ipc() string
	Start(ctx context.Context) error
}
