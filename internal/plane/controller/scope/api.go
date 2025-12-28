package scope

import (
	"fmt"
	"mals/internal/plane/controller"
	"mals/internal/scope"
	"mals/pkg/config"

	"github.com/puzpuzpuz/xsync/v4"
)

func (s *ScopeController) Shutdown() error {
	s.state.statusRW.RLock()
	cancel := s.state.statusCancel
	s.state.statusRW.RUnlock()

	cancel()

	return nil
}

func (s *ScopeController) ModelRegister(config config.Model) error {
	name := config.Name

	if _, ok := s.state.models.Load(name); ok {
		return fmt.Errorf("model %v exists", name)
	}

	s.state.models.Store(name, &config)

	return nil
}

func (s *ScopeController) ModelAcquire(name string, scope *scope.Scope) (string, controller.ScopeToken, error) {
	current := s.state.root

	for _, sc := range scope.Path() {
		next, ok := current.children.Load(sc)
		if !ok {
			next = NewSpace(sc)
			current.children.Store(sc, next)
		}
		current = next
	}

	resource, exists := current.models.Load(name)
	if !exists {
		config, ok := s.state.models.Load(name)

		if !ok {
			return "", nil, fmt.Errorf("model %v does not exist", name)
		}

		r := &ResourceModel{
			fullname:     generateName(name, scope),
			config:       config,
			dependencies: xsync.NewMap[string, *ScopeToken](),
		}

		current.models.Store(name, r)
		resource = r
	}

	resource.rw.Lock()
	defer resource.rw.Unlock()

	empty := true
	resource.dependencies.Range(func(key string, value *ScopeToken) bool {
		empty = false
		return false
	})

	if empty {
		if err := s.plane.Model().Register(resource.fullname, resource.config); err != nil {
			return "", nil, err
		}
		if err := s.plane.Model().Create(resource.fullname); err != nil {
			return "", nil, err
		}
		if err := s.plane.Model().Start(resource.fullname); err != nil {
			return "", nil, err
		}
	}

	token := newToken()
	resource.dependencies.Store(token.Token(), token)

	return resource.fullname, token, nil
}

func (s *ScopeController) ModelRelease(fullname string, token controller.ScopeToken) error {
	name, scope := decodeName(fullname)

	current := s.state.root

	for _, sc := range scope.Path() {
		next, ok := current.children.Load(sc)
		if !ok {
			return fmt.Errorf("scope %v path not found", scope)
		}
		current = next
	}

	resource, ok := current.models.Load(name)
	if !ok {
		return fmt.Errorf("resource %v in scope %v not found", name, scope)
	}

	resource.rw.Lock()
	defer resource.rw.Unlock()

	resource.dependencies.Delete(token.Token())

	empty := true
	resource.dependencies.Range(func(key string, value *ScopeToken) bool {
		empty = false
		return false
	})

	if empty {
		if err := s.plane.Model().Stop(fullname); err != nil {
			return err
		}
		if err := s.plane.Model().Delete(fullname); err != nil {
			return err
		}
		if err := s.plane.Model().Unregister(fullname); err != nil {
			return err
		}
		current.models.Delete(name)
	}

	return nil
}
