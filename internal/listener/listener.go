package listener

import (
	"context"
)

type Listener interface {
	Name() string
	Kind() string
	Ipc() string
	Run(ctx context.Context) error
}
