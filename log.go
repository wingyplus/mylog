//go:generate stringer -type=Severity

package mylog

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/mgutz/ansi"
)

func init() {
	if os.Getenv("MYLOGCOLOR_DISABLED") == "1" {
		ansi.DisableColors(true)
	}
}

var (
	now = time.Now
	pid = os.Getpid()
	red = ansi.ColorFunc("red")
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

const logTimeFormat = "15:04:05.999999"

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
	timestamp := now().Format(logTimeFormat)
	msg := fmt.Sprint(v...)
	if (severity & logger.allowedSeverities) != 0 {
		logger.mu.Lock()
		fmt.Fprintf(logger.writer, "%s|%s|%d|%s\n", red(severity.String()), timestamp, pid, msg)
		logger.mu.Unlock()
	}
}

func Info(v ...interface{})  { logger.Write(INFO, v...) }
func Debug(v ...interface{}) { logger.Write(DEBUG, v...) }
func Error(v ...interface{}) { logger.Write(ERROR, v...) }
func Warn(v ...interface{})  { logger.Write(WARN, v...) }
func Fatal(v ...interface{}) { logger.Write(FATAL, v...) }

func SetOutput(w io.Writer) {
	logger.SetOutput(w)
}

func SetAllowedSeverities(severities Severity) {
	logger.SetAllowedSeverities(severities)
}
