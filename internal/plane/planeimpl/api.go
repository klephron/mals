package planeimpl

import "mals/internal/plane/controller"

func (s *Plane) Listener() controller.ListenerController {
	return s.listener
}

func (s *Plane) Log() controller.LogController {
	return s.log
}

func (s *Plane) Lsp() controller.LspController {
	return s.lsp
}

func (s *Plane) Model() controller.ModelController {
	return s.model
}

func (s *Plane) Scope() controller.ScopeController {
	return s.scope
}

func (s *Plane) Handler() controller.HandlerController {
	return s.handler
}

func (s *Plane) Debugf(format string, a ...any) error {
	return s.log.Debugf(format, a...)
}

func (s *Plane) Infof(format string, a ...any) error {
	return s.log.Infof(format, a...)
}

func (s *Plane) Warnf(format string, a ...any) error {
	return s.log.Warnf(format, a...)
}

func (s *Plane) Errorf(format string, a ...any) error {
	return s.log.Errorf(format, a...)
}
