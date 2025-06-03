package state

import (
	"fmt"
	"mals-engine/pkg/message"
	"strings"
)

type Workspace struct {
	Root      string
	Documents map[string]string
}

func NewWorkspace(root string) *Workspace {
	return &Workspace{Root: root, Documents: make(map[string]string)}
}

func (w *Workspace) OpenDocument(filepath string, content string) {
	w.Documents[filepath] = content
}

func (w *Workspace) ChangeDocument(filepath string, content string) {
	w.OpenDocument(filepath, content)
}

// NOTE: maybe no need to delete because of file resolution
func (w *Workspace) CloseDocument(filepath string) {
	delete(w.Documents, filepath)
}

func (w *Workspace) GetCompletionList(filepath string, position message.Position) (*message.CompletionList, bool) {
	documentText, exists := w.Documents[filepath]
	if !exists {
		return nil, false
	}

	// TODO: change when deleging to real LSP
	words := strings.Fields(documentText)

	items := make([]message.CompletionItem, len(words))
	for i, s := range words {
		items[i] = message.CompletionItem{
			Label:         s,
			Detail:        fmt.Sprintf("%s (%d)", s, i),
			Documentation: "see dictionary",
		}
	}

	return &message.CompletionList{
		IsIncomplete: false,
		Items:        items,
	}, true
}
