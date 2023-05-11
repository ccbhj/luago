package llog

import (
	"fmt"
	"log"
	"os"
)

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

type LogLevel int

func (l LogLevel) String() string {
	var m = [...]string{
		"DEBUG",
		"INFO",
		"WARN",
		"ERROR",
		"FATAL",
	}
	if l < 0 || int(l) >= len(m) {
		return "DEBUG"
	}
	return m[l]
}

var (
	Debug func(f string, args ...interface{})
	Info  func(f string, args ...interface{})
	Warn  func(f string, args ...interface{})
	Error func(f string, args ...interface{})
	Fatal func(f string, args ...interface{})
)

var logLevel LogLevel
var logger *log.Logger

func init() {
	initLogLevel()
	initLogger()
}

func initLogLevel() {
	var m = map[string]LogLevel{
		LogLevelDebug.String(): LogLevelDebug,
		LogLevelInfo.String():  LogLevelInfo,
		LogLevelWarn.String():  LogLevelWarn,
		LogLevelError.String(): LogLevelError,
		LogLevelFatal.String(): LogLevelFatal,
	}
	l, in := m[os.Getenv("LOG_LEVEL")]
	if !in {
		logLevel = LogLevelDebug
		return
	}
	logLevel = l
}

func initLogger() {
	logger = log.New(os.Stdout, "[LUAGO]", 0)

	Debug = wrapLogFn(LogLevelDebug, "D")
	Info = wrapLogFn(LogLevelInfo, "I")
	Warn = wrapLogFn(LogLevelWarn, "W")
	Error = wrapLogFn(LogLevelError, "E")
	Fatal = wrapLogFn(LogLevelFatal, "F")
}

func wrapLogFn(level LogLevel, prefix string) func(string, ...interface{}) {
	return func(f string, args ...interface{}) {
		if level < logLevel {
			return
		}
		s := fmt.Sprintf("["+prefix+"] "+f, args...)
		if level == LogLevelFatal {
			logger.Fatal(s)
		} else {
			logger.Print(s)
		}
	}
}
