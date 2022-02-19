package main

import (
	"io"
	"log"
	"os"
	"time"

	remotelog "github.com/fabiankachlock/remote-log"
)

func main() {
	logger := log.Default()
	mw := io.MultiWriter(os.Stdout, remotelog.RemoteLog)
	logger.SetOutput(mw)

	go remotelog.Start("127.0.0.1", 10341)

	for i := 0; i <= 100; i++ {
		logger.Println(i)
		<-time.After(time.Second)
	}
}
