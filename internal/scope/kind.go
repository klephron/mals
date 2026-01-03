package scope

func NewScopeGlobal() *Scope {
	return NewScope()
}

func NewScopeClient(client string) *Scope {
	return NewScope(NewSpace(client))
}

func NewScopeWorkflow(client string, workflow string) *Scope {
	return NewScope(NewSpace(client), NewSpace(workflow))
}

func NewScopeStep(client string, workflow string, step string) *Scope {
	return NewScope(NewSpace(client), NewSpace(workflow), NewSpace(step))
}
