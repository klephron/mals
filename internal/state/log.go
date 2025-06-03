package state

import "fmt"

func (c *State) LogInfoPrintf(format string, v ...any) {
	c.logger.Printf("info: state: %s", fmt.Sprintf(format, v...))
}

func (c *State) LogWarnPrintf(format string, v ...any) {
	c.logger.Printf("warn: state: %s", fmt.Sprintf(format, v...))
}

func (c *State) LogErrorPrintf(format string, v ...any) {
	c.logger.Printf("error: state: %s", fmt.Sprintf(format, v...))
}
