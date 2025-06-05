package workspace

import (
	"fmt"
	"slices"
	"strings"
)

func (w *Workspace) GetCompletionList(filepath string, position Position) ([]CompletionItem, bool) {
	documentText, exists := w.Documents[filepath]
	if !exists {
		return nil, false
	}

	// TODO: change when deleging to real LSP
	words := strings.Fields(documentText)
	slices.Sort(words)
	words = slices.Compact(words)
	words = append(words, filepath)

	items := make([]CompletionItem, len(words))
	for i, s := range words {
		items[i] = CompletionItem{
			Label:         s,
			Detail:        fmt.Sprintf("%s (%d)", s, i),
			Documentation: "see dictionary",
		}
	}

	return items, true
}
