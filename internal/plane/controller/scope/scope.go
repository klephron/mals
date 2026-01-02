package scope

import (
	"mals/internal/plane/controller"
	"mals/internal/scope"
	"strings"

	"github.com/google/uuid"
)

type ScopeToken struct {
	controller.ScopeToken
	id string
}

func (s *ScopeToken) Token() string {
	return s.id
}

func newToken() *ScopeToken {
	return &ScopeToken{
		id: uuid.New().String(),
	}
}

func mangleName(name string, scope *scope.Scope) string {
	fullname := ""

	for _, sp := range scope.Path() {
		fullname += sp.Name()
		fullname += "#"
	}

	fullname += name

	return fullname
}

func unmangleName(name string) (string, *scope.Scope) {
	parts := strings.Split(name, "#")
	if len(parts) == 0 {
		return "", nil
	}

	baseName := parts[len(parts)-1]

	if len(parts) == 1 {
		return baseName, scope.NewScope()
	}

	scopePath := make([]scope.Space, len(parts)-1)
	for i := 0; i < len(parts)-1; i++ {
		scopePath[i] = scope.NewSpace(parts[i])
	}

	return baseName, scope.NewScope(scopePath...)
}
