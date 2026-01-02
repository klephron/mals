package scope

type Space struct {
	name string
}

func NewSpace(name string) Space {
	return Space{
		name: name,
	}
}

func (s *Space) Name() string {
	return s.name
}

type Scope struct {
	path []Space
}

func NewScope(spaces ...Space) *Scope {
	return &Scope{
		path: spaces,
	}
}

func (s *Scope) Path() []Space {
	return s.path
}

func NewScopeGlobal() *Scope {
	return NewScope()
}

func NewScopeClient(name string) *Scope {
	return NewScope(NewSpace(name))
}
