package middleware

import (
	"fmt"
	"mals/internal/lsp/protocol"
)

func (s *Middleware) Initialized(params *protocol.InitializedParams) error {
	if s.initialized {
		return fmt.Errorf("%v: already initialized", s.Name())
	}
	s.initialized = true
	return nil
}
