package log

import (
	"bytes"
	"testing"
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
