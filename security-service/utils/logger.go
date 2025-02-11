package utils

import (
	"log"
	"os"
)

/**
* Gin Debugger vs own custom logger (?)
* gin is limited and no control over it
* own custom logger can be customized and controlled
* will be use to track in the server
**/

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
)

type Logger struct {
	loggers map[LogLevel]*log.Logger
}

func NewLogger() (*Logger, error) {
	return &Logger{
		loggers: map[LogLevel]*log.Logger{
			DEBUG: log.New(os.Stdout, colorBlue+"[DEBUG] "+colorReset, log.LstdFlags),
			INFO:  log.New(os.Stdout, colorGreen+"[INFO]  "+colorReset, log.LstdFlags),
			WARN:  log.New(os.Stdout, colorYellow+"[WARN]  "+colorReset, log.LstdFlags),
			ERROR: log.New(os.Stdout, colorRed+"[ERROR] "+colorReset, log.LstdFlags),
			FATAL: log.New(os.Stdout, colorRed+"[FATAL] "+colorReset, log.LstdFlags),
		},
	}, nil
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.loggers[DEBUG].Printf(format, v...)
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.loggers[INFO].Printf(format, v...)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	l.loggers[WARN].Printf(format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.loggers[ERROR].Printf(format, v...)
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	l.loggers[FATAL].Printf(format, v...)
	os.Exit(1)
}