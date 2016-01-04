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
	caller = func() (fn string, file string, line int) {
		return "yourfunc", "test.go", 24
	}
}

func TestLog(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetAllowedSeverities(ALL)

	Debug("Hello World")

	expected := "DEBUG|17:50:22.615673|1234|test.go|yourfunc|24|Hello World\n"
	if s := buf.String(); s != expected {
		t.Errorf("Expect %s but got %s", expected, s)
	}
}

func TestLog_DoesNotPrintWhenNotAllowed(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetAllowedSeverities(INFO | WARN)

	Debug("Hello World")

	if s := buf.String(); s != "" {
		t.Errorf("Expect empty string but got %s", s)
	}
}

func TestLog_Colors(t *testing.T) {
	if os.Getenv("MYLOGCOLOR_DISABLED") == "1" {
		t.SkipNow()
	}

	colorTestCases := []struct {
		log    func(...interface{})
		output string
	}{
		{Debug, "\033[31mDEBUG|17:50:22.615673|1234|test.go|yourfunc|24|Hello World\033[0m\n"},
		{Info, "\033[32mINFO|17:50:22.615673|1234|test.go|yourfunc|24|Hello World\033[0m\n"},
		{Warn, "\033[33mWARN|17:50:22.615673|1234|test.go|yourfunc|24|Hello World\033[0m\n"},
		{Error, "\033[31mERROR|17:50:22.615673|1234|test.go|yourfunc|24|Hello World\033[0m\n"},
		{Fatal, "\033[31mFATAL|17:50:22.615673|1234|test.go|yourfunc|24|Hello World\033[0m\n"},
	}

	SetAllowedSeverities(ALL)

	for _, testcase := range colorTestCases {
		var buf bytes.Buffer
		SetOutput(&buf)

		testcase.log("Hello World")

		if s := buf.String(); s != testcase.output {
			t.Errorf("Expect %s but got %s", testcase.output, s)
		}
	}

}

func TestLog_LineNo(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetAllowedSeverities(ALL)

	Debug("Test Line Number.")

	expected := "DEBUG|17:50:22.615673|1234|test.go|yourfunc|24|Test Line Number.\n"
	if s := buf.String(); s != expected {
		t.Errorf(`Expect
			%s
			but got
			%s`, expected, s)
	}
}

func TestLog_Filename(t *testing.T) {
	caller = func() (fn string, file string, line int) {
		return "yourfunc", "/path/to/gopath/src/ourpackage/test.go", 24
	}
	var buf bytes.Buffer
	SetOutput(&buf)
	SetAllowedSeverities(ALL)

	Debug("Hello World")

	expected := "DEBUG|17:50:22.615673|1234|ourpackage/test.go|yourfunc|24|Hello World\n"
	if out := buf.String(); out != expected {
		t.Errorf("Expect %s but got %s", out, expected)
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

func BenchmarkParallelLog_Lumberjack(b *testing.B) {
	logger := newLogger("hello-lumberjack.log")
	defer logger.Close()

	SetOutput(logger)
	SetAllowedSeverities(ALL)

	benchmarkLog(b)
}

func BenchmarkParallelLog_BufferedLumberjack(b *testing.B) {
	logger := newLogger("hello-buffered-lumberjack.log")
	defer logger.Close()

	SetOutput(bufio.NewWriter(logger))
	SetAllowedSeverities(ALL)

	benchmarkLog(b)
}

func BenchmarkLog_BufferedLumberjack_OneSeverity(b *testing.B) {
	logger := newLogger("hello-buffered-lumberjack-not-parallel.log")
	defer logger.Close()

	SetOutput(bufio.NewWriter(logger))
	SetAllowedSeverities(ALL)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Debug("Hello World")
	}
}
