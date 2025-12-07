package config

type Condition struct {
	Filetypes []string
	Paths     []string
	Events    []string
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

type Config struct {
	Loggers   []Log
	Listeners []Listener
	Models    []*Model
	Lsps      []*Lsp
	Usages    []*Usage
}
