package server

import "context"

type LspServer interface {
	Name() string
	Kind() string
	Run(ctx context.Context) error
}
