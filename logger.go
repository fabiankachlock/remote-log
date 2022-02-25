package remotelog

import (
	"fmt"
	"log"
)

type logWriter struct {
	servers *sharedServerSlice
	results chan error
	done    chan bool
}

func (r RemoteLog) NewLogger() *log.Logger {
	logger := log.Default()
	logger.SetOutput(r.Writer)
	return logger
}

func (r logWriter) Write(p []byte) (n int, err error) {
	n = len(p)

	if len(r.servers.slc) == 0 {
		fmt.Println("no servers")
		return
	}

	for _, server := range r.servers.slc {
		server.write(p)
	}

	errCount := 0
	select {
	case errRes := <-r.results:
		// there might appear multiple errors, but they must get shadowed to conform the the writer interface
		errCount += 1
		err = fmt.Errorf("[%d] errors; last: %s", errCount, errRes.Error())
	case <-r.done:
		break
	}

	return
}
