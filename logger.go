package remotelog

import (
	"io"
	"log"
)

var (
	Writer io.Writer = remoteLogger{}
)

func NewLogger() *log.Logger {
	logger := log.Default()
	logger.SetOutput(Writer)
	return logger
}

type remoteLogger struct{}

func (r remoteLogger) Write(p []byte) (n int, err error) {
	n = len(p)

	for _, server := range activeServers {
		server.write(p)
	}

	select {
	case errRes := <-resultsChan:
		err = errRes
	case <-doneChan:
		break
	}

	return
}
