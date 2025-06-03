package client

import (
	"encoding/json"
	"mals-engine/internal/jsonrpc"
	"mals-engine/pkg/message"
)

func (c *Client) HandleClientRequest(bytes []byte) {
	msg, data, err := jsonrpc.DecodeRequest(bytes)

	if err != nil {
		c.LogErrorPrintf("unable to decode %s", err)
		return
	}

	c.LogInfoPrintf("handling method %s", msg.Method)

	switch msg.Method {
	case "initialize":
		c.initialize(data)
		break
	default:
		c.LogWarnPrintf("unhandled method %s", msg.Method)
		break
	}
}

func (c *Client) initialize(data []byte) {
	var initializeRequest message.InitializeRequest
	json.Unmarshal(data, &initializeRequest)
	if bytes, err := json.Marshal(initializeRequest); err == nil {
		c.LogInfoPrintf("parsed %s", string(bytes))
	}

	initializeResponse := message.InitializeResponse{
		Response: message.Response{
			Message: message.Message{
				JsonRpc: "2.0",
			},
			Id: initializeRequest.Id,
		},
		Result: message.InitializeResult{
			Capabilities: message.ServerCapabilities{
				TextDocumentSync: message.TextDocumentSyncOptions{
					OpenClose: true,
					Change:    message.FULL,
				},
				CompletionProvider: message.CompletionOptions{
					CompletionItem: message.CompletionItem{
						LabelDetailsSupport: true,
					},
				},
				HoverProvider: message.HoverOptions{},
				CodeActionProvider: message.CodeActionOptions{
					CodeActionKinds: []message.CodeActionKind{message.REFACTOR, message.QUICKFIX},
				},
				CodeLensProvider: message.CodeLensOptions{},
			},
			ServerInfo: message.ServerInfo{
				Name:    "mals-engine",
				Version: "0.0.1",
			},
		},
	}

	if bytes, err := jsonrpc.EncodeResponse(initializeResponse); err == nil {
		c.writer.Write(bytes)
		c.writer.Flush()
	} else {
		c.LogErrorPrintf("cannot encode message to send")
	}
}
