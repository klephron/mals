package workspace

import (
	"mals-engine/internal/model"
	"sync"
)

type Workspace struct {
	Root        string
	Documents   map[string]string
	model       model.ModelService
	modelRespCh chan model.ModelResponse
	closeOnce   sync.Once
}

type Position struct {
	Line uint32
	Char uint32
}

type CompletionItem struct {
	Label         string
	Detail        string
	Documentation string
}

func NewWorkspace(root string, m model.ModelService) *Workspace {
	return &Workspace{
		Root:        root,
		Documents:   make(map[string]string),
		model:       m,
		modelRespCh: make(chan model.ModelResponse),
	}
}

func (w *Workspace) Close() {
	w.closeOnce.Do(func() {
		close(w.modelRespCh)
	})
}
