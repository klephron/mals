package config

type ModelSpec interface {
	modelspec()
}

type Model struct {
	Name     string
	Settings ModelSpec
}

type LspSpec interface {
	lspspec()
}

type Lsp struct {
	Name     string
	Settings LspSpec
}

type Condition struct {
	Filetypes []string
	Paths     []string
	Types     []string
}

type Step interface {
	step()
}

type StepGeneric struct {
	Step
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

type Workflow struct {
	Name  string
	Steps []Step
}

type Usage struct {
	Name       string
	Conditions []*Condition
	Workflow   *Workflow
}

type Log interface {
	log()
}

type LogGeneric struct {
	Log
	Name string
}

type Listener struct {
	Name string
	Type string
	Port int
}

type Config struct {
	Loggers   []Log
	Listeners []*Listener
	Models    []*Model
	Lsps      []*Lsp
	Usages    []*Usage
}
