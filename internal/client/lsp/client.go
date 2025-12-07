package lsp

import (
	"bufio"
	"context"
	"fmt"
	"mals/internal/client"
	"mals/internal/jsonrpc"
	"mals/internal/plane"
)

type ClientLsp struct {
	client.Client
	plane   plane.Plane
	scanner *bufio.Scanner
	writer  *bufio.Writer
}

func New(plane plane.Plane) *ClientLsp {
	s := &ClientLsp{
		plane:   plane,
		scanner: nil,
		writer:  nil,
	}
	return s
}

func (s *ClientLsp) Bind(scanner *bufio.Scanner, writer *bufio.Writer) error {
	if err := s.Unbind(); err != nil {
		return err
	}
	s.scanner = scanner
	s.scanner.Split(jsonrpc.ScannerSplit)
	s.writer = writer
	return nil
}

func (s *ClientLsp) Unbind() error {
	if s.scanner == nil || s.writer == nil {
		return nil
	}

	s.scanner = nil
	if err := s.writer.Flush(); err != nil {
		return err
	}
	s.writer = nil
	return nil
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

	for {
		select {
		case <-scanCtx.Done():
			s.plane.Log().Infof("%s: done", s.Name())
			return nil
		case bytes := <-ch:
			s.handle(bytes)
		}
	}
}

func (s *ClientLsp) Close() error {
	if err := s.Unbind(); err != nil {
		return err
	}
	s.plane.Log().Infof("%s: closed", s.Name())
	return nil
}
