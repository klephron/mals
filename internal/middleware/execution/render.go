package execution

import (
	"bytes"
	"fmt"
	"mals/internal/util"
	"strconv"
	"strings"
	"text/template"
	"text/template/parse"
)

func rewriteNode(node parse.Node) {
	switch n := node.(type) {
	case *parse.ListNode:
		if n == nil {
			return
		}
		for _, child := range n.Nodes {
			rewriteNode(child)
		}
	case *parse.ActionNode:
		rewritePipe(n.Pipe)
	case *parse.IfNode:
		rewritePipe(n.Pipe)
		rewriteNode(n.List)
		rewriteNode(n.ElseList)
	case *parse.RangeNode:
		rewritePipe(n.Pipe)
		rewriteNode(n.List)
		rewriteNode(n.ElseList)
	case *parse.WithNode:
		rewritePipe(n.Pipe)
		rewriteNode(n.List)
		rewriteNode(n.ElseList)
	case *parse.TemplateNode:
		rewritePipe(n.Pipe)
	}
}

func rewritePipe(pipe *parse.PipeNode) {
	if pipe == nil {
		return
	}
	for _, cmd := range pipe.Cmds {
		rewriteCmd(cmd)
	}
}

func rewriteCmd(cmd *parse.CommandNode) {
	if len(cmd.Args) == 0 {
		return
	}

	// {{.foo.bar}} -> get "foo" "bar"
	if field, ok := cmd.Args[0].(*parse.FieldNode); ok {
		pos := field.Position()
		newArgs := []parse.Node{
			&parse.IdentifierNode{NodeType: parse.NodeIdentifier, Pos: pos, Ident: "get"},
		}
		for _, seg := range field.Ident {
			newArgs = append(newArgs, &parse.StringNode{
				NodeType: parse.NodeString,
				Pos:      pos,
				Quoted:   `"` + seg + `"`,
				Text:     seg,
			})
		}
		cmd.Args = newArgs
		return
	}

	// {{index .foo.bar 0}} -> get "foo" "bar" "0"
	ident, ok := cmd.Args[0].(*parse.IdentifierNode)
	if !ok || ident.Ident != "index" {
		return
	}
	if len(cmd.Args) < 3 {
		return
	}
	field, ok := cmd.Args[1].(*parse.FieldNode)
	if !ok {
		return
	}
	idxNode, ok := cmd.Args[2].(*parse.NumberNode)
	if !ok {
		return
	}

	pos := field.Position()
	newArgs := []parse.Node{
		&parse.IdentifierNode{NodeType: parse.NodeIdentifier, Pos: pos, Ident: "get"},
	}
	for _, seg := range field.Ident {
		newArgs = append(newArgs, &parse.StringNode{
			NodeType: parse.NodeString,
			Pos:      pos,
			Quoted:   `"` + seg + `"`,
			Text:     seg,
		})
	}
	newArgs = append(newArgs, &parse.StringNode{
		NodeType: parse.NodeString,
		Pos:      pos,
		Quoted:   `"` + idxNode.Text + `"`,
		Text:     idxNode.Text,
	})
	cmd.Args = newArgs
}

func (s *ExecutionEnvironment) renderString(tmpl string, memory map[string]any) (*string, error) {
	funcs := template.FuncMap{
		"get": func(segments ...string) (any, error) {
			return s.get(memory, segments...)
		},
		"index": func(a any, i any) (any, error) { return nil, nil },
	}

	trees, err := parse.Parse("root", tmpl, "", "", funcs)
	if err != nil {
		return nil, err
	}

	t := template.New("root").Funcs(funcs)
	for name, tree := range trees {
		rewriteNode(tree.Root)
		if name == "root" {
			t.Tree = tree
		} else {
			nt, _ := t.New(name).Parse("") // initialize the slot
			nt.Tree = tree
		}
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, nil); err != nil {
		return nil, err
	}
	return util.Ptr(buf.String()), nil
}

func (s *ExecutionEnvironment) renderBool(tmpl string, data map[string]any) (*bool, error) {
	str, err := s.renderString(tmpl, data)
	if err != nil {
		return nil, err
	}
	if str == nil {
		return nil, nil
	}
	switch *str {
	case "true":
		return util.Ptr(true), nil
	case "false":
		return util.Ptr(false), nil
	default:
		return nil, fmt.Errorf("cannot convert '%v' to bool", *str)
	}
}

func (s *ExecutionEnvironment) renderInt(tmpl string, data map[string]any) (*int, error) {
	str, err := s.renderString(tmpl, data)
	if err != nil {
		return nil, err
	}
	if str == nil {
		return nil, nil
	}
	i, err := strconv.Atoi(strings.TrimSpace(*str))
	if err != nil {
		return nil, fmt.Errorf("cannot convert '%v' to int", *str)
	}
	return util.Ptr(i), nil
}
