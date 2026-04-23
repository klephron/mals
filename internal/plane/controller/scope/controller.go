package scope

import (
	"context"
	"fmt"
	"mals/internal/plane"
	"mals/internal/scope"
	"mals/pkg/config"
	"strings"
	"sync"

	"github.com/puzpuzpuz/xsync/v4"
)

type state struct {
	statusRW     sync.RWMutex
	statusCancel context.CancelFunc

	lsps   *xsync.Map[string, *config.Lsp]
	models *xsync.Map[string, *config.Model]

	root *stateSpace
}

type stateSpace struct {
	space scope.Space

	children *xsync.Map[scope.Space, *stateSpace]

	lsps   *xsync.Map[string, *stateLsp]
	models *xsync.Map[string, *stateModel]
}

type stateLsp struct {
	rw           sync.RWMutex
	fullname     string
	config       *config.Lsp
	dependencies *xsync.Map[string, *scopeToken]
}

type stateModel struct {
	rw           sync.RWMutex
	fullname     string
	config       *config.Model
	dependencies *xsync.Map[string, *scopeToken]
}

type ScopeController struct {
	state state
	plane plane.Plane
}

func newStateSpace(space scope.Space) *stateSpace {
	return &stateSpace{
		space:    space,
		children: xsync.NewMap[scope.Space, *stateSpace](),
		lsps:     xsync.NewMap[string, *stateLsp](),
		models:   xsync.NewMap[string, *stateModel](),
	}
}

func newStateSpaceRoot() *stateSpace {
	return newStateSpace(scope.NewSpace(""))
}

func New(plane plane.Plane) *ScopeController {
	return &ScopeController{
		state: state{
			statusRW:     sync.RWMutex{},
			statusCancel: nil,

			lsps:   xsync.NewMap[string, *config.Lsp](),
			models: xsync.NewMap[string, *config.Model](),

			root: newStateSpaceRoot(),
		},
		plane: plane,
	}
}

func nameMangle(name string, scope *scope.Scope) string {
	fullname := ""

	for _, sp := range scope.Path() {
		fullname += sp.Name()
		fullname += "#"
	}

	fullname += name

	return fullname
}

func nameUnmangle(name string) (string, *scope.Scope) {
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

func (s *ScopeController) ControllerRun(onReady func()) error {
	s.state.statusRW.Lock()

	if s.state.statusCancel != nil {
		s.state.statusRW.Unlock()

		err := fmt.Errorf("%T is already serving", s)
		s.plane.Errorf("%T: %v", s, err)
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.state.statusCancel = cancel
	s.state.statusRW.Unlock()

	onReady()
	<-ctx.Done()

	s.state.statusRW.Lock()
	s.state.statusCancel = nil
	s.state.statusRW.Unlock()

	return nil
}

func (s *ScopeController) ControllerShutdown() error {
	s.Close(scope.NewScopeGlobal())

	s.state.statusRW.RLock()
	cancel := s.state.statusCancel
	s.state.statusRW.RUnlock()

	cancel()

	return nil
}
