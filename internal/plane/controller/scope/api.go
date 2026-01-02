package scope

import (
	"fmt"
	"mals/internal/plane/controller"
	"mals/internal/scope"
	"mals/pkg/config"

	"github.com/puzpuzpuz/xsync/v4"
)

func scopeDescent(current *Space, scope *scope.Scope, create bool) *Space {
	for _, sc := range scope.Path() {
		next, ok := current.children.Load(sc)
		if !ok {
			if create {
				next = newSpace(sc)
				current.children.Store(sc, next)
			} else {
				return nil
			}
		}
		current = next
	}
	return current
}

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
	current := scopeDescent(s.state.root, scope, true)

	resource, exist := current.models.Load(name)

	if !exist {
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

		s.plane.Infof("ModelAcquire new stored %v", resource.fullname)
	}

	resource.rw.Lock()
	defer resource.rw.Unlock()

	if !exist {
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

	current := scopeDescent(s.state.root, scope, false)
	if current == nil {
		return fmt.Errorf("scope %v path not found", scope)
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

func (s *ScopeController) ScopeLspRegister(config config.Lsp) error {
	name := config.Name

	if _, ok := s.state.lsps.Load(name); ok {
		return fmt.Errorf("lsp %v exists", name)
	}

	s.state.lsps.Store(name, &config)

	return nil
}

func (s *ScopeController) ScopeLspAcquire(name string, scope *scope.Scope) (string, controller.ScopeToken, error) {
	current := scopeDescent(s.state.root, scope, true)

	resource, exist := current.lsps.Load(name)
	if !exist {
		config, ok := s.state.lsps.Load(name)

		if !ok {
			return "", nil, fmt.Errorf("lsp %v does not exist", name)
		}

		r := &ResourceLsp{
			fullname:     mangleName(name, scope),
			config:       config,
			dependencies: xsync.NewMap[string, *ScopeToken](),
		}

		resource = r
		current.lsps.Store(name, resource)

		s.plane.Infof("LspAcquire new stored %v", resource.fullname)
	}

	resource.rw.Lock()
	defer resource.rw.Unlock()

	if !exist {
		if err := s.plane.LspRegister(resource.fullname, resource.config); err != nil {
			return "", nil, err
		}
		if err := s.plane.LspCreate(resource.fullname); err != nil {
			return "", nil, err
		}
		if err := s.plane.LspStart(resource.fullname); err != nil {
			return "", nil, err
		}
	}

	token := newToken()
	resource.dependencies.Store(token.Token(), token)

	s.plane.Infof("LspAcquire %v in scope %v token %v assigned", resource.fullname, scope, token)

	return resource.fullname, token, nil
}

func (s *ScopeController) ScopeLspRelease(fullname string, token controller.ScopeToken) error {
	name, scope := unmangleName(fullname)

	current := scopeDescent(s.state.root, scope, false)
	if current == nil {
		return fmt.Errorf("scope %v path not found", scope)
	}

	resource, ok := current.lsps.Load(name)
	if !ok {
		return fmt.Errorf("resource %v in scope %v not found", name, scope)
	}

	resource.rw.Lock()
	defer resource.rw.Unlock()

	resource.dependencies.Delete(token.Token())

	s.plane.Infof("LspRelease %v in scope %v token %v released", resource.fullname, scope, token)

	return nil
}

func (s *ScopeController) scopeClose(errors *[]error, current *Space) {
	current.lsps.Range(func(name string, resource *ResourceLsp) bool {
		resource.rw.Lock()
		defer resource.rw.Unlock()

		if resource.dependencies.Size() > 0 {
			err := fmt.Errorf("cannot close scope, resource %v has %d active tokens",
				name, resource.dependencies.Size())

			*errors = append(*errors, err)

			return true
		}

		if err := s.plane.LspStop(resource.fullname); err != nil {
			*errors = append(*errors, err)
			return true
		}
		if err := s.plane.LspDelete(resource.fullname); err != nil {
			*errors = append(*errors, err)
			return true
		}
		if err := s.plane.LspUnregister(resource.fullname); err != nil {
			*errors = append(*errors, err)
			return true
		}

		current.lsps.Delete(name)

		s.plane.Infof("closed lsp %v", resource.fullname)

		return true
	})

	current.models.Range(func(name string, resource *ResourceModel) bool {
		resource.rw.Lock()
		defer resource.rw.Unlock()

		if resource.dependencies.Size() > 0 {
			err := fmt.Errorf("cannot close scope, resource %v has %d active tokens",
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

		s.plane.Infof("close model %v", resource.fullname)

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
