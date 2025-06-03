package state

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
