package client

func (c *Client) LogInfoPrintf(format string, v ...any) {
	c.logger.Printf("info: %v: "+format, append([]any{c.conn.RemoteAddr()}, v...)...)
}

func (c *Client) LogWarnPrintf(format string, v ...any) {
	c.logger.Printf("warn: %v: "+format, append([]any{c.conn.RemoteAddr()}, v...)...)
}

func (c *Client) LogErrorPrintf(format string, v ...any) {
	c.logger.Printf("error: %v: "+format, append([]any{c.conn.RemoteAddr()}, v...)...)
}
