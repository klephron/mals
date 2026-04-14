package scope

func NewScopeGlobal() *Scope {
	return NewScope()
}

func NewScopeClient(listener string, client string) *Scope {
	return NewScope(NewSpace(listener), NewSpace(client))
}
