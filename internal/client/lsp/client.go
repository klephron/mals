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

	state State

	plane   plane.Plane
	scanner *bufio.Scanner
	writer  *bufio.Writer
}

func New(plane plane.Plane, scanner *bufio.Scanner, writer *bufio.Writer) *ClientLsp {
	s := &ClientLsp{
		state: State{
			initialized: false,
		},
		plane:   plane,
		scanner: scanner,
		writer:  writer,
	}
	s.scanner.Split(jsonrpc.ScannerSplit)
	return s
}

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

			if err := s.Client.Close(); err != nil {
				s.plane.Log().Errorf("%s: error while closing %v", err)
				return err
			}

			return nil

		case bytes := <-ch:
			s.handle(bytes)
		}
	}
}

func (s *ClientLsp) Close() error {
	if err := s.writer.Flush(); err != nil {
		return err
	}
	return nil
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
