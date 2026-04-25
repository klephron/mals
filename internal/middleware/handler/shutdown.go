package handler

import "mals/pkg/config"

func (s *Handler) ShutdownDefault() error {
	s.resources.Range(func(key string, value *config.HandlerLspResource) bool {
		scope, err := s.getResourceScope(value.Scope)
		if err != nil {
			s.plane.Errorf("%T %v: Shutdown %v", s, s.Name(), err)
			return true
		}

		switch vs := value.Spec.(type) {
		case *config.HandlerLspResourceSpecLsp:
			lspName, token, err := s.plane.Scope().LspAcquire(vs.Name, scope)
			if err != nil {
				s.plane.Errorf("%T %v: Shutdown %v", s, s.Name(), err)
				return true
			}

			defer func() {
				if err := s.plane.Scope().LspRelease(lspName, token); err != nil {
					s.plane.Errorf("%T %v: Shutdown %v", s, s.Name(), err)
				}
			}()

			err = s.plane.Lsp().Shutdown(lspName)
			if err != nil {
				s.plane.Errorf("%T %v: Shutdown %v", s, s.Name(), err)
				return true
			}

			s.plane.Debugf("%T %v: Shutdown %T", s, s.Name(), vs)

		case *config.HandlerLspResourceSpecModel:
		default:
			s.plane.Errorf("%T %v: Shutdown unexpected spec %T", s, s.Name(), vs)
		}

		return true
	})

	return nil
}

func (s *Handler) Shutdown() error {
	if *s.endpoints.Shutdown.Default {
		err := s.ShutdownDefault()
		if err != nil {
			return err
		}
	}

	s.plane.Infof("%T %v: Shutdown done", s, s.Name())

	scope, err := s.getResourceScope(config.HandlerLspResourceScopeHandler)
	if err != nil {
		s.plane.Errorf("%T %v: Shutdown %v", s, s.Name(), err)
		return err
	}

	s.plane.Scope().Close(scope)

	return nil
}
