package stdio

import "fmt"

func (s *LspServerStdio) errorNotRunning() error {
	return fmt.Errorf("%v: not running", s.Name())
}

func (s *LspServerStdio) errorRunning() error {
	return fmt.Errorf("%v: running", s.Name())
}
