package client

import (
	"bufio"
	"context"
	"fmt"
	"mals/internal/control/controller"
	"mals/internal/jsonrpc"
	"mals/internal/lsp/client"
)

type ClientGeneric struct {
	client.Client
	controller *controller.Controller
	scanner    *bufio.Scanner
	writer     *bufio.Writer
}

func New(controller *controller.Controller) *ClientGeneric {
	s := &ClientGeneric{
		controller: controller,
		scanner:    nil,
		writer:     nil,
	}
	return s
}

func (s *ClientGeneric) Name() string {
	return "client"
}

func (s *ClientGeneric) Bind(scanner *bufio.Scanner, writer *bufio.Writer) {
	s.Unbind()
	s.scanner = scanner
	s.scanner.Split(jsonrpc.ScannerSplit)
	s.writer = writer
}

func (s *ClientGeneric) Unbind() error {
	s.scanner = nil
	if err := s.writer.Flush(); err != nil {
		return err
	}
	s.writer = nil
	return nil
}

func (s *ClientGeneric) ServeBinded(ctx context.Context) error {
	if s.scanner == nil || s.writer == nil {
		return fmt.Errorf("%s: must be binded before serve: scanner or writer is nil", s.Name())
	}

	ch := make(chan []byte)
	defer close(ch)

	s.controller.Info(fmt.Sprintf("%s: listening", s.Name()))

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
				s.controller.Debug(fmt.Sprintf("%s: scanned %s", s.Name(), string(s.scanner.Bytes())))
			}
		}
	}()

	for {
		select {
		case <-scanCtx.Done():
			s.controller.Info("%s: scanner done", s.Name())
			return nil
		case bytes := <-ch:
			s.LspHandle(bytes)
		}
	}
}

func (s *ClientGeneric) Close() error {
	s.Unbind()
	s.controller.Info("%s: closed", s.Name())
	return nil
}
