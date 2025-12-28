package lsp

import (
	"fmt"
	"mals/internal/lsp/protocol"
	"mals/internal/scope"
	"mals/pkg/config"

	"github.com/invopop/jsonschema"
)

const (
	completionTemplate = `
<|begin_of_text|><|start_header_id|>system<|end_header_id|>

You are an intelligent autocompletion engine. Your task is to suggest relevant text completions based on document context.

Rules:
- Analyze the document context and writing style
- Generate 3-5 contextually appropriate completions
- Match the existing tone and style
- Provide completions of varying lengths (words to phrases)
- Return ONLY a valid JSON array of strings
- No explanations, no additional text

<|eot_id|><|start_header_id|>user<|end_header_id|>

Document text: The weather today is quite nice. I think I'll go for a walk in the park and maybe stop by the
Current context: "I think I'll go for a walk in the park and maybe stop by the"

<|eot_id|><|start_header_id|>assistant<|end_header_id|>

["coffee shop", "library", "store to buy groceries", "lake to feed the ducks"]

<|eot_id|><|start_header_id|>user<|end_header_id|>

Document text: She opened her laptop and started typing her report. The first chapter discussed the importance of renewable energy sources. Solar panels and wind turbines are becoming more
Current context: "Solar panels and wind turbines are becoming more"

<|eot_id|><|start_header_id|>assistant<|end_header_id|>

["efficient", "affordable each year", "popular worldwide", "cost-effective solutions", "widely adopted"]

<|eot_id|><|start_header_id|>user<|end_header_id|>

Document text:
%s

Current context:
%s

<|eot_id|><|start_header_id|>assistant<|end_header_id|>
`
)

func completionPrompt(documentContext string, currentContext string) string {
	prompt := fmt.Sprintf(completionTemplate, documentContext, currentContext)
	return prompt
}

func completionSchema[T any]() any {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}

func (s *ClientLsp) completionWorkflow(params *protocol.CompletionParams, document *Document, workflow *config.Workflow) (*protocol.CompletionList, error) {

	for _, step := range workflow.Steps {
		switch step := step.(type) {
		case *config.StepModel:
			// TODO: make it work not only on global scope
			if step.Scope != "global" {
				s.plane.Log().Warnf("unsupported step %T %v scope %v", step, step, step.Scope)
			}

			// acquire model with scope
			modelName, token, err := s.plane.Scope().ModelAcquire(step.Model, scope.NewScopeGlobal())
			if err != nil {
				s.plane.Log().Warnf("while step %T %v: %v", step, step, err)
				return nil, err
			}

			defer s.plane.Scope().ModelRelease(modelName, token)

		default:
			s.plane.Log().Warnf("unhandled step %T %v", step, step)
		}
	}

	list := protocol.CompletionList{
		IsIncomplete: true,
		Items: []protocol.CompletionItem{
			{
				Label:         "lsp-test",
				Detail:        "lsp-details",
				Documentation: &protocol.Or_CompletionItem_documentation{Value: "lsp-documentation"},
			},
		},
	}

	return &list, nil
}
