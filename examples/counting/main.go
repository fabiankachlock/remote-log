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

	s := remotelog.NewTcp()
	s.Listen("127.0.0.1", 10341)

	go func() {
		<-time.After(time.Second * 10)
		s.Close()
	}()

	for i := 0; i <= 20; i++ {
		logger.Println(i)
		<-time.After(time.Second)
	}
}
