package client

// import (
// 	"encoding/json"
// 	"mals/internal/jsonrpc"
// 	"mals/internal/workspace"
// 	"mals/pkg/lsp_message"
// 	"mals/pkg/url"
// )

// func defaultResponse(request *lsp_message.Request) lsp_message.Response {
// 	return lsp_message.Response{
// 		Message: lsp_message.Message{
// 			JsonRpc: "2.0",
// 		},
// 		Id: request.Id,
// 	}
// }

// func (s *Client) writeResponse(response any) {
// 	if bytes, err := jsonrpc.Encode(response); err == nil {
// 		s.writer.Write(bytes)
// 		s.writer.Flush()
// 	} else {
// 		s.LogErrorPrintf("cannot encode message to send")
// 	}
// }

// func (s *Client) initialize(data []byte) {
// 	var request lsp_message.InitializeRequest
// 	json.Unmarshal(data, &request)

// 	for _, workspace := range request.Params.WorkspaceFolders {
// 		path, err := url.UriToPath(workspace.URI)
// 		if err != nil {
// 			s.LogErrorPrintf("unable to get workspace %s path", workspace.URI)
// 			continue
// 		}
// 		if _, created := s.NewWorkspace(path, s.config.Workspace.DefaultModel); created {
// 			s.LogInfoPrintf("workspace %s: created", path)
// 		} else {
// 			s.LogInfoPrintf("workspace %s: recreated", path)
// 		}
// 	}

// 	response := lsp_message.InitializeResponse{
// 		Response: defaultResponse(&request.Request),
// 		Result: lsp_message.InitializeResult{
// 			Capabilities: lsp_message.ServerCapabilities{
// 				TextDocumentSync: lsp_message.TextDocumentSyncOptions{
// 					OpenClose: true,
// 					Change:    lsp_message.FULL,
// 				},
// 				CompletionProvider: lsp_message.CompletionOptions{},
// 				// HoverProvider: lsp_message.HoverOptions{},
// 				// CodeActionProvider: lsp_message.CodeActionOptions{
// 				// 	CodeActionKinds: []lsp_message.CodeActionKind{lsp_message.REFACTOR, lsp_message.QUICKFIX},
// 				// },
// 				// CodeLensProvider: lsp_message.CodeLensOptions{},
// 			},
// 			ServerInfo: lsp_message.ServerInfo{
// 				Name:    "mals-engine",
// 				Version: "0.0.1",
// 			},
// 		},
// 	}

// 	s.writeResponse(response)
// }

// func (s *Client) textDocumentDidOpen(data []byte) {
// 	var notification lsp_message.DidOpenTextDocumentNotification
// 	json.Unmarshal(data, &notification)

// 	fileUri := notification.Params.TextDocument.Uri
// 	path, err := url.UriToPath(fileUri)
// 	if err != nil {
// 		s.LogErrorPrintf("unable to resolve file uri: %s", err)
// 		return
// 	}

// 	if workspace, found := s.FindWorkspace(path); found {
// 		workspace.OpenDocument(path, notification.Params.TextDocument.Text)
// 		s.LogInfoPrintf("workspace %s: document %s: opened", workspace.Root, path)
// 	}
// }

// func (s *Client) textDocumentDidChange(data []byte) {
// 	var notification lsp_message.DidChangeTextDocumentNotification
// 	json.Unmarshal(data, &notification)

// 	fileUri := notification.Params.TextDocument.Uri
// 	path, err := url.UriToPath(fileUri)
// 	if err != nil {
// 		s.LogErrorPrintf("unable to resolve file uri: %s", err)
// 		return
// 	}

// 	if workspace, found := s.FindWorkspace(path); found {
// 		for _, change := range notification.Params.ContentChanges {
// 			// NOTE: assumes change == FULL
// 			workspace.ChangeDocument(path, change.Text)
// 			s.LogInfoPrintf("workspace %s: document %s: changed", workspace.Root, path)
// 		}
// 	}
// }

// func (s *Client) textDocumentDidClose(data []byte) {
// 	var notification lsp_message.DidCloseTextDocumentNotification
// 	json.Unmarshal(data, &notification)

// 	fileUri := notification.Params.TextDocument.Uri
// 	path, err := url.UriToPath(fileUri)
// 	if err != nil {
// 		s.LogErrorPrintf("unable to resolve file uri: %s", err)
// 		return
// 	}

// 	if workspace, found := s.FindWorkspace(path); found {
// 		workspace.CloseDocument(path)
// 		s.LogInfoPrintf("workspace %s: document %s: closed", workspace.Root, path)
// 	}
// }

// func (s *Client) textDocumentCompletion(data []byte) {
// 	var request lsp_message.CompletionRequest
// 	json.Unmarshal(data, &request)

// 	fileUri := request.Params.TextDocument.Uri
// 	path, err := url.UriToPath(fileUri)
// 	if err != nil {
// 		s.LogErrorPrintf("unable to resolve file uri: %s", err)
// 		return
// 	}

// 	w, found := s.FindWorkspace(path)
// 	if !found {
// 		s.LogErrorPrintf("workspace %s not found", err)
// 		return
// 	}

// 	s.LogInfoPrintf("workspace %s: document %s: completion", w.Root, path)

// 	// in separate goroutine because it is a long process
// 	go func() {
// 		items, err := w.GenerateCompletionList(path, workspace.Position{
// 			Line: request.Params.Position.Line,
// 			Char: request.Params.Position.Character,
// 		})

// 		if err != nil {
// 			s.LogErrorPrintf("workspace %s: %s: ", w.Root, err)
// 			return
// 		}

// 		list_items := make([]lsp_message.CompletionItem, len(items))
// 		for i, s := range items {
// 			list_items[i] = lsp_message.CompletionItem{
// 				Label:         s.Label,
// 				Detail:        s.Detail,
// 				Documentation: s.Documentation,
// 			}
// 		}

// 		response := lsp_message.CompletionResponse{
// 			Response: defaultResponse(&request.Request),
// 			Result: lsp_message.CompletionList{
// 				IsIncomplete: true,
// 				Items:        list_items,
// 			},
// 		}
// 		s.writeResponse(response)
// 	}()
// }

func (s *ClientGeneric) LspHandle(bytes []byte) {
	// msg, data, err := jsonrpc.DecodeNotification(bytes)

	// if err != nil {
	// 	s.LogErrorPrintf("unable to decode %s", err)
	// 	return
	// }

	// switch msg.Method {
	// case "initialize":
	// 	s.initialize(data)
	// case "initialized":
	// case "textDocument/didOpen":
	// 	s.textDocumentDidOpen(data)
	// case "textDocument/didChange":
	// 	s.textDocumentDidChange(data)
	// case "textDocument/didClose":
	// 	s.textDocumentDidClose(data)
	// case "textDocument/completion":
	// 	s.textDocumentCompletion(data)
	// default:
	// 	s.LogWarnPrintf("unhandled method %s", msg.Method)
	// }
}
