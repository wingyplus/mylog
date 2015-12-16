//go:generate stringer -type=Severity

package mylog

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
)

type Severity int

const (
	ERROR Severity = 1 << iota
	DEBUG
	FATAL
	INFO
	WARN

	ALL = ERROR | DEBUG | FATAL | INFO | WARN
)

type Logger struct {
	writer            io.Writer
	allowedSeverities Severity

	mu sync.Mutex
}

var logger = &Logger{
	writer:            bufio.NewWriter(os.Stdout),
	allowedSeverities: ALL,
}

func (logger *Logger) SetOutput(w io.Writer) {
	logger.writer = w
}

func (logger *Logger) SetAllowedSeverities(severities Severity) {
	logger.allowedSeverities = severities
}

func (logger *Logger) Write(severity Severity, v ...interface{}) {
	if (severity & logger.allowedSeverities) != 0 {
		logger.mu.Lock()
		fmt.Fprint(logger.writer, fmt.Sprintf("%s|%s\n", severity, fmt.Sprint(v...)))
		logger.mu.Unlock()
	}
}

func Log(severity Severity, v ...interface{}) {
	logger.Write(severity, v...)
}

func Info(v ...interface{})  { Log(INFO, v...) }
func Debug(v ...interface{}) { Log(DEBUG, v...) }
func Error(v ...interface{}) { Log(ERROR, v...) }
func Warn(v ...interface{})  { Log(WARN, v...) }
func Fatal(v ...interface{}) { Log(FATAL, v...) }

func SetOutput(w io.Writer) {
	logger.SetOutput(w)
}

func SetAllowedSeverities(severities Severity) {
	logger.SetAllowedSeverities(severities)
}
