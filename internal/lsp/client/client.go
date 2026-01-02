package client

import (
	"bufio"
	"context"
	"mals/internal/client"
	"mals/internal/jsonrpc"
	"mals/internal/middleware"
	"mals/internal/plane"
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
		middleware: nil,
	}

	s.middleware = middleware.New(s.plane, s)

	s.scanner.Split(jsonrpc.ScannerSplit)

	return s
}

// abstract
func (s *ClientLsp) Name() string {
	return s.Client.Name()
}

func (s *ClientLsp) Serve(ctx context.Context) error {
	ch := make(chan []byte)
	defer close(ch)

	s.plane.Infof("%s: listening", s.Name())

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
				s.plane.Debugf("%s: scanned %s", s.Name(), string(s.scanner.Bytes()))
			}
		}
	}()

	defer func() {
		s.plane.Infof("%s: closed", s.Name())
	}()

	for {
		select {
		case <-scanCtx.Done():
			s.plane.Infof("%s: scanner done", s.Name())

			if err := s.writer.Flush(); err != nil {
				s.plane.Errorf("%s: flush: %v", err)
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
		s.plane.Errorf("%v", err)
		return err
	}
	if _, err := s.writer.Write(bytes); err != nil {
		s.plane.Errorf("%v", err)
		return err
	}
	if err := s.writer.Flush(); err != nil {
		s.plane.Errorf("%v", err)
		return err
	}

	s.plane.Debugf("%s: sent %s", s.Name(), string(bytes))

	return nil
}
