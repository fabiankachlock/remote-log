package main

import (
	"io"
	"log"
	"os"
	"time"

	remotelog "github.com/fabiankachlock/remote-log"
)

func main() {
	rm := remotelog.New(remotelog.InstanceOptions{
		EnableLogging: true,
	})
	logger := log.Default()
	mw := io.MultiWriter(os.Stdout, rm.Writer)
	logger.SetOutput(mw)
	// or logger := remotelog.NewLogger()

	s := rm.NewUdpServer()
	s.Listen(remotelog.ServerOptions{Host: "127.0.0.1", Port: 10341})

	go func() {
		<-time.After(time.Second * 10)
		s.Close()
	}()

	for i := 0; i <= 20; i++ {
		logger.Println(i)
		<-time.After(time.Second)
	}
}
