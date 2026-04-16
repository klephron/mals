package config

type Step struct {
	Name       string
	Assign     string
	Definition StepDefinition
}

type StepDefinition interface {
	StepDefinition() string
}

// lsp/completion
type StepLspCompetion struct {
	Resource string
}

func (s *StepLspCompetion) StepDefinition() string {
	return "lsp/completion"
}

// json/dumps
type StepJsonDumps struct {
	Input string
}

func (s *StepJsonDumps) StepDefinition() string {
	return "json/dumps"
}

// json/parse
type StepJsonParse struct {
	Input string
}

func (s *StepJsonParse) StepDefinition() string {
	return "json/parse"
}

// json/parse/completion
type StepJsonParseCompletion struct {
	Input string
}

func (s *StepJsonParseCompletion) StepDefinition() string {
	return "json/parse/completion"
}

// model/simple
type StepModelSimple struct {
	Resource string
	Prompt   string
}

func (s *StepModelSimple) StepDefinition() string {
	return "model/simple"
}

// model/template
type StepModelTemplate struct {
	Resource string
	Prompt   string
}

func (s *StepModelTemplate) StepDefinition() string {
	return "model/template"
}

// return
type StepReturn struct {
	Input string
}

func (s *StepReturn) StepDefinition() string {
	return "return"
}

// if
type StepIf struct {
	Condition string
	Then      []Step
	Else      []Step
}

func (s *StepIf) StepDefinition() string {
	return "if"
}

// while
type StepWhile struct {
	Condition string
	Do        []Step
	Max       *int
}

func (s *StepWhile) StepDefinition() string {
	return "while"
}
