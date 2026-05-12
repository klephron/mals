package handler

import (
	"mals/pkg/config"
	"mals/third_party/lsp"
)

func (s *Handler) InitializedDefault(params *lsp.InitializedParams) error {
	s.resources.Range(func(key string, value *config.HandlerLspResource) bool {
		scope, err := s.getResourceScope(value.Scope)
		if err != nil {
			s.plane.Errorf("%T %v: Initialized %v", s, s.Name(), err)
			return true
		}

		switch vs := value.Spec.(type) {
		case *config.HandlerLspResourceSpecLsp:
			lspName, token, err := s.plane.Scope().LspAcquire(vs.Name, scope)
			if err != nil {
				s.plane.Errorf("%T %v: Initialized %v", s, s.Name(), err)
				return true
			}

			defer func() {
				if err := s.plane.Scope().LspRelease(lspName, token); err != nil {
					s.plane.Errorf("%T %v: Initialized %v", s, s.Name(), err)
				}
			}()

			lspParams := &lsp.InitializedParams{}

			err = s.plane.Lsp().Initialized(lspName, lspParams)
			if err != nil {
				s.plane.Errorf("%T %v: Initialized %v", s, s.Name(), err)
				return true
			}

			s.plane.Debugf("%T %v: Initialized %T", s, s.Name(), vs)

		case *config.HandlerLspResourceSpecModel:
		default:
			s.plane.Errorf("%T %v: Initialized unexpected spec %T", s, s.Name(), vs)
		}
		return true
	})

	return nil
}

func (s *Handler) Initialized(params *lsp.InitializedParams) error {
	if *s.endpoints.Initialized.Default {
		err := s.InitializedDefault(params)
		if err != nil {
			return err
		}
	}

	s.plane.Infof("%T %v: Initialized done", s, s.Name())

	return nil
}
