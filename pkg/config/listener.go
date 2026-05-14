package config

type Listener struct {
	Name     string
	Ipc      ListenerIpc
	Protocol ListenerProtocol
}

type ListenerIpc interface {
	ListenerIpcKind() string
}

type ListenerIpcTcp struct {
	Port *int32
}

func (s *ListenerIpcTcp) ListenerIpcKind() string {
	return "tcp"
}

type ListenerProtocol interface {
	ListenerProtocolKind() string
}

type ListenerProtocolApi struct {
}

func (s *ListenerProtocolApi) ListenerProtocolKind() string {
	return "api"
}

type ListenerProtocolLsp struct {
	Handlers []*ListenerProtocolLspHandler
}

func (s *ListenerProtocolLsp) ListenerProtocolKind() string {
	return "lsp"
}

type ListenerProtocolLspHandler struct {
	Name    string
	Handler string
}

func (s *Listener) Default() {
	switch cp := s.Protocol.(type) {
	case *ListenerProtocolLsp:
		if cp.Handlers == nil {
			cp.Handlers = make([]*ListenerProtocolLspHandler, 0)
		}
	case *ListenerProtocolApi:
	}
}
