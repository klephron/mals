package client

import (
	"bufio"
	"context"
	"mals/internal/client"
	"mals/internal/jsonrpc"
	"mals/internal/middleware"
	"mals/internal/plane"
	"mals/internal/scope"
)

type ClientLsp struct {
	client.Client

	plane   plane.Plane
	scanner *bufio.Scanner
	writer  *bufio.Writer

	middleware *middleware.Middleware
}

func New(plane plane.Plane, scanner *bufio.Scanner, writer *bufio.Writer) *ClientLsp {
	s := &ClientLsp{
		plane:      plane,
		scanner:    scanner,
		writer:     writer,
		middleware: middleware.New(plane),
	}
	s.Client = s

	s.scanner.Split(jsonrpc.ScannerSplit)

	return s
}

func (s *ClientLsp) Name() string {
	return s.Client.Name()
}

func (s *ClientLsp) Serve(ctx context.Context) error {
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

	errs := s.plane.ScopeClose(scope.NewScopeClient(s.Name()))
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

func (s *ClientLsp) send(msg jsonrpc.Message) error {
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
