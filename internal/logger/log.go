package logger

import (
	"fmt"
	"io"
	"log"
)

var Log *Lcg

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

var logNames = []string{"DEBUG", "INFO", "WARN", "ERROR"}

type Lcg struct {
	level  Level
	logger *log.Logger
}

func NewLogger(level Level, writers ...io.Writer) *Lcg {
	multiWriter := io.MultiWriter(writers...)
	return &Lcg{
		level:  level,
		logger: log.New(multiWriter, "", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *Lcg) SetLevel(level Level) {
	l.level = level
}

func (l *Lcg) log(level Level, msg string) {
	if level >= l.level {
		l.logger.SetPrefix(logNames[level] + ": ")
		_ = l.logger.Output(5, msg)
	}
}

func (l *Lcg) Debug(msg string) {
	l.log(DEBUG, msg)
}

func (l *Lcg) DebugF(format string, arg ...any) {
	l.log(DEBUG, fmt.Sprintf(format, arg))
}

func (l *Lcg) Info(msg string) {
	l.log(INFO, msg)
}

func (l *Lcg) InfoF(format string, args ...any) {
	l.log(INFO, fmt.Sprintf(format, args))
}

func (l *Lcg) Warn(msg string) {
	l.log(WARN, msg)
}

func (l *Lcg) WarnF(format string, arg ...any) {
	l.log(WARN, fmt.Sprintf(format, arg))
}

func (l *Lcg) Error(msg string) {
	l.log(ERROR, msg)
}

func (l *Lcg) ErrorF(format string, args ...any) {
	l.log(ERROR, fmt.Sprintf(format, args))
}

func (l *Lcg) FatalF(format string, args ...any) {
	l.log(FATAL, fmt.Sprintf(format, args))
}
