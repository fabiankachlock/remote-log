package remotelog

import (
	"fmt"
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

	errCount := 0
	select {
	case errRes := <-resultsChan:
		// there might appear multiple errors, but they must get shadowed to conform the the writer interface
		errCount += 1
		err = fmt.Errorf("[%d] errors; last: %s", errCount, errRes.Error())
	case <-doneChan:
		break
	}

	return
}
