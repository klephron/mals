package stdio

import "mals/internal/lsp/protocol"

func (s *LspServerStdio) Initialize(params *protocol.InitializeParams) (*protocol.InitializeResult, error) {

}

func (s *LspServerStdio) Initialized(params *protocol.InitializedParams) error {
}

func (s *LspServerStdio) TextDocumentDidOpen(params *protocol.DidOpenTextDocumentParams) error {
}

func (s *LspServerStdio) TextDocumentDidChange(params *protocol.DidChangeTextDocumentParams) error {
}

func (s *LspServerStdio) TextDocumentDidClose(params *protocol.DidCloseTextDocumentParams) error {
}

func (s *LspServerStdio) TextDocumentCompletion(params *protocol.CompletionParams) (*protocol.CompletionList, error) {
}

func (s *LspServerStdio) Shutdown() error {
}
