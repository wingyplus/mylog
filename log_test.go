package mylog

import (
	"bufio"
	"bytes"
	"io"
	"sync"
	"testing"

	"gopkg.in/natefinch/lumberjack.v2"
)

func TestLog(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetAllowedLevel(ALL)

	Log(DEBUG, "Hello World")

	if s := buf.String(); s != "DEBUG|Hello World\n" {
		t.Errorf("Expect DEBUG|Hello World but got %s", s)
	}
}

func TestLog_DoesNotPrintWhenNotAllowed(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetAllowedLevel(INFO | WARN)

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

type writer struct {
	io.Writer
	mu sync.Mutex
}

func (w *writer) Write(b []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.Writer.Write(b)
}

func newWriter(w io.Writer) *writer {
	return &writer{
		Writer: w,
	}
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
	SetAllowedLevel(ALL)

	benchmarkLog(b)
}

func BenchmarkLog_BufferedLumberjack(b *testing.B) {
	logger := newLogger("hello-buffered-lumberjack.log")
	defer logger.Close()

	SetOutput(bufio.NewWriter(logger))
	SetAllowedLevel(ALL)

	benchmarkLog(b)
}

func BenchmarkLog_MutexBufferedLumberjack(b *testing.B) {
	logger := newLogger("hello-mutex-buffered-lumberjack.log")
	defer logger.Close()

	SetOutput(newWriter(bufio.NewWriter(logger)))
	SetAllowedLevel(ALL)

	benchmarkLog(b)
}
