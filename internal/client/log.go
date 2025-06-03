package client

import "fmt"

func (c *Client) logPrintf(level, format string, v ...any) {
	c.logger.Printf("%s: %v: %s", level, c.conn.RemoteAddr(), fmt.Sprintf(format, v...))
}

func (c *Client) LogInfoPrintf(format string, v ...any) {
	c.logPrintf("info", format, v...)
}

func (c *Client) LogWarnPrintf(format string, v ...any) {
	c.logPrintf("warn", format, v...)
}

func (c *Client) LogErrorPrintf(format string, v ...any) {
	c.logPrintf("error", format, v...)
}
