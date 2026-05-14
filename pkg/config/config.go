package config

type Config struct {
	Logs      []*Log
	Models    []*Model
	Lsps      []*Lsp
	Handlers  []*Handler
	Listeners []*Listener
}

func (s *Config) Default() {
	if s.Logs == nil {
		s.Logs = make([]*Log, 0)
	}
	if s.Models == nil {
		s.Models = make([]*Model, 0)
	}
	if s.Lsps == nil {
		s.Lsps = make([]*Lsp, 0)
	}
	if s.Handlers == nil {
		s.Handlers = make([]*Handler, 0)
	}
	if s.Listeners == nil {
		s.Listeners = make([]*Listener, 0)
	}

	for _, log := range s.Logs {
		log.Default()
	}
	for _, model := range s.Models {
		model.Default()
	}
	for _, lsp := range s.Lsps {
		lsp.Default()
	}
	for _, handler := range s.Handlers {
		handler.Default()
	}
	for _, listener := range s.Listeners {
		listener.Default()
	}
}
