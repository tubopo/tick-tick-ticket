package logger

import (
	"log"
	"os"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}

type LogLevel int

// Log levels constants.
const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

type SimpleLogger struct {
	*log.Logger
	level LogLevel
}

func New(level LogLevel) *SimpleLogger {
	return &SimpleLogger{
		Logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
		level:  level,
	}
}

func (l *SimpleLogger) Debug(args ...interface{}) {
	if l.level <= LogLevelDebug {
		l.SetPrefix("DEBUG: ")
		l.Println(args...)
	}
}

func (l *SimpleLogger) Info(args ...interface{}) {
	if l.level <= LogLevelInfo {
		l.SetPrefix("INFO: ")
		l.Println(args...)
	}
}

func (l *SimpleLogger) Warn(args ...interface{}) {
	if l.level <= LogLevelWarn {
		l.SetPrefix("WARN: ")
		l.Println(args...)
	}
}

func (l *SimpleLogger) Error(args ...interface{}) {
	if l.level <= LogLevelError {
		l.SetPrefix("ERROR: ")
		l.Println(args...)
	}
}

func (l *SimpleLogger) Fatal(args ...interface{}) {
	l.SetPrefix("FATAL: ")
	l.Println(args...)
	os.Exit(1)
}
