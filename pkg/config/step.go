package config

type Step interface {
	step()
}

type StepGeneric struct {
	Step       `json:"step,omitempty"`
	Name       string
	Conditions []*Condition
	Scope      string
}

type StepModel struct {
	StepGeneric
	Model    string
	Template string
}

type StepLsp struct {
	StepGeneric
	Lsp      string
	Template string
}
