package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reminder/etc"
	"runtime"
	"strings"
	"time"
)

var Log *Lcg

type Level int

const (
	TimeFormat       = time.TimeOnly
	DEBUG      Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

var logLevels = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}

var levelMap = map[string]Level{"INFO": INFO, "WARN": WARN, "ERROR": ERROR, "FATAL": FATAL}

type Lcg struct {
	level  Level
	logger *log.Logger
}

func getLogFilePath() string {
	return fmt.Sprintf("%s", etc.AppConfig.Server.LogPath)
}

func getLogFileFullPath() string {
	prefixPath := getLogFilePath()
	suffixPath := fmt.Sprintf("%s%s.%s", time.Now().Format(time.DateOnly), etc.AppConfig.Server.LogName, etc.AppConfig.Server.LogExt)
	fmt.Printf("suffixPath: %s", etc.AppConfig.Server.LogExt)
	return fmt.Sprintf("%s%s", prefixPath, suffixPath)
}

func OpenLogFile() *os.File {
	filePath := getLogFileFullPath()
	_, err := os.Stat(filePath)
	switch {
	case os.IsNotExist(err):
		mkDir()
	case os.IsPermission(err):
		FatalF("Permission :%v", err)
	}

	handle, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		FatalF("Fail to OpenFile :%v", err)
	}
	return handle
}

func mkDir() {
	dir, _ := os.Getwd()
	err := os.MkdirAll(dir+"/"+getLogFilePath(), os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func NewLogger(level Level, writers ...io.Writer) *Lcg {
	multiWriter := io.MultiWriter(writers...)
	return &Lcg{
		level:  level,
		logger: log.New(multiWriter, "", log.LstdFlags),
	}
}

func SetLevel(level Level) {
	Log.level = level
}

func lg(level Level, msg string) {
	v, ok := levelMap[strings.ToUpper(etc.AppConfig.Server.LogLevel)]
	if !ok {
		v = INFO
	}
	if level >= v {
		_, file, line, ok := runtime.Caller(3)
		if ok {
			Log.logger.SetPrefix(fmt.Sprintf("[%s][%s:%d]", logLevels[level], filepath.Base(file), line))
		} else {
			Log.logger.SetPrefix(fmt.Sprintf("[%s]", logLevels[level]))
		}
		_ = Log.logger.Output(3, msg)
	}
}

func Debug(msg string) {
	lg(DEBUG, msg)
}

func DebugF(format string, arg ...any) {
	lg(DEBUG, fmt.Sprintf(format, arg))
}

func Info(msg string) {
	lg(INFO, msg)
}

func InfoF(format string, args ...any) {
	lg(INFO, fmt.Sprintf(format, args))
}

func Warn(msg string) {
	lg(WARN, msg)
}

func WarnF(format string, arg ...any) {
	lg(WARN, fmt.Sprintf(format, arg))
}

func Error(msg string) {
	lg(ERROR, msg)
}

func ErrorF(format string, args ...any) {
	lg(ERROR, fmt.Sprintf(format, args))
}

func FatalF(format string, args ...any) {
	lg(FATAL, fmt.Sprintf(format, args))
}
