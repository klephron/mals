package handler

import (
	"fmt"
	"mals/internal/lsp/protocol"
	"mals/internal/middleware/execution"
	"mals/pkg/config"
	"mals/pkg/info"
	"strings"
)

func (s *Handler) TextDocumentCompletionDefault(params *protocol.CompletionParams) (*protocol.CompletionList, error) {
	list := protocol.CompletionList{
		IsIncomplete: false,
		Items:        make([]protocol.CompletionItem, 0),
	}

	s.resources.Range(func(key string, value *config.HandlerLspResource) bool {
		scope, err := s.getResourceScope(value.Scope)
		if err != nil {
			s.plane.Errorf("%T %v: TextDocumentCompletion %v", s, s.Name(), err)
			return true
		}

		switch vs := value.Spec.(type) {
		case *config.HandlerLspResourceSpecLsp:
			lspName, token, err := s.plane.Scope().LspAcquire(vs.Name, scope)
			if err != nil {
				s.plane.Errorf("%T %v: TextDocumentDidChange %v", s, s.Name(), err)
				return true
			}

			defer func() {
				if err := s.plane.Scope().LspRelease(lspName, token); err != nil {
					s.plane.Errorf("%T %v: TextDocumentDidChange %v", s, s.Name(), err)
				}
			}()

			lspList, err := s.plane.Lsp().TextDocumentCompletion(lspName, params)
			if err != nil {
				s.plane.Errorf("%T %v: TextDocumentCompletion %v", s, s.Name(), err)
				return true
			}

			lspItems := make([]protocol.CompletionItem, len(lspList.Items))
			for i, s := range lspList.Items {
				lspItems[i] = protocol.CompletionItem{
					Label:         strings.TrimSpace(s.Label),
					Detail:        fmt.Sprintf("%v(%v)", info.MiddlewareServerName, lspName),
					Documentation: &protocol.Or_CompletionItem_documentation{Value: fmt.Sprintf("%v", lspName)},
				}
			}

			list.Items = append(list.Items, lspItems...)

		case *config.HandlerLspResourceSpecModel:
		default:
			s.plane.Errorf("%T %v: TextDocumentCompletion unexpected spec %T", s, s.Name(), vs)
		}

		return true
	})

	return &list, nil
}

func (s *Handler) TextDocumentCompletionCustom(params *protocol.CompletionParams) (*protocol.CompletionList, error) {

	exec := execution.New(s.plane)

	exec.SetResources(s.resources)

	uri := params.TextDocument.URI
	exec.SetFileUri(uri)

	workspaces := s.workspaceFindAllByPrefix(uri)
	exec.SetWorkspaces(workspaces)

	defer exec.ResetContext()

	if err := exec.Build(s.endpoints.TextDocumentCompletion.Execution); err != nil {
		s.plane.Errorf("%T %v: %v", s, s.Name(), err)
		return nil, err
	}

	completionItemsAny, err := exec.Execute()
	if err != nil {
		s.plane.Errorf("%T %v: %v", s, s.Name(), err)
		return nil, err
	}

	completionItems, ok := completionItemsAny.([]protocol.CompletionItem)
	if !ok {
		err := fmt.Errorf("output type is not of type []protocol.CompletionItem")
		s.plane.Errorf("%T %v: %v", s, s.Name(), err)
		return nil, err
	}

	list := &protocol.CompletionList{
		IsIncomplete: false,
		Items:        completionItems,
	}

	return list, nil
}

func (s *Handler) TextDocumentCompletion(params *protocol.CompletionParams) (*protocol.CompletionList, error) {
	var list *protocol.CompletionList
	var err error

	if *s.endpoints.TextDocumentCompletion.Default {
		list, err = s.TextDocumentCompletionDefault(params)
	} else {
		list, err = s.TextDocumentCompletionCustom(params)
	}

	if err == nil {
		s.plane.Infof("%T %v: TextDocumentCompletion done", s, s.Name())
		s.plane.Debugf("%T %v: TextDocumentCompletion %+v", s, s.Name(), list)
	}

	return list, err
}
