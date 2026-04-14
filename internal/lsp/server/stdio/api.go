package stdio

import (
	"errors"
	"mals/internal/jsonrpc"
	"mals/internal/lsp/protocol"
	"mals/internal/util"
)

func (s *LspServerStdio) Capabilities() (*protocol.ServerCapabilities, error) {
	s.rw.RLock()
	defer s.rw.RUnlock()

	if !s.running {
		return nil, s.errorNotRunning()
	}

	return s.capabilities, nil
}

func (s *LspServerStdio) Info() (*protocol.ServerInfo, error) {
	s.rw.RLock()
	defer s.rw.RUnlock()

	if !s.running {
		return nil, s.errorNotRunning()
	}

	return s.info, nil
}

func (s *LspServerStdio) Initialize(params *protocol.InitializeParams) (*protocol.InitializeResult, error) {
	paramsRaw, err := util.JsonMarshal(params)
	if err != nil {
		s.plane.Errorf("%T %v: %v", s, s.Name(), err)
		return nil, err
	}

	request := &jsonrpc.Request{
		Method: "initialize",
		Params: paramsRaw,
	}

	res, err := s.sendRequest(request)
	if err != nil {
		return nil, err
	}

	resp := <-res

	if resp == nil {
		return nil, s.plane.Errorf("%T %v: Initialize response nil", s, s.Name())
	}

	if resp.Error != nil {
		return nil, errors.New(resp.Error.Message)
	}

	result, err := util.JsonUnmarshal[*protocol.InitializeResult](resp.Result)
	if err != nil {
		s.plane.Warnf("%T %v: %v", s, s.Name(), err)
	}

	s.rw.Lock()
	defer s.rw.Unlock()

	if s.running {
		s.capabilities = &result.Capabilities
		s.info = result.ServerInfo
	}

	return result, nil
}

func (s *LspServerStdio) Initialized(params *protocol.InitializedParams) error {
	paramsRaw, err := util.JsonMarshal(params)
	if err != nil {
		s.plane.Errorf("%T %v: %v", s, s.Name(), err)
		return err
	}

	notification := &jsonrpc.Notification{
		Method: "initialized",
		Params: paramsRaw,
	}

	err = s.sendNotification(notification)
	if err != nil {
		return err
	}

	return nil
}

func (s *LspServerStdio) TextDocumentDidOpen(params *protocol.DidOpenTextDocumentParams) error {
	paramsRaw, err := util.JsonMarshal(params)
	if err != nil {
		s.plane.Errorf("%T %v: %v", s, s.Name(), err)
		return err
	}

	notification := &jsonrpc.Notification{
		Method: "textDocument/didOpen",
		Params: paramsRaw,
	}

	err = s.sendNotification(notification)
	if err != nil {
		return err
	}

	return nil
}

func (s *LspServerStdio) TextDocumentDidChange(params *protocol.DidChangeTextDocumentParams) error {
	paramsRaw, err := util.JsonMarshal(params)
	if err != nil {
		s.plane.Errorf("%T %v: %v", s, s.Name(), err)
		return err
	}

	notification := &jsonrpc.Notification{
		Method: "textDocument/didChange",
		Params: paramsRaw,
	}

	err = s.sendNotification(notification)
	if err != nil {
		return err
	}

	return nil
}

func (s *LspServerStdio) TextDocumentDidClose(params *protocol.DidCloseTextDocumentParams) error {
	paramsRaw, err := util.JsonMarshal(params)
	if err != nil {
		s.plane.Errorf("%T %v: %v", s, s.Name(), err)
		return err
	}

	notification := &jsonrpc.Notification{
		Method: "textDocument/didClose",
		Params: paramsRaw,
	}

	err = s.sendNotification(notification)
	if err != nil {
		return err
	}

	return nil
}

func (s *LspServerStdio) TextDocumentCompletion(params *protocol.CompletionParams) (*protocol.CompletionList, error) {
	paramsRaw, err := util.JsonMarshal(params)
	if err != nil {
		s.plane.Errorf("%T %v: %v", s, s.Name(), err)
		return nil, err
	}

	request := &jsonrpc.Request{
		Method: "textDocument/completion",
		Params: paramsRaw,
	}

	res, err := s.sendRequest(request)
	if err != nil {
		return nil, err
	}

	resp := <-res

	if resp == nil {
		return nil, s.plane.Errorf("%T %v: TextDocumentCompletion response nil", s, s.Name())
	}

	if resp.Error != nil {
		return nil, errors.New(resp.Error.Message)
	}

	result, err := util.JsonUnmarshal[*protocol.CompletionList](resp.Result)
	if err != nil {
		s.plane.Warnf("%T %v: %v", s, s.Name(), err)
	}

	return result, nil
}

func (s *LspServerStdio) Shutdown() error {
	notification := &jsonrpc.Notification{
		Method: "shutdown",
		Params: nil,
	}

	err := s.sendNotification(notification)
	if err != nil {
		return err
	}

	return nil
}
