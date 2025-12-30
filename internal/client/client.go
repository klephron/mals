package client

import "context"

type Client interface {
	Name() string
	Serve(ctx context.Context) error
}
