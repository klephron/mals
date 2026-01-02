package stdio

import (
	"errors"
	"mals/internal/jsonrpc"
	"mals/internal/lsp/protocol"
	"mals/internal/util"
	"mals/pkg/config"
)

func (s *LspServerStdio) Initialize(params *protocol.InitializeParams) (*protocol.InitializeResult, error) {
	paramsRaw, err := util.JsonMarshal(params)
	if err != nil {
		s.plane.Errorf("%T %v: %v", s, s.Name(), err)
		return nil, err
	}

	request := &jsonrpc.Request{
		Method: string(config.EventInitialize),
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
		s.plane.Errorf("%T %v: %v", s, s.Name(), err)
		return nil, err
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
		Method: string(config.EventInitialized),
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
		Method: string(config.EventTextDocumentDidOpen),
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
		Method: string(config.EventTextDocumentDidChange),
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
		Method: string(config.EventTextDocumentDidClose),
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
		Method: string(config.EventTextDocumentCompletion),
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
		s.plane.Errorf("%T %v: %v", s, s.Name(), err)
		return nil, err
	}

	return result, nil
}

func (s *LspServerStdio) Shutdown() error {
	notification := &jsonrpc.Notification{
		Method: string(config.EventShutdown),
		Params: nil,
	}

	err := s.sendNotification(notification)
	if err != nil {
		return err
	}

	return nil
}
