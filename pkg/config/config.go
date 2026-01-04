package config

type Condition struct {
	Filetypes []string
	Paths     []string
}

type Workflow struct {
	Steps []*Step
}

type Usage struct {
	Name       string
	Events     []Event
	Conditions []*Condition
	Workflow   *Workflow
}

type Config struct {
	Logs      []*Log
	Models    []*Model
	Lsps      []*Lsp
	Usages    []*Usage
	Listeners []*Listener
}
