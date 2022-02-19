package remotelog

import (
	"io"
)

var (
	RemoteLog io.Writer = remoteLogger{}
)

type remoteLogger struct{}

func (r remoteLogger) Write(p []byte) (n int, err error) {
	n = len(p)
	server.write(p)

	select {
	case errRes := <-resultsChan:
		err = errRes
	case <-doneChan:
		break
	}

	return
}
