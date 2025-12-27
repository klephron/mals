package lsp

// import (
// 	"encoding/json"
// 	"mals/internal/jsonrpc"
// 	"mals/internal/workspace"
// 	"mals/pkg/lsp_message"
// 	"mals/pkg/url"
// )

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
