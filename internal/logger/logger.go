package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

type Logger struct {
	level  Level
	output io.Writer
	logger *log.Logger
}

var defaultLogger *Logger

func init() {
	defaultLogger = New(INFO, os.Stderr)
}

func New(level Level, output io.Writer) *Logger {
	return &Logger{
		level:  level,
		output: output,
		logger: log.New(output, "", 0), // No default prefix/flags
	}
}

func (l *Logger) log(level Level, msg string, args ...interface{}) {
	if level < l.level {
		return
	}

	timestamp := time.Now().Format("15:04:05")
	prefix := fmt.Sprintf("[%s %s]", timestamp, level.String())
	
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	
	l.logger.Printf("%s %s", prefix, msg)
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.log(DEBUG, msg, args...)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.log(INFO, msg, args...)
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	l.log(WARN, msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.log(ERROR, msg, args...)
}

func (l *Logger) SetLevel(level Level) {
	l.level = level
}

// Package-level convenience functions
func Debug(msg string, args ...interface{}) {
	defaultLogger.Debug(msg, args...)
}

func Info(msg string, args ...interface{}) {
	defaultLogger.Info(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	defaultLogger.Warn(msg, args...)
}

func Error(msg string, args ...interface{}) {
	defaultLogger.Error(msg, args...)
}

func SetLevel(level Level) {
	defaultLogger.SetLevel(level)
}

func SetVerbose(verbose bool) {
	if verbose {
		SetLevel(DEBUG)
	} else {
		SetLevel(INFO)
	}
}
