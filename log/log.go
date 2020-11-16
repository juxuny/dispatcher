package log

import (
	"fmt"
)

type LoggerLevel int

func (t LoggerLevel) Int() int {
	return int(t)
}

func SetLevel(level LoggerLevel) {
	globalLoggerLevel = level
}

var globalLoggerLevel = LoggerLevelInfo

const (
	LoggerLevelDebug = LoggerLevel(1)
	LoggerLevelInfo  = LoggerLevel(2)
	LoggerLevelWarn  = LoggerLevel(3)
	LoggerLevelError = LoggerLevel(4)
)

type Logger struct {
	level  LoggerLevel
	prefix string
}

func NewLogger(prefix ...string) *Logger {
	if len(prefix) > 0 {
		return &Logger{
			level:  globalLoggerLevel,
			prefix: prefix[0],
		}
	}
	return &Logger{
		level:  globalLoggerLevel,
		prefix: "[default]",
	}
}

func (t *Logger) output(v ...interface{}) {
	if t.prefix != "" {
		v = append([]interface{}{t.prefix}, v...)
	}
	fmt.Println(v...)
}

func (t *Logger) Debug(v ...interface{}) {
	if t.level < LoggerLevelDebug {
		return
	}
	t.output(append([]interface{}{"[debug]"}, v...)...)
}

func (t *Logger) Info(v ...interface{}) {
	if t.level < LoggerLevelInfo {
		return
	}
	t.output(append([]interface{}{"[info]"}, v...)...)
}

func (t *Logger) Warn(v ...interface{}) {
	if t.level < LoggerLevelWarn {
		return
	}
	t.output(append([]interface{}{"warn"}, v...)...)
}

func (t *Logger) Error(v ...interface{}) {
	if t.level < LoggerLevelError {
		return
	}
	t.output(append([]interface{}{"error"}, v...)...)
}
