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

type SimpleLogger struct {
	*log.Logger
}

func New() *SimpleLogger {
	return &SimpleLogger{
		Logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}
}

func (l *SimpleLogger) Debug(args ...interface{}) {
	l.SetPrefix("DEBUG: ")
	l.Println(args...)
}

func (l *SimpleLogger) Info(args ...interface{}) {
	l.SetPrefix("INFO: ")
	l.Println(args...)
}

func (l *SimpleLogger) Warn(args ...interface{}) {
	l.SetPrefix("WARN: ")
	l.Println(args...)
}

func (l *SimpleLogger) Error(args ...interface{}) {
	l.SetPrefix("ERROR: ")
	l.Println(args...)
}

func (l *SimpleLogger) Fatal(args ...interface{}) {
	l.SetPrefix("FATAL: ")
	l.Println(args...)
	os.Exit(1)
}
