package client

import (
	"bufio"
	"context"
	"fmt"
	"mals/internal/jsonrpc"
	"mals/internal/middleware"
	"mals/internal/plane"
	"mals/internal/scope"
)

type LspClient struct {
	listenerName string
	clientName   string
	plane        plane.Plane
	scanner      *bufio.Scanner
	writer       *bufio.Writer
	middleware   *middleware.Middleware
}

func New(listenerName string, clientName string, plane plane.Plane, scanner *bufio.Scanner, writer *bufio.Writer) *LspClient {
	s := &LspClient{
		listenerName: listenerName,
		clientName:   clientName,
		plane:        plane,
		scanner:      scanner,
		writer:       writer,
		middleware:   middleware.New(plane),
	}

	s.scanner.Split(jsonrpc.ScannerSplit)

	return s
}

func (s *LspClient) Name() string {
	return fmt.Sprintf("%v:%v", s.listenerName, s.clientName)
}

func (s *LspClient) Serve(ctx context.Context) error {
	s.plane.Infof("%v: listening", s.Name())

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

				bytes := s.scanner.Bytes()
				s.plane.Debugf("%v: scanned %v", s.Name(), string(bytes))

				s.handle(bytes)
			}
		}
	}()

	<-scanCtx.Done()

	errs := s.plane.ScopeClose(scope.NewScopeClient(s.listenerName, s.clientName))
	for _, err := range errs {
		s.plane.Warnf("%v: %v", s.Name(), err)
	}

	s.plane.Infof("%v: done", s.Name())

	if err := s.writer.Flush(); err != nil {
		s.plane.Errorf("%v: flush: %v", s.Name(), err)
		return err
	}

	s.plane.Infof("%s: closed", s.Name())

	return nil
}

func (s *LspClient) send(msg jsonrpc.Message) error {
	bytes, err := jsonrpc.EncodeMessage(msg)
	if err != nil {
		s.plane.Errorf("%T %v: %v", s, s.Name(), err)
		return err
	}
	if _, err := s.writer.Write(bytes); err != nil {
		s.plane.Errorf("%T %v: %v", s, s.Name(), err)
		return err
	}
	if err := s.writer.Flush(); err != nil {
		s.plane.Errorf("%T %v: %v", s, s.Name(), err)
		return err
	}

	s.plane.Debugf("%s: sent %s", s.Name(), string(bytes))

	return nil
}
