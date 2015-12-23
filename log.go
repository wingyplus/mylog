//go:generate stringer -type=Severity

package mylog

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
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
	now    = time.Now
	pid    = os.Getpid()
	caller = func() (funcname string, file string, line int) {
		pc, file, line, ok := runtime.Caller(3)
		if !ok {
			line = 0
			file = ""
		}
		funcname = runtime.FuncForPC(pc).Name()
		return funcname, file, line
	}

	red    = ansi.ColorFunc("red")
	green  = ansi.ColorFunc("green")
	yellow = ansi.ColorFunc("yellow")

	severityColors = map[Severity]func(string) string{
		DEBUG: ansi.ColorFunc("red"),
		ERROR: ansi.ColorFunc("red"),
		FATAL: ansi.ColorFunc("red"),
		INFO:  ansi.ColorFunc("green"),
		WARN:  ansi.ColorFunc("yellow"),
	}
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
	var colorFn func(string) string
	switch severity {
	case ERROR, FATAL, DEBUG:
		colorFn = red
	case INFO:
		colorFn = green
	case WARN:
		colorFn = yellow
	default:
		colorFn = func(s string) string { return s }
	}

	funcname, file, line := caller()
	timestamp := now().Format(logTimeFormat)
	msg := colorFn(fmt.Sprintf("%s|%s|%d|%s|%s|%d|%s", severity.String(), timestamp, pid, file, funcname, line, fmt.Sprint(v...)))

	if (severity & logger.allowedSeverities) != 0 {
		logger.mu.Lock()
		logger.writer.Write([]byte(msg + "\n"))
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
