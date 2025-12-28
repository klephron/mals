package scope

func NewScopeGlobal() *Scope {
	return &Scope{
		path: make([]Space, 0),
	}
}
