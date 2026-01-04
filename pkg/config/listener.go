package config

type Listener struct {
	Name string
	Kind ListenerKind
	Ipc  ListenerIpc
}

type ListenerKind interface {
	Kind() string
}

type ListenerKindApi struct {
}

func (s *ListenerKindApi) Kind() string {
	return "api"
}

type ListenerKindLsp struct {
	Usages []string
}

func (s *ListenerKindLsp) Kind() string {
	return "lsp"
}

type ListenerIpc interface {
	Ipc() string
}

type ListenerIpcStdio struct {
}

func (s *ListenerIpcStdio) Ipc() string {
	return "stdio"
}

type ListenerIpcTcp struct {
	Port int
}

func (s *ListenerIpcTcp) Ipc() string {
	return "tcp"
}

type ListenerIpcHttp struct {
	Port int
}

func (s *ListenerIpcHttp) Ipc() string {
	return "http"
}
