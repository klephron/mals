package scope

type SpaceKind int32

const (
	SpaceClient   SpaceKind = iota
	SpaceWorkflow SpaceKind = iota
	SpaceStep     SpaceKind = iota
)

type Space struct {
	kind SpaceKind
	name string
}

type Scope struct {
	spaces []Space
}

func NewSpace(kind SpaceKind, name string) Space {
	return Space{
		kind: kind,
		name: name,
	}
}

func NewScope(spaces ...Space) *Scope {
	return &Scope{
		spaces: spaces,
	}
}

func (s *Scope) IsGlobal() bool {
	return len(s.spaces) == 0
}
