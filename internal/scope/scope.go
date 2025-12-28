package scope

type ScopeKind int32

const (
	ScopeClient   ScopeKind = iota
	ScopeWorkflow ScopeKind = iota
	ScopeStep     ScopeKind = iota
)

type Scope struct {
	kind ScopeKind
	name string
}

type Path struct {
	scopes []Scope
}

func NewScope(kind ScopeKind, name string) Scope {
	return Scope{
		kind: kind,
		name: name,
	}
}

func NewPath(scopes ...Scope) *Path {
	return &Path{
		scopes: scopes,
	}
}
