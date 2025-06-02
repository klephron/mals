package logger

import (
	"log"
	"os"
)

func GetLogger() (*log.Logger, error) {
	return log.New(os.Stdout, "", log.LUTC|log.Lshortfile|log.Ldate|log.Ltime), nil
}
