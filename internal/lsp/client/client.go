package client

import (
	"bufio"
	"context"
	"fmt"
	"mals/internal/jsonrpc"
	"mals/internal/middleware"
	"mals/internal/plane"
	"mals/third_party/lsp"
)

type LspClient struct {
	plane      plane.Plane
	scanner    *bufio.Scanner
	writer     *bufio.Writer
	middleware *middleware.Middleware
}

func errorParseUnexpectedType[T jsonrpc.Message](s *LspClient) {
	var dummy T

	resp := jsonrpc.Response{
		Error: &jsonrpc.Error{
			Code:    int32(lsp.ParseError),
			Message: fmt.Sprintf("message is not of type %T", dummy),
		},
	}

	s.plane.Warnf("%T %v: %v", s, s.Name(), resp.Error.Message)

	s.send(&resp)
}

func New(listenerName string, clientName string, plane plane.Plane, scanner *bufio.Scanner, writer *bufio.Writer) *LspClient {
	s := &LspClient{
		plane:      plane,
		scanner:    scanner,
		writer:     writer,
		middleware: middleware.New(listenerName, clientName, plane),
	}

	s.scanner.Split(jsonrpc.ScannerSplit)

	return s
}

func (s *LspClient) Name() string {
	return s.middleware.Name()
}

func (s *LspClient) Serve(ctx context.Context) error {
	s.plane.Infof("%T %v: listen", s, s.Name())

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
				s.plane.Debugf("%T %v: scanned %v", s, s.Name(), string(bytes))

				s.handle(bytes)
			}
		}
	}()

	<-scanCtx.Done()

	s.plane.Infof("%T %v: done", s, s.Name())

	if err := s.writer.Flush(); err != nil {
		s.plane.Errorf("%T %v: flush: %v", s, s.Name(), err)
		return err
	}

	s.plane.Infof("%T %v: closed", s, s.Name())

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

	s.plane.Debugf("%T %s: sent %s", s, s.Name(), string(bytes))

	return nil
}
