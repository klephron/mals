package config

type ModelSpec interface {
	modelspec()
}

type Model struct {
	Name     *string   `json:"name"`
	Spec     *string   `json:"spec"`
	Settings ModelSpec `json:"settings"`
}

type LspSpec interface {
	lspspec()
}

type Lsp struct {
	Name     *string `json:"name"`
	Spec     *string `json:"spec"`
	Settings LspSpec `json:"settings"`
}

type Condition struct {
	Filetypes []string `json:"filetypes" default:"[]"`
	Paths     []string `json:"paths" default:"[]"`
	Types     []string `json:"types" default:"[]"`
}

type Step interface {
	step()
}

type StepGeneric struct {
	Step
	Name       *string      `json:"name"`
	Conditions []*Condition `json:"conditions"`
	Scope      *string      `json:"scope"`
}

type StepModel struct {
	StepGeneric
	Model    *string `json:"model"`
	Template *string `json:"template"`
}

type StepLsp struct {
	StepGeneric
	Lsp      *string `json:"lsp"`
	Template *string `json:"template"`
}

type Workflow struct {
	Name  *string `json:"name"`
	Steps []Step  `json:"steps" default:"[]"`
}

type Usage struct {
	Name       *string      `json:"name"`
	Conditions []*Condition `json:"conditions" default:"[]"`
	Workflow   Workflow     `json:"workflow"`
}

type Config struct {
	Models []*Model `json:"models" default:"[]"`
	Lsps   []*Lsp   `json:"lsps" default:"[]"`
	Usages []*Usage `'json:"usages" default:"[]"`
}
