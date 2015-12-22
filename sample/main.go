package main

import (
	"bufio"
	"io"
	"net/http"
	"os"

	"github.com/wingyplus/mylog"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	logger := &lumberjack.Logger{
		Filename: "./test.log",
		MaxAge:   1,   // days
		MaxSize:  100, // MB
	}
	mylog.SetOutput(io.MultiWriter(os.Stdout, bufio.NewWriter(logger)))

	mylog.Info("Listen on 0.0.0.0:9000")
	http.HandleFunc("/", index)
	http.ListenAndServe(":9000", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	mylog.Info("Incomming request from path ", r.URL.Path)
	mylog.Error("Has error")
}
