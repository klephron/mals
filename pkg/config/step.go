package config

type Step struct {
	Name       string
	Assign     string
	Definition StepDefinition
}

type StepDefinition interface {
	Kind() string
}

// lsp/completion
type StepLspCompetion struct {
	Resource string
}

func (s *StepLspCompetion) Kind() string {
	return "lsp/completion"
}

// json/dumps
type StepJsonDumps struct {
	Input string
}

func (s *StepJsonDumps) Kind() string {
	return "json/dumps"
}

// json/parse
type StepJsonParse struct {
	Input string
}

func (s *StepJsonParse) Kind() string {
	return "json/parse"
}

// json/parse/completion
type StepJsonParseCompletion struct {
	Input string
}

func (s *StepJsonParseCompletion) Kind() string {
	return "json/parse/completion"
}

// model/simple
type StepModelSimple struct {
	Resource string
	Prompt   string
}

func (s *StepModelSimple) Kind() string {
	return "model/simple"
}

// model/template
type StepModelTemplate struct {
	Resource string
	Prompt   string
}

func (s *StepModelTemplate) Kind() string {
	return "model/template"
}

// return
type StepReturn struct {
	Input string
}

func (s *StepReturn) Kind() string {
	return "return"
}

// if
type StepIf struct {
	Condition string
	Then      []Step
	Else      []Step
}

func (s *StepIf) Kind() string {
	return "if"
}

// while
type StepWhile struct {
	Condition string
	Do        []Step
	Max       *int
}

func (s *StepWhile) Kind() string {
	return "while"
}
