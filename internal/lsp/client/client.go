package client

import "context"

type Client interface {
	Name() string
	ServeBinded(ctx context.Context) error
	Close() error
}
