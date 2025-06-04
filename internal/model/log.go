package model

import "fmt"

func (c *Model) logPrintf(level, format string, v ...any) {
	c.logger.Printf("%s: model %v: %s", level, c.Id, fmt.Sprintf(format, v...))
}

func (c *Model) LogInfoPrintf(format string, v ...any) {
	c.logPrintf("info", format, v...)
}

func (c *Model) LogWarnPrintf(format string, v ...any) {
	c.logPrintf("warn", format, v...)
}

func (c *Model) LogErrorPrintf(format string, v ...any) {
	c.logPrintf("error", format, v...)
}
