package mylog

import (
	"bufio"
	"bytes"
	"os"
	"testing"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	now = func() time.Time {
		t, _ := time.Parse(time.RFC3339Nano, "2015-12-21T17:50:22.615673Z")
		return t
	}
	pid = 1234
}

func TestLog(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetAllowedSeverities(ALL)

	logger.Write(DEBUG, "Hello World")

	expected := "DEBUG|17:50:22.615673|1234|Hello World\n"
	if s := buf.String(); s != expected {
		t.Errorf("Expect %s but got %s", expected, s)
	}
}

func TestLog_DoesNotPrintWhenNotAllowed(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetAllowedSeverities(INFO | WARN)

	logger.Write(DEBUG, "Hello World")

	if s := buf.String(); s != "" {
		t.Errorf("Expect empty string but got %s", s)
	}
}

func TestLog_Colors(t *testing.T) {
	if os.Getenv("MYLOGCOLOR_DISABLED") == "1" {
		t.SkipNow()
	}

	colorTestCases := []struct {
		severity Severity
		output   string
	}{
		{DEBUG, "\033[31mDEBUG\033[0m|17:50:22.615673|1234|Hello World\n"},
		{INFO, "\033[32mINFO\033[0m|17:50:22.615673|1234|Hello World\n"},
		{WARN, "\033[33mWARN\033[0m|17:50:22.615673|1234|Hello World\n"},
		{ERROR, "\033[31mERROR\033[0m|17:50:22.615673|1234|Hello World\n"},
		{FATAL, "\033[31mFATAL\033[0m|17:50:22.615673|1234|Hello World\n"},
	}

	SetAllowedSeverities(ALL)

	for _, testcase := range colorTestCases {
		var buf bytes.Buffer
		SetOutput(&buf)

		logger.Write(testcase.severity, "Hello World")

		if s := buf.String(); s != testcase.output {
			t.Errorf("Expect %s but got %s", testcase.output, s)
		}
	}

}

func benchmarkLog(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Debug("Hello World")
			Info("Hello World")
			Warn("Hello World")
			Error("Hello World")
			Fatal("Hello World")
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
