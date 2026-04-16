package config

type Listener struct {
	Name     string
	Ipc      ListenerIpc
	Protocol ListenerProtocol
}

type ListenerIpc interface {
	ListenerIpc() string
}

type ListenerIpcTcp struct {
	Port int
}

func (s *ListenerIpcTcp) ListenerIpc() string {
	return "tcp"
}

type ListenerProtocol interface {
	ListenerProtocol() string
}

type ListenerProtocolApi struct {
}

func (s *ListenerProtocolApi) ListenerProtocol() string {
	return "api"
}

type ListenerProtocolLsp struct {
	Handlers []ListenerProtocolLspHandler
}

func (s *ListenerProtocolLsp) ListenerProtocol() string {
	return "lsp"
}

type ListenerProtocolLspHandler struct {
	Name      string
	Condition ListenerProtocolLspHandlerCondition
	Handler   string
}

type ListenerProtocolLspHandlerCondition struct {
	Filetypes []string
	Paths     []string
}
