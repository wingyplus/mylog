package mylog

import (
	"bufio"
	"bytes"
	"testing"

	"gopkg.in/natefinch/lumberjack.v2"
)

func TestLog(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetAllowedSeverities(ALL)

	Log(DEBUG, "Hello World")

	if s := buf.String(); s != "DEBUG|Hello World\n" {
		t.Errorf("Expect DEBUG|Hello World but got %s", s)
	}
}

func TestLog_DoesNotPrintWhenNotAllowed(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetAllowedSeverities(INFO | WARN)

	Log(DEBUG, "Hello World")

	if s := buf.String(); s != "" {
		t.Errorf("Expect empty string but got %s", s)
	}
}

func benchmarkLog(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Log(DEBUG, "Hello World")
			Log(INFO, "Hello World")
			Log(WARN, "Hello World")
			Log(ERROR, "Hello World")
			Log(FATAL, "Hello World")
		}
	})
}

func newLogger(filename string) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    100,
		MaxAge:     1,
		MaxBackups: 10,
	}
}

func BenchmarkLog_Lumberjack(b *testing.B) {
	logger := newLogger("hello-lumberjack.log")
	defer logger.Close()

	SetOutput(logger)
	SetAllowedSeverities(ALL)

	benchmarkLog(b)
}

func BenchmarkLog_BufferedLumberjack(b *testing.B) {
	logger := newLogger("hello-buffered-lumberjack.log")
	defer logger.Close()

	SetOutput(bufio.NewWriter(logger))
	SetAllowedSeverities(ALL)

	benchmarkLog(b)
}
