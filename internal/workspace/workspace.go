package workspace

type Workspace struct {
	Root      string
	Documents map[string]string
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

func NewWorkspace(root string) *Workspace {
	return &Workspace{Root: root, Documents: make(map[string]string)}
}
