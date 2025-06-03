package client

import (
	"encoding/json"
	"mals-engine/internal/jsonrpc"
	"mals-engine/pkg/message"
	"mals-engine/pkg/uri"
)

func defaultResponse(request *message.Request) message.Response {
	return message.Response{
		Message: message.Message{
			JsonRpc: "2.0",
		},
		Id: request.Id,
	}
}

func (c *Client) HandleClientRequest(bytes []byte) {
	msg, data, err := jsonrpc.DecodeNotification(bytes)

	if err != nil {
		c.LogErrorPrintf("unable to decode %s", err)
		return
	}

	switch msg.Method {
	case "initialize":
		c.initialize(data)
		break
	case "initialized":
		break
	case "textDocument/didOpen":
		c.textDocumentDidOpen(data)
		break
	case "textDocument/didChange":
		c.textDocumentDidChange(data)
		break
	case "textDocument/didClose":
		c.textDocumentDidClose(data)
		break
	case "textDocument/completion":
		c.textDocumentCompletion(data)
		break
	default:
		c.LogWarnPrintf("unhandled method %s", msg.Method)
		break
	}
}

func (c *Client) writeResponse(response any) {
	if bytes, err := jsonrpc.EncodeResponse(response); err == nil {
		c.writer.Write(bytes)
		c.writer.Flush()
	} else {
		c.LogErrorPrintf("cannot encode message to send")
	}
}

func (c *Client) initialize(data []byte) {
	var request message.InitializeRequest
	json.Unmarshal(data, &request)

	for _, workspace := range request.Params.WorkspaceFolders {
		path, err := uri.UriToPath(workspace.URI)
		if err != nil {
			c.LogErrorPrintf("unable to get workspace %s path", workspace.URI)
			continue
		}
		if _, created := c.state.NewWorkspace(path); created {
			c.LogInfoPrintf("workspace %s: created", path)
		} else {
			c.LogInfoPrintf("workspace %s: recreated", path)
		}
	}

	response := message.InitializeResponse{
		Response: defaultResponse(&request.Request),
		Result: message.InitializeResult{
			Capabilities: message.ServerCapabilities{
				TextDocumentSync: message.TextDocumentSyncOptions{
					OpenClose: true,
					Change:    message.FULL,
				},
				CompletionProvider: message.CompletionOptions{},
				// HoverProvider: message.HoverOptions{},
				// CodeActionProvider: message.CodeActionOptions{
				// 	CodeActionKinds: []message.CodeActionKind{message.REFACTOR, message.QUICKFIX},
				// },
				// CodeLensProvider: message.CodeLensOptions{},
			},
			ServerInfo: message.ServerInfo{
				Name:    "mals-engine",
				Version: "0.0.1",
			},
		},
	}

	c.writeResponse(response)
}

func (c *Client) textDocumentDidOpen(data []byte) {
	var notification message.DidOpenTextDocumentNotification
	json.Unmarshal(data, &notification)

	fileUri := notification.Params.TextDocument.Uri
	path, err := uri.UriToPath(fileUri)
	if err != nil {
		c.LogErrorPrintf("unable to resolve file uri: %s", err)
		return
	}

	if workspace, found := c.state.FindWorkspace(path); found {
		workspace.OpenDocument(path, notification.Params.TextDocument.Text)
		c.LogInfoPrintf("workspace %s: document %s: opened", workspace.Root, path)
	}
}

func (c *Client) textDocumentDidChange(data []byte) {
	var notification message.DidChangeTextDocumentNotification
	json.Unmarshal(data, &notification)

	fileUri := notification.Params.TextDocument.Uri
	path, err := uri.UriToPath(fileUri)
	if err != nil {
		c.LogErrorPrintf("unable to resolve file uri: %s", err)
		return
	}

	if workspace, found := c.state.FindWorkspace(path); found {
		for _, change := range notification.Params.ContentChanges {
			// NOTE: assumes change == FULL
			workspace.ChangeDocument(path, change.Text)
			c.LogInfoPrintf("workspace %s: document %s: changed", workspace.Root, path)
		}
	}
}

func (c *Client) textDocumentDidClose(data []byte) {
	var notification message.DidCloseTextDocumentNotification
	json.Unmarshal(data, &notification)

	fileUri := notification.Params.TextDocument.Uri
	path, err := uri.UriToPath(fileUri)
	if err != nil {
		c.LogErrorPrintf("unable to resolve file uri: %s", err)
		return
	}

	if workspace, found := c.state.FindWorkspace(path); found {
		workspace.CloseDocument(path)
		c.LogInfoPrintf("workspace %s: document %s: closed", workspace.Root, path)
	}
}

func (c *Client) textDocumentCompletion(data []byte) {
	var request message.CompletionRequest
	json.Unmarshal(data, &request)

	fileUri := request.Params.TextDocument.Uri
	path, err := uri.UriToPath(fileUri)
	if err != nil {
		c.LogErrorPrintf("unable to resolve file uri: %s", err)
		return
	}

	if workspace, found := c.state.FindWorkspace(path); found {
		list, ok := workspace.GetCompletionList(path, request.Params.Position)
		if !ok {
			c.LogErrorPrintf("workspace %s: document %s: not found for completion", workspace.Root, path)
		} else {
			c.LogInfoPrintf("workspace %s: document %s: completion", workspace.Root, path)

			response := message.CompletionResponse{
				Response: defaultResponse(&request.Request),
				Result:   *list,
			}
			c.writeResponse(response)
		}
	}
}
