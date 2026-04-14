package config

type Listener struct {
	Name     string
	Ipc      ListenerIpc
	Protocol ListenerProtocol
}

type ListenerIpc interface {
	Kind() string
}

type ListenerIpcTcp struct {
	Port int
}

func (s *ListenerIpcTcp) Kind() string {
	return "tcp"
}

type ListenerProtocol interface {
	Kind() string
}

type ListenerProtocolApi struct {
}

func (s *ListenerProtocolApi) Kind() string {
	return "api"
}

type ListenerProtocolLsp struct {
	Handlers []ListenerProtocolLspHandler
}

func (s *ListenerProtocolLsp) Kind() string {
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
