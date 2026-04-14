package middleware

import (
	// "encoding/json"
	// "fmt"
	"mals/internal/lsp/protocol"
	// "mals/internal/model"
	// "mals/internal/scope"
	// "mals/internal/util"
	// "mals/pkg/config"
	// "mals/pkg/info"
	// "strings"
	// "sync"
	// "github.com/invopop/jsonschema"
)

// const (
// 	eventCompletionModelTemplate = `
// <|begin_of_text|><|start_header_id|>system<|end_header_id|>

// You are an intelligent autocompletion engine. Your task is to suggest relevant text completions based on document context.

// Rules:
// - Analyze the document context and writing style
// - Generate 3-5 contextually appropriate completions
// - Match the existing tone and style
// - Provide completions of varying lengths (words to phrases)
// - Return ONLY a valid JSON array of strings
// - No explanations, no additional text

// <|eot_id|><|start_header_id|>user<|end_header_id|>

// Document text: The weather today is quite nice. I think I'll go for a walk in the park and maybe stop by the
// Current context: "I think I'll go for a walk in the park and maybe stop by the"

// <|eot_id|><|start_header_id|>assistant<|end_header_id|>

// ["coffee shop", "library", "store to buy groceries", "lake to feed the ducks"]

// <|eot_id|><|start_header_id|>user<|end_header_id|>

// Document text: She opened her laptop and started typing her report. The first chapter discussed the importance of renewable energy sources. Solar panels and wind turbines are becoming more
// Current context: "Solar panels and wind turbines are becoming more"

// <|eot_id|><|start_header_id|>assistant<|end_header_id|>

// ["efficient", "affordable each year", "popular worldwide", "cost-effective solutions", "widely adopted"]

// <|eot_id|><|start_header_id|>user<|end_header_id|>

// Document text:
// %s

// Current context:
// %s

// <|eot_id|><|start_header_id|>assistant<|end_header_id|>
// `
// )

// func eventCompletionModelPrompt(documentContext string, currentContext string) string {
// 	prompt := fmt.Sprintf(eventCompletionModelTemplate, documentContext, currentContext)
// 	return prompt
// }

// func eventCompletionModelSchema[T any]() any {
// 	reflector := jsonschema.Reflector{
// 		AllowAdditionalProperties: false,
// 		DoNotReference:            true,
// 	}
// 	var v T
// 	schema := reflector.Reflect(v)
// 	return schema
// }

// func (s *Middleware) eventCompletionModel(params *protocol.CompletionParams, workspace *Workspace, step *config.Step) (*protocol.CompletionList, error) {
// 	if step.Scope != "global" {
// 		s.plane.Warnf("%T %v: step %T %v scope %v unsupported, set to global",
// 			s, s.Name(), step, step, step.Scope)
// 	}
// 	scope := scope.NewScopeGlobal()

// 	modelName := step.Kind.(*config.StepKindModel).Name

// 	modelKey, token, err := s.plane.Scope().ModelAcquire(modelName, scope)
// 	if err != nil {
// 		s.plane.Errorf("%T %v: step %T %v: %v", s, s.Name(), step, step, err)
// 		return nil, err
// 	}
// 	defer s.plane.Scope().ModelRelease(modelKey, token)

// 	document := s.documentGet(workspace, params.TextDocument.URI)

// 	if document == nil {
// 		return nil, fmt.Errorf("document %v not found", params.TextDocument.URI)
// 	}

// 	task := model.NewTask(
// 		eventCompletionModelPrompt(
// 			document.Text(),
// 			document.LastNChars(params.Position.Line, params.Position.Character, 200),
// 		),
// 		eventCompletionModelSchema[[]string](),
// 		"completion_items",
// 		"Generated completion items",
// 	)

// 	text, err := s.plane.Model().TaskExecClient(modelKey, task, s.clientName)
// 	if err != nil {
// 		s.plane.Errorf("%T %v: step %T %v: %v", s, s.Name(), step, step, err)
// 		return nil, err
// 	}

// 	var resItems []string
// 	if err := json.Unmarshal([]byte(text), &resItems); err != nil {
// 		return nil, fmt.Errorf("step %T %v: %v", step, step, err)
// 	}

// 	items := make([]protocol.CompletionItem, len(resItems))
// 	for i, s := range resItems {
// 		items[i] = protocol.CompletionItem{
// 			Label:         strings.TrimSpace(s),
// 			Detail:        fmt.Sprintf("%v(%v)", info.MiddlewareServerName, modelName),
// 			Documentation: &protocol.Or_CompletionItem_documentation{Value: fmt.Sprintf("%v", modelName)},
// 		}
// 	}
// 	result := &protocol.CompletionList{
// 		IsIncomplete: false,
// 		Items:        items,
// 	}

// 	return result, nil
// }

// func (s *Middleware) eventCompletionLsp(params *protocol.CompletionParams, _ *Workspace, step *config.Step) (*protocol.CompletionList, error) {
// 	if step.Scope != "client" {
// 		s.plane.Warnf("TextDocumentCompletion %T %v scope %v unsupported, set to client", step, step, step.Scope)
// 	}
// 	scope := scope.NewScopeClient(s.listenerName, s.clientName)

// 	lspName := step.Kind.(*config.StepKindLsp).Name

// 	lspKey, token, err := s.plane.Scope().LspAcquire(lspName, scope)
// 	if err != nil {
// 		s.plane.Errorf("%T %v: TextDocumentCompletion %T %v: %v", s, s.Name(), step, step, err)
// 		return nil, err
// 	}
// 	defer s.plane.Scope().LspRelease(lspKey, token)

// 	list, err := s.plane.Lsp().EventTextDocumentCompletion(lspKey, params)
// 	if err != nil {
// 		s.plane.Errorf("%T %v: TextDocumentCompletion %T %v: %v", s, s.Name(), step, step, err)
// 		return nil, err
// 	}

// 	s.plane.Debugf("%T %v: TextDocumentCompletion %T %v: %+v", s, s.Name(), step, step, list)

// 	// only return label, detail and documentation
// 	items := make([]protocol.CompletionItem, len(list.Items))
// 	for i, s := range list.Items {
// 		items[i] = protocol.CompletionItem{
// 			Label:         strings.TrimSpace(s.Label),
// 			Detail:        fmt.Sprintf("%v(%v)", info.MiddlewareServerName, lspName),
// 			Documentation: &protocol.Or_CompletionItem_documentation{Value: fmt.Sprintf("%v", lspName)},
// 		}
// 	}
// 	result := &protocol.CompletionList{
// 		IsIncomplete: list.IsIncomplete,
// 		Items:        items,
// 	}

// 	return result, nil
// }

// // TODO: make more complex workflow strategy
// func (s *Middleware) eventCompletionWorkflow(params *protocol.CompletionParams, workspace *Workspace, workflow *config.Workflow) (*protocol.CompletionList, error) {

// 	lists := make([]*protocol.CompletionList, 0)

// 	for _, step := range workflow.Steps {
// 		switch step.Kind.(type) {

// 		case *config.StepKindModel:
// 			list, err := s.eventCompletionModel(params, workspace, step)
// 			if err != nil {
// 				continue
// 			}
// 			lists = append(lists, list)

// 		case *config.StepKindLsp:
// 			list, err := s.eventCompletionLsp(params, workspace, step)
// 			if err != nil {
// 				continue
// 			}
// 			lists = append(lists, list)

// 		default:
// 			err := fmt.Errorf("unhandled step %T %v", step, step)
// 			return nil, err
// 		}
// 	}

// 	var result protocol.CompletionList

// 	for _, list := range lists {
// 		result.Items = append(result.Items, list.Items...)
// 		result.IsIncomplete = result.IsIncomplete || list.IsIncomplete
// 	}

// 	return &result, nil
// }

// func (s *Middleware) eventCompletion(params *protocol.CompletionParams, workspaces []*Workspace) (*protocol.CompletionList, error) {
// 	var wg sync.WaitGroup

// 	listCh := make(chan *protocol.CompletionList)

// 	for _, workspace := range workspaces {
// 		usages := s.plane.Usage().GetFilteredClient(
// 			usage.ConditionFilter{Filetype: nil, Path: util.Ptr(params.TextDocument.URI)},
// 			usage.EventFilter{Event: util.Ptr(config.EventTextDocumentCompletion)}, s.listenerName, s.clientName)

// 		for _, usage := range usages {
// 			s.plane.Infof("%T %v: usage %v completion", s, s.Name(), usage.Name)

// 			wg.Go(func() {
// 				list, err := s.eventCompletionWorkflow(params, workspace, usage.Workflow)

// 				if err != nil {
// 					s.plane.Warnf("%v", err)
// 					return
// 				}

// 				listCh <- list
// 			})
// 		}
// 	}

// 	go func() {
// 		wg.Wait()
// 		close(listCh)
// 	}()

// 	var result protocol.CompletionList

// 	for list := range listCh {
// 		if list != nil {
// 			result.Items = append(result.Items, list.Items...)
// 			result.IsIncomplete = result.IsIncomplete || list.IsIncomplete
// 		}
// 	}

// 	s.plane.Infof("%T %v: TextDocumentCompletion %+v", s, s.Name(), result)

// 	return &result, nil
// }

func (s *Middleware) TextDocumentCompletion(params *protocol.CompletionParams) (*protocol.CompletionList, error) {
	workspaces := s.workspaceFindAllByPrefix(params.TextDocument.URI)

	if len(workspaces) == 0 {
		s.plane.Warnf("%T %v: file %v is not bound to any workspace", s, s.Name(), params.TextDocument.URI)
	}

	// s.eventCompletion(params, workspaces)
	return &protocol.CompletionList{}, nil
}
