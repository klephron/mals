package config

type Step struct {
	Name       string
	Conditions []*Condition
	Kind       StepKind
	Scope      ScopeKind
}

type StepKind interface {
	Kind() string
}

type StepKindModel struct {
	Name string
}

func (s *StepKindModel) Kind() string {
	return "model"
}

type StepKindLsp struct {
	Name string
}

func (s *StepKindLsp) Kind() string {
	return "lsp"
}
