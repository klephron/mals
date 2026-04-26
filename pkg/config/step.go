package config

type Step struct {
	Name       *string
	Assign     *string
	Definition StepDefinition
}

type StepDefinition interface {
	StepDefinitionKind() string
}

// lsp/completion
type StepLspCompletion struct {
	Resource *string
}

func (s *StepLspCompletion) StepDefinitionKind() string {
	return "lsp/completion"
}

// json/dumps
type StepJsonDumps struct {
	Input *string
}

func (s *StepJsonDumps) StepDefinitionKind() string {
	return "json/dumps"
}

// json/parse
type StepJsonParse struct {
	Input *string
}

func (s *StepJsonParse) StepDefinitionKind() string {
	return "json/parse"
}

// json/parse/completion
type StepJsonParseCompletion struct {
	Input *string
}

func (s *StepJsonParseCompletion) StepDefinitionKind() string {
	return "json/parse/completion"
}

// model/raw
type StepModelRaw struct {
	Resource *string
	Prompt   *string
}

func (s *StepModelRaw) StepDefinitionKind() string {
	return "model/raw"
}

// model
type StepModel struct {
	Resource *string
	Prompt   *string
}

func (s *StepModel) StepDefinitionKind() string {
	return "model"
}

// return
type StepReturn struct {
	Input *string
}

func (s *StepReturn) StepDefinitionKind() string {
	return "return"
}

// if
type StepIf struct {
	Condition *string
	Then      []*Step
	Else      []*Step
}

func (s *StepIf) StepDefinitionKind() string {
	return "if"
}

// for
type StepFor struct {
	Condition *string
	Do        []*Step
}

func (s *StepFor) StepDefinitionKind() string {
	return "for"
}
