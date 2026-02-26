package scope

func NewScopeGlobal() *Scope {
	return NewScope()
}

func NewScopeClient(listener string, client string) *Scope {
	return NewScope(NewSpace(listener), NewSpace(client))
}

func NewScopeWorkflow(listener string, client string, workflow string) *Scope {
	return NewScope(NewSpace(listener), NewSpace(client), NewSpace(workflow))
}

func NewScopeStep(listener string, client string, workflow string, step string) *Scope {
	return NewScope(NewSpace(listener), NewSpace(client), NewSpace(workflow), NewSpace(step))
}
