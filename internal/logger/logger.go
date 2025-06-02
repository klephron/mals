package logger

import (
	"errors"
	"log"
	"os"
)

func GetLogger(filename string) (*log.Logger, error) {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, errors.New("unable to open log file" + filename)
	}
	return log.New(logfile, "", log.LUTC|log.Llongfile|log.Ldate|log.Ltime), nil
}
