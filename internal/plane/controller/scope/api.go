package scope

import (
	"fmt"
	"mals/internal/plane/controller"
	"mals/internal/scope"
	"mals/pkg/config"

	"github.com/puzpuzpuz/xsync/v4"
)

func (s *ScopeController) Shutdown() error {
	s.ScopeClose(scope.NewScopeGlobal())

	s.state.statusRW.RLock()
	cancel := s.state.statusCancel
	s.state.statusRW.RUnlock()

	cancel()

	return nil
}

func (s *ScopeController) ScopeModelRegister(config config.Model) error {
	name := config.Name

	if _, ok := s.state.models.Load(name); ok {
		return fmt.Errorf("model %v exists", name)
	}

	s.state.models.Store(name, &config)

	return nil
}

func (s *ScopeController) ScopeModelAcquire(name string, scope *scope.Scope) (string, controller.ScopeToken, error) {
	current := s.state.root

	for _, sc := range scope.Path() {
		next, ok := current.children.Load(sc)
		if !ok {
			next = newSpace(sc)
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
			fullname:     mangleName(name, scope),
			config:       config,
			dependencies: xsync.NewMap[string, *ScopeToken](),
		}

		resource = r
		current.models.Store(name, resource)

		s.plane.Infof("ModelAcquire new model stored %v", resource.fullname)
	}

	resource.rw.Lock()
	defer resource.rw.Unlock()

	if !exists {
		if err := s.plane.ModelRegister(resource.fullname, resource.config); err != nil {
			return "", nil, err
		}
		if err := s.plane.ModelCreate(resource.fullname); err != nil {
			return "", nil, err
		}
		if err := s.plane.ModelStart(resource.fullname); err != nil {
			return "", nil, err
		}
	}

	token := newToken()
	resource.dependencies.Store(token.Token(), token)

	s.plane.Infof("ModelAcquire %v in scope %v token %v assigned", resource.fullname, scope, token)

	return resource.fullname, token, nil
}

func (s *ScopeController) ScopeModelRelease(fullname string, token controller.ScopeToken) error {
	name, scope := unmangleName(fullname)

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

	s.plane.Infof("ModelRelease %v in scope %v token %v released", resource.fullname, scope, token)

	return nil
}

func (s *ScopeController) scopeClose(errors *[]error, current *Space) {
	current.models.Range(func(name string, resource *ResourceModel) bool {
		resource.rw.Lock()
		defer resource.rw.Unlock()

		if resource.dependencies.Size() > 0 {
			err := fmt.Errorf("cannot close scope, resource %s has %d active tokens",
				name, resource.dependencies.Size())

			*errors = append(*errors, err)

			return true
		}

		if err := s.plane.ModelStop(resource.fullname); err != nil {
			*errors = append(*errors, err)
			return true
		}
		if err := s.plane.ModelDelete(resource.fullname); err != nil {
			*errors = append(*errors, err)
			return true
		}
		if err := s.plane.ModelUnregister(resource.fullname); err != nil {
			*errors = append(*errors, err)
			return true
		}

		current.models.Delete(name)

		s.plane.Infof("Close %v", resource.fullname)

		return true
	})
}

func (s *ScopeController) scopeCloseDFS(errors *[]error, current *Space) {
	current.children.Range(func(key scope.Space, value *Space) bool {
		s.scopeCloseDFS(errors, value)
		return true
	})

	s.scopeClose(errors, current)
}

func (s *ScopeController) ScopeClose(scope *scope.Scope) []error {
	errors := make([]error, 0)

	current := s.state.root

	for _, sc := range scope.Path() {
		next, ok := current.children.Load(sc)
		if !ok {
			errors = append(errors, fmt.Errorf("scope %v does not exist", scope.Path()))
			return errors
		}
		current = next
	}

	s.scopeCloseDFS(&errors, current)

	return errors
}
