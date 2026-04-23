package stdio

import (
	"bufio"
	"context"
	"fmt"
	"mals/internal/jsonrpc"
	"mals/internal/lsp/protocol"
	"mals/internal/plane"
	"mals/pkg/config"
	"os/exec"
	"sync"
	"sync/atomic"

	"github.com/puzpuzpuz/xsync/v4"
)

type lspTask struct {
	request  *jsonrpc.Request
	response chan *jsonrpc.Response
}

type LspServerStdio struct {
	name  string
	plane plane.Plane

	cmd []string

	rw           sync.RWMutex
	running      bool
	reader       *bufio.Reader
	writer       *bufio.Writer
	capabilities *protocol.ServerCapabilities
	info         *protocol.ServerInfo

	tasks *xsync.Map[int32, *lspTask]
	taskc atomic.Int32
}

func newLspTask(request *jsonrpc.Request) *lspTask {
	return &lspTask{
		request:  request,
		response: make(chan *jsonrpc.Response, 1),
	}
}

func (s *LspServerStdio) errorNotRunning() error {
	return fmt.Errorf("%v: not running", s.Name())
}

func (s *LspServerStdio) errorRunning() error {
	return fmt.Errorf("%v: running", s.Name())
}

func New(name string, api *config.LspApiStdio, plane plane.Plane) *LspServerStdio {
	cmd := api.Cmd

	return &LspServerStdio{
		name:         name,
		plane:        plane,
		cmd:          cmd,
		rw:           sync.RWMutex{},
		running:      false,
		reader:       nil,
		writer:       nil,
		capabilities: nil,
		info:         nil,
		tasks:        xsync.NewMap[int32, *lspTask](),
		taskc:        atomic.Int32{},
	}
}

func (s *LspServerStdio) Name() string {
	return s.name
}

func (s *LspServerStdio) Kind() string {
	var settings config.LspApiStdio
	return settings.LspApiKind()
}

func (s *LspServerStdio) Run(ctx context.Context, onReady func()) error {
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
				s.plane.Debugf("%T %v: scanned %v", s, s.Name(), string(bytes))

				s.handle(bytes)
			}
		}
	}(s.reader)

	s.running = true
	s.rw.Unlock()

	s.plane.Infof("%T %v: started", s, s.Name())

	onReady()

	<-ctx.Done()

	s.rw.Lock()
	defer s.rw.Unlock()

	if err := s.writer.Flush(); err != nil {
		s.plane.Errorf("%T %v: flush: %v", s, s.Name(), err)
		return err
	}

	s.running = false
	s.reader = nil
	s.writer = nil

	s.tasks.Range(func(key int32, value *lspTask) bool {
		value.response <- nil
		return true
	})
	s.tasks.Clear()

	s.taskc.Store(0)

	s.plane.Infof("%T %v: done", s, s.Name())

	return nil
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

func (s *LspServerStdio) sendRequest(request *jsonrpc.Request) (<-chan *jsonrpc.Response, error) {
	s.rw.RLock()
	defer s.rw.RUnlock()

	if !s.running {
		return nil, s.errorNotRunning()
	}

	request.Id = s.taskc.Add(1)

	value := newLspTask(request)
	s.tasks.Store(request.Id, value)

	return value.response, s.send(request)
}

func (s *LspServerStdio) sendNotification(notification *jsonrpc.Notification) error {
	s.rw.RLock()
	defer s.rw.RUnlock()

	if !s.running {
		return s.errorNotRunning()
	}

	return s.send(notification)
}

func (s *LspServerStdio) handle(bytes []byte) {
	msg, err := jsonrpc.DecodeMessage(bytes)

	if err != nil {
		s.plane.Errorf("%T %v: %v", s, s.Name(), err)
		return
	}

	switch msg := msg.(type) {
	case *jsonrpc.Response:
		value, ok := s.tasks.LoadAndDelete(msg.Id)
		if !ok {
			s.plane.Warnf("%T %v: no request with Id=%d in map", s, s.Name(), msg.Id)
			return
		}

		value.response <- msg
		close(value.response)

	case *jsonrpc.Notification:
		s.plane.Debugf("%T %v: received notification: %+v", s, s.Name(), msg)

	default:
		s.plane.Warnf("%T %v: unhandled type %T", s, s.Name(), msg)
	}
}
