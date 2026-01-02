package stdio

import (
	"bufio"
	"context"
	"fmt"
	"mals/internal/jsonrpc"
	"mals/internal/lsp/server"
	"mals/internal/plane"
	"mals/pkg/config"
	"os/exec"
	"sync"

	"github.com/puzpuzpuz/xsync/v4"
)

type RequestValue struct {
	request *jsonrpc.Request
	result  chan *jsonrpc.Response
}

type LspServerStdio struct {
	server.LspServer
	name  string
	plane plane.Plane

	cmd []string

	rw sync.RWMutex

	running  bool
	reader   *bufio.Reader
	writer   *bufio.Writer
	requests *xsync.Map[int32, *RequestValue]
}

func New(name string, settings *config.LspSettingsStdio, plane plane.Plane) *LspServerStdio {
	cmd := settings.Cmd

	return &LspServerStdio{
		name:     name,
		plane:    plane,
		cmd:      cmd,
		rw:       sync.RWMutex{},
		running:  false,
		requests: xsync.NewMap[int32, *RequestValue](),
	}
}

func (s *LspServerStdio) Name() string {
	return s.name
}

func (s *LspServerStdio) Kind() string {
	var settings config.LspSettingsStdio
	return settings.Kind()
}

func (s *LspServerStdio) Run(ctx context.Context) error {
	{
		s.rw.Lock()

		if s.running {
			s.rw.Unlock()
			return s.errorRunning()
		}

		cmd := exec.CommandContext(ctx, s.cmd[0], s.cmd[1:]...)

		stdin, err := cmd.StdinPipe()
		if err != nil {
			s.rw.Unlock()
			return fmt.Errorf("%v: stdin pipe: %v", s.Name(), err)
		}

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			s.rw.Unlock()
			return fmt.Errorf("%v: stdout pipe: %v", s.Name(), err)
		}

		if err := cmd.Start(); err != nil {
			s.rw.Unlock()
			return fmt.Errorf("%v: start: %v", s.Name(), err)
		}

		s.reader = bufio.NewReader(stdout)
		s.writer = bufio.NewWriter(stdin)

		go func(reader *bufio.Reader) {
			scanner := bufio.NewScanner(reader)
			scanner.Split(jsonrpc.ScannerSplit)

			for {
				select {
				case <-ctx.Done():
					return
				default:
					if !scanner.Scan() {
						return
					}

					bytes := scanner.Bytes()
					s.plane.Debugf("%v: scanned %v", s.Name(), string(bytes))

					s.handleResponse(bytes)
				}
			}
		}(s.reader)

		s.running = true
		s.rw.Unlock()

		s.plane.Infof("%v: started", s.Name())
	}

	<-ctx.Done()

	{
		s.rw.Lock()
		defer s.rw.Unlock()

		if err := s.writer.Flush(); err != nil {
			s.plane.Errorf("%v: flush: %v", s.Name(), err)
			return err
		}

		s.running = false
		s.reader = nil
		s.writer = nil

		s.requests.Range(func(key int32, value *RequestValue) bool {
			value.result <- nil
			return true
		})
		s.requests.Clear()

		s.plane.Infof("%v: done", s.Name())

		return nil
	}
}

func (s *LspServerStdio) send(msg jsonrpc.Message) error {
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

func (s *LspServerStdio) sendRequest(request *RequestValue) error {
	s.rw.RLock()
	defer s.rw.RUnlock()

	if !s.running {
		return s.errorNotRunning()
	}

	s.requests.Store(request.request.Id, request)

	return s.send(request.request)
}

func (s *LspServerStdio) sendNotification(notification *jsonrpc.Notification) error {
	s.rw.Lock()
	defer s.rw.RUnlock()

	if !s.running {
		return s.errorNotRunning()
	}

	return s.send(notification)
}

func (s *LspServerStdio) handleResponse(bytes []byte) {
	msg, err := jsonrpc.DecodeMessage(bytes)

	if err != nil {
		s.plane.Errorf("%T %v: %v", s, s.Name(), err)
		return
	}

	switch msg := msg.(type) {
	case *jsonrpc.Response:
		value, ok := s.requests.LoadAndDelete(msg.Id)
		if !ok {
			s.plane.Warnf("%T %v: no request with Id=%d in map", s, s.Name(), msg.Id)
			return
		}

		value.result <- msg
		close(value.result)

	case *jsonrpc.Notification:
		s.plane.Debugf("%T %v: received notification: %+v", s, s.Name(), msg)

	default:
		s.plane.Warnf("%T %v: unhandled type %T", s, s.Name(), msg)
	}
}
