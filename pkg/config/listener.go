package config

type Listener interface {
	Name() string
	Kind() string
}

type ListenerGeneric struct {
	Listener
	name string
	kind string
}

func (s *ListenerGeneric) Name() string {
	return s.name
}

func (s *ListenerGeneric) Kind() string {
	return s.kind
}

func NewListenerGeneric(name string, kind string) ListenerGeneric {
	return ListenerGeneric{
		name: name,
		kind: kind,
	}
}

type ListenerTcp struct {
	ListenerGeneric
	Port int
}

type ListenerStdio struct {
	ListenerGeneric
}
