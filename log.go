//go:generate stringer -type=Level

package mylog

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
)

type Level int

const (
	ERROR Level = 1 << iota
	DEBUG
	FATAL
	INFO
	WARN

	ALL = ERROR | DEBUG | FATAL | INFO | WARN
)

type Logger struct {
	writer        io.Writer
	allowedLevels Level

	mu sync.Mutex
}

var logger = &Logger{
	writer:        bufio.NewWriter(os.Stdout),
	allowedLevels: ALL,
}

func (logger *Logger) SetOutput(w io.Writer) {
	logger.writer = w
}

func (logger *Logger) SetAllowedLevel(lvls Level) {
	logger.allowedLevels = lvls
}

func (logger *Logger) Write(lv Level, v ...interface{}) {
	if (lv & logger.allowedLevels) != 0 {
		logger.mu.Lock()
		defer logger.mu.Unlock()

		fmt.Fprint(logger.writer, fmt.Sprintf("%s|%s\n", lv, fmt.Sprint(v...)))
	}
}

func Log(lv Level, v ...interface{}) {
	logger.Write(lv, v...)
}

func SetOutput(w io.Writer) {
	logger.SetOutput(w)
}

func SetAllowedLevel(lv Level) {
	logger.SetAllowedLevel(lv)
}
