package common

import (
	"context"
)

type Listener interface {
	Type() string
	ListenAndServe(ctx context.Context) error
}
