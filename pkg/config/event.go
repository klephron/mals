package config

type Event string

const (
	EventInitialize             Event = "initialize"
	EventInitialized            Event = "initialized"
	EventTextDocumentDidOpen    Event = "textDocument/didOpen"
	EventTextDocumentDidChange  Event = "textDocument/didChange"
	EventTextDocumentDidClose   Event = "textDocument/didClose"
	EventTextDocumentCompletion Event = "textDocument/completion"
	EventShutdown               Event = "shutdown"
)
