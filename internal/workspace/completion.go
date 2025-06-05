package workspace

import (
	"encoding/json"
	"errors"
	"fmt"
	"mals-engine/internal/model"

	"github.com/invopop/jsonschema"
)

func GenerateSchema[T any]() any {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}

// long performing task
func (w *Workspace) GenerateCompletionList(filepath string, position Position) ([]CompletionItem, error) {
	documentText, exists := w.Documents[filepath]
	if !exists {
		return nil, errors.New(fmt.Sprintf("document %s doesn't exist", filepath))
	}

	// TODO: change when deleging to real LSP
	// words := strings.Fields(documentText)
	// slices.Sort(words)
	// words = slices.Compact(words)
	// words = append(words, filepath)

	request := model.NewModelRequest(
		GetCompletionPrompt(documentText),
		GenerateSchema[[]string](),
		"completion_items",
		"Generated completion items",
	)

	task := model.NewModelTask(request, w.modelRespCh)
	w.model.SubmitTask(task)

	resp := <-w.modelRespCh

	if resp.Error != nil {
		return nil, resp.Error
	}

	var respItemArray []string
	if err := json.Unmarshal([]byte(resp.Text), &respItemArray); err != nil {
		return nil, err
	}

	items := make([]CompletionItem, len(respItemArray))
	for i, s := range respItemArray {
		items[i] = CompletionItem{
			Label:         s,
			Detail:        fmt.Sprintf("%s (%d)", s, i),
			Documentation: "see dictionary",
		}
	}

	return items, nil
}
