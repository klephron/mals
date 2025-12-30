package client

import (
	"bufio"
	"context"
	"fmt"
	"mals/internal/client"
	"mals/internal/jsonrpc"
	"mals/internal/plane"
	"strings"

	"github.com/puzpuzpuz/xsync/v4"
)

type ClientLsp struct {
	client.Client

	plane   plane.Plane
	scanner *bufio.Scanner
	writer  *bufio.Writer

	initialized bool
	workspaces  *xsync.Map[string, *Workspace]
}

func New(plane plane.Plane, scanner *bufio.Scanner, writer *bufio.Writer) *ClientLsp {
	s := &ClientLsp{
		plane:       plane,
		scanner:     scanner,
		writer:      writer,
		initialized: false,
		workspaces:  xsync.NewMap[string, *Workspace](),
	}
	s.scanner.Split(jsonrpc.ScannerSplit)
	return s
}

// abstract
func (s *ClientLsp) Name() string {
	return s.Client.Name()
}

func (s *ClientLsp) Serve(ctx context.Context) error {
	if s.scanner == nil || s.writer == nil {
		return fmt.Errorf("%s: must be binded before serve: scanner or writer is nil", s.Name())
	}

	ch := make(chan []byte)
	defer close(ch)

	s.plane.Log().Infof("%s: listening", s.Name())

	scanCtx, scanCancel := context.WithCancel(ctx)

	go func() {
		defer scanCancel()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if !s.scanner.Scan() {
					return
				}
				ch <- s.scanner.Bytes()
				s.plane.Log().Debugf("%s: scanned %s", s.Name(), string(s.scanner.Bytes()))
			}
		}
	}()

	defer func() {
		s.plane.Log().Infof("%s: closed", s.Name())
	}()

	for {
		select {
		case <-scanCtx.Done():
			s.plane.Log().Infof("%s: scanner done", s.Name())

			if err := s.writer.Flush(); err != nil {
				s.plane.Log().Errorf("%s: flush: %v", err)
				return err
			}

			return nil

		case bytes := <-ch:
			s.handle(bytes)
		}
	}
}

func (s *ClientLsp) send(msg jsonrpc.Message) error {
	bytes, err := jsonrpc.EncodeMessage(msg)
	if err != nil {
		s.plane.Log().Errorf("%v", err)
		return err
	}
	if _, err := s.writer.Write(bytes); err != nil {
		s.plane.Log().Errorf("%v", err)
		return err
	}
	if err := s.writer.Flush(); err != nil {
		s.plane.Log().Errorf("%v", err)
		return err
	}

	s.plane.Log().Debugf("%s: sent %s", s.Name(), string(bytes))

	return nil
}

func (s *ClientLsp) workspaceAdd(uri string, name string) {
	workspace := newWorkspace(uri, name)
	s.workspaces.Store(uri, workspace)

	s.plane.Log().Infof("%s: workspace %s added", s.Name(), workspace.name)
}

func (s *ClientLsp) workspaceDelete(uri string) {
	workspace, ok := s.workspaces.LoadAndDelete(uri)

	if !ok {
		s.plane.Log().Warnf("%s: workspace by uri %s is not present", s.Name(), uri)
	}
	s.plane.Log().Infof("%s: workspace %s deleted", s.Name(), workspace.name)
}

func (s *ClientLsp) workspaceFindAll(uri string) []*Workspace {
	workspaces := make([]*Workspace, 0)

	s.workspaces.Range(func(key string, value *Workspace) bool {
		if strings.HasPrefix(uri, key) {
			workspaces = append(workspaces, value)
		}
		return true
	})

	return workspaces
}
