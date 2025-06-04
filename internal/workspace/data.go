package workspace

type Position struct {
	Line uint32
	Char uint32
}

type CompletionItem struct {
	Label string
	Detail string
	Documentation string
}
