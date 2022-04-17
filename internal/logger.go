package internal

import (
	"log"
	"os"
)

var (
	logger     *log.Logger
	loggerPath = "logs.txt"
)

func init() {
	os.Create(loggerPath)
}

func newLogger(path string) *log.Logger {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}

	logger := &log.Logger{}
	logger.SetFlags(log.LstdFlags | log.Ldate | log.Llongfile)
	logger.SetOutput(file)

	return logger
}

func GetGlobalLoggerInstance() *log.Logger {
	if logger == nil {
		logger = newLogger(loggerPath)
	}

	return logger
}
