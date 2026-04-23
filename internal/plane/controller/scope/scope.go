package scope

import (
	"mals/internal/plane/controller"
	"time"

	"github.com/google/uuid"
)

type scopeToken struct {
	controller.ScopeToken
	id   string
	from time.Time
}

func (s *scopeToken) Token() string {
	return s.id
}

func (s *scopeToken) From() time.Time {
	return s.from
}

func newScopeToken() *scopeToken {
	return &scopeToken{
		id:   uuid.New().String(),
		from: time.Now(),
	}
}
