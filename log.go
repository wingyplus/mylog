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
	writer        io.Writer
	allowedLevels Severity

	mu sync.Mutex
}

var logger = &Logger{
	writer:        bufio.NewWriter(os.Stdout),
	allowedLevels: ALL,
}

func (logger *Logger) SetOutput(w io.Writer) {
	logger.writer = w
}

func (logger *Logger) SetAllowedSeverities(lvls Severity) {
	logger.allowedLevels = lvls
}

func (logger *Logger) Write(lv Severity, v ...interface{}) {
	if (lv & logger.allowedLevels) != 0 {
		logger.mu.Lock()
		defer logger.mu.Unlock()

		fmt.Fprint(logger.writer, fmt.Sprintf("%s|%s\n", lv, fmt.Sprint(v...)))
	}
}

func Log(lv Severity, v ...interface{}) {
	logger.Write(lv, v...)
}

func SetOutput(w io.Writer) {
	logger.SetOutput(w)
}

func SetAllowedSeverities(lv Severity) {
	logger.SetAllowedSeverities(lv)
}
