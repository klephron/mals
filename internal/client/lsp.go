package client

import (
	"encoding/json"
	"mals/internal/jsonrpc"
	"mals/internal/workspace"
	"mals/pkg/lsp_message"
	"mals/pkg/url"
)

func defaultResponse(request *lsp_message.Request) lsp_message.Response {
	return lsp_message.Response{
		Message: lsp_message.Message{
			JsonRpc: "2.0",
		},
		Id: request.Id,
	}
}

func (c *Client) writeResponse(response any) {
	if bytes, err := jsonrpc.Encode(response); err == nil {
		c.writer.Write(bytes)
		c.writer.Flush()
	} else {
		c.LogErrorPrintf("cannot encode message to send")
	}
}

func (c *Client) initialize(data []byte) {
	var request lsp_message.InitializeRequest
	json.Unmarshal(data, &request)

	for _, workspace := range request.Params.WorkspaceFolders {
		path, err := url.UriToPath(workspace.URI)
		if err != nil {
			c.LogErrorPrintf("unable to get workspace %s path", workspace.URI)
			continue
		}
		if _, created := c.NewWorkspace(path, c.config.Workspace.DefaultModel); created {
			c.LogInfoPrintf("workspace %s: created", path)
		} else {
			c.LogInfoPrintf("workspace %s: recreated", path)
		}
	}

	response := lsp_message.InitializeResponse{
		Response: defaultResponse(&request.Request),
		Result: lsp_message.InitializeResult{
			Capabilities: lsp_message.ServerCapabilities{
				TextDocumentSync: lsp_message.TextDocumentSyncOptions{
					OpenClose: true,
					Change:    lsp_message.FULL,
				},
				CompletionProvider: lsp_message.CompletionOptions{},
				// HoverProvider: lsp_message.HoverOptions{},
				// CodeActionProvider: lsp_message.CodeActionOptions{
				// 	CodeActionKinds: []lsp_message.CodeActionKind{lsp_message.REFACTOR, lsp_message.QUICKFIX},
				// },
				// CodeLensProvider: lsp_message.CodeLensOptions{},
			},
			ServerInfo: lsp_message.ServerInfo{
				Name:    "mals-engine",
				Version: "0.0.1",
			},
		},
	}

	c.writeResponse(response)
}

func (c *Client) textDocumentDidOpen(data []byte) {
	var notification lsp_message.DidOpenTextDocumentNotification
	json.Unmarshal(data, &notification)

	fileUri := notification.Params.TextDocument.Uri
	path, err := url.UriToPath(fileUri)
	if err != nil {
		c.LogErrorPrintf("unable to resolve file uri: %s", err)
		return
	}

	if workspace, found := c.FindWorkspace(path); found {
		workspace.OpenDocument(path, notification.Params.TextDocument.Text)
		c.LogInfoPrintf("workspace %s: document %s: opened", workspace.Root, path)
	}
}

func (c *Client) textDocumentDidChange(data []byte) {
	var notification lsp_message.DidChangeTextDocumentNotification
	json.Unmarshal(data, &notification)

	fileUri := notification.Params.TextDocument.Uri
	path, err := url.UriToPath(fileUri)
	if err != nil {
		c.LogErrorPrintf("unable to resolve file uri: %s", err)
		return
	}

	if workspace, found := c.FindWorkspace(path); found {
		for _, change := range notification.Params.ContentChanges {
			// NOTE: assumes change == FULL
			workspace.ChangeDocument(path, change.Text)
			c.LogInfoPrintf("workspace %s: document %s: changed", workspace.Root, path)
		}
	}
}

func (c *Client) textDocumentDidClose(data []byte) {
	var notification lsp_message.DidCloseTextDocumentNotification
	json.Unmarshal(data, &notification)

	fileUri := notification.Params.TextDocument.Uri
	path, err := url.UriToPath(fileUri)
	if err != nil {
		c.LogErrorPrintf("unable to resolve file uri: %s", err)
		return
	}

	if workspace, found := c.FindWorkspace(path); found {
		workspace.CloseDocument(path)
		c.LogInfoPrintf("workspace %s: document %s: closed", workspace.Root, path)
	}
}

func (c *Client) textDocumentCompletion(data []byte) {
	var request lsp_message.CompletionRequest
	json.Unmarshal(data, &request)

	fileUri := request.Params.TextDocument.Uri
	path, err := url.UriToPath(fileUri)
	if err != nil {
		c.LogErrorPrintf("unable to resolve file uri: %s", err)
		return
	}

	w, found := c.FindWorkspace(path)
	if !found {
		c.LogErrorPrintf("workspace %s not found", err)
		return
	}

	c.LogInfoPrintf("workspace %s: document %s: completion", w.Root, path)

	// in separate goroutine because it is a long process
	go func() {
		items, err := w.GenerateCompletionList(path, workspace.Position{
			Line: request.Params.Position.Line,
			Char: request.Params.Position.Character,
		})

		if err != nil {
			c.LogErrorPrintf("workspace %s: %s: ", w.Root, err)
			return
		}

		list_items := make([]lsp_message.CompletionItem, len(items))
		for i, s := range items {
			list_items[i] = lsp_message.CompletionItem{
				Label:         s.Label,
				Detail:        s.Detail,
				Documentation: s.Documentation,
			}
		}

		response := lsp_message.CompletionResponse{
			Response: defaultResponse(&request.Request),
			Result: lsp_message.CompletionList{
				IsIncomplete: true,
				Items:        list_items,
			},
		}
		c.writeResponse(response)
	}()
}

func (c *Client) HandleLspRequest(bytes []byte) {
	msg, data, err := jsonrpc.DecodeNotification(bytes)

	if err != nil {
		c.LogErrorPrintf("unable to decode %s", err)
		return
	}

	switch msg.Method {
	case "initialize":
		c.initialize(data)
	case "initialized":
	case "textDocument/didOpen":
		c.textDocumentDidOpen(data)
	case "textDocument/didChange":
		c.textDocumentDidChange(data)
	case "textDocument/didClose":
		c.textDocumentDidClose(data)
	case "textDocument/completion":
		c.textDocumentCompletion(data)
	default:
		c.LogWarnPrintf("unhandled method %s", msg.Method)
	}
}
