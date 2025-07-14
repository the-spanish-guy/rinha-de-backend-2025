package config

import (
	"io"
	"log"
	"os"
)

var (
	logger *Logger
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

type Logger struct {
	err     *log.Logger
	info    *log.Logger
	debug   *log.Logger
	warning *log.Logger
	writer  io.Writer
}

func NewLogger(p string) *Logger {
	writer := io.Writer(os.Stdout)
	logger := log.New(writer, p, log.Ldate|log.Ltime)

	return &Logger{
		err:     log.New(writer, Red+"ERROR: ", logger.Flags()),
		debug:   log.New(writer, Cyan+"DEBUG: ", logger.Flags()),
		info:    log.New(writer, Reset+"INFO: ", logger.Flags()),
		warning: log.New(writer, Yellow+"WARNING: ", logger.Flags()),
		writer:  writer,
	}
}

func GetLogger(p string) *Logger {
	return NewLogger(p)
}

func (l *Logger) Debug(v ...interface{}) {
	l.debug.Println(v...)
}

func (l *Logger) Info(v ...interface{}) {
	l.info.Println(v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.err.Println(v...)
}

func (l *Logger) Warning(v ...interface{}) {
	l.warning.Println(v...)
}

func (l *Logger) Debugf(f string, v ...interface{}) {
	l.debug.Printf(f, v...)
}

func (l *Logger) Infof(f string, v ...interface{}) {
	l.info.Printf(f, v...)
}

func (l *Logger) Errorf(f string, v ...interface{}) {
	l.err.Printf(f, v...)
}

func (l *Logger) Warningf(f string, v ...interface{}) {
	l.warning.Printf(f, v...)
}
