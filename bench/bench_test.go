package bench

import (
	"bufio"
	"testing"

	"github.com/Sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func BenchmarkLogrus(b *testing.B) {
	writer := &lumberjack.Logger{
		Filename: "./logrus.log",
		MaxAge:   1,
		MaxSize:  100,
	}
	defer writer.Close()

	logrus.SetOutput(bufio.NewWriter(writer))
	logrus.SetLevel(logrus.DebugLevel)
	logger := logrus.WithFields(logrus.Fields{
		"app": "myapp",
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debug("Hello World")
	}
}
