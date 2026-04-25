package handler

import (
	"fmt"
	"mals/internal/lsp/protocol"
	"mals/pkg/config"
	"mals/pkg/info"
	"strings"
)

type executionNode struct {
	Step *config.Step
	Then *executionNode
	Else *executionNode
}

type executionContext struct {
	graph *executionNode
	// resources
	// get (dict)
	// set (dict)
	// file variables
}

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

// recursively builds CFG (cuz execution graph is rather small no optimizations are done)
// returns (start, end pair execution)
// if: Then - then, Else - else
// while: Then - do, Else - <statement_after>
// simple: Then - <statement_after>
func (s *Handler) executionBuild(steps []*config.Step, i int) (*executionNode, *executionNode) {
	if i >= len(steps) {
		return nil, nil
	}

	curr := steps[i]

	switch cd := curr.Definition.(type) {
	case *config.StepIf:
		execThenStart, execThenEnd := s.executionBuild(cd.Then, 0)
		execElseStart, execElseEnd := s.executionBuild(cd.Else, 0)

		if execThenStart == nil {
			s.plane.Errorf("%T %v: TextDocumentCompletion step %v then is nil", s, s.Name(), curr)
			return nil, nil
		}

		execCurrStart := &executionNode{
			Step: curr,
			Then: execThenStart,
			Else: execElseStart,
		}

		execNextStart, execNextEnd := s.executionBuild(steps, i+1)

		if execThenEnd != nil {
			execThenEnd.Then = execNextStart
		}
		if execElseEnd != nil {
			execElseEnd.Then = execNextStart
		}

		execCurrEnd := execNextEnd
		if execCurrEnd == nil {
			execCurrEnd = execThenEnd // any branch is valid
		}
		return execCurrStart, execCurrEnd

	case *config.StepWhile:
		execDoStart, execDoEnd := s.executionBuild(cd.Do, 0)

		if execDoStart == nil {
			s.plane.Errorf("%T %v: TextDocumentCompletion step %v do is nil", s, s.Name(), curr)
			return nil, nil
		}

		execNextStart, execNextEnd := s.executionBuild(steps, i+1)

		execCurrStart := &executionNode{
			Step: curr,
			Then: execDoStart,
			Else: execNextStart,
		}

		if execDoEnd != nil {
			execDoEnd.Then = execCurrStart
		}

		execCurrEnd := execNextEnd
		if execCurrEnd == nil {
			execCurrEnd = execCurrStart // loop re-entry
		}

		return execCurrStart, execCurrEnd

	case *config.StepReturn:
		execCurrStart := &executionNode{
			Step: curr,
			Then: nil,
			Else: nil,
		}

		return execCurrStart, execCurrStart

	case *config.StepJsonDumps:
	case *config.StepJsonParse:
	case *config.StepJsonParseCompletion:
	case *config.StepLspCompetion:
	case *config.StepModelSimple:
	case *config.StepModelTemplate:
	default:
		s.plane.Errorf("%T %v: unexpected config.StepDefinition: %#v", s, s.Name(), cd)
		return nil, nil
	}

	// handle simple statements
	execNextStart, execNextEnd := s.executionBuild(steps, i+1)

	execCurrStart := &executionNode{
		Step: curr,
		Then: execNextStart,
		Else: nil,
	}

	execCurrEnd := execNextEnd
	if execCurrEnd == nil {
		execCurrEnd = execCurrStart
	}

	return execCurrStart, execCurrEnd
}

func (s *Handler) TextDocumentCompletionCustom(params *protocol.CompletionParams) (*protocol.CompletionList, error) {

	return &protocol.CompletionList{}, nil
}

func (s *Handler) TextDocumentCompletion(params *protocol.CompletionParams) (*protocol.CompletionList, error) {
	var list *protocol.CompletionList
	var err error

	s.plane.Infof("executing here")

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
