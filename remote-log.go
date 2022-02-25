package remotelog

import (
	"io"
	"sync"
)

type (
	RemoteLog struct {
		resultsChan   chan error
		doneChan      chan bool
		activeServers *sharedServerSlice
		Writer        io.Writer
	}

	sharedServerSlice struct {
		lock sync.Mutex
		slc  []Server
	}

	ServerOptions struct {
		Host string
		Port int
	}
)

func New() RemoteLog {
	instance := RemoteLog{
		resultsChan: make(chan error, 1),
		doneChan:    make(chan bool, 1),
		activeServers: &sharedServerSlice{
			sync.Mutex{},
			[]Server{},
		},
	}

	instance.Writer = logWriter{
		instance.activeServers,
		instance.resultsChan,
		instance.doneChan,
	}

	return instance
}

func (r RemoteLog) ExitAll() []error {
	errs := []error{}
	for _, s := range r.activeServers.slc {
		err := s.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

type Server interface {
	Listen(options ServerOptions) error
	Close() error
	write(p []byte)
	closedChan() <-chan bool
}

func (r RemoteLog) NewTcpServer() Server {
	server := newTcpServer(r.resultsChan, r.doneChan)

	r.activeServers.lock.Lock()
	r.activeServers.slc = append(r.activeServers.slc, server)
	r.activeServers.lock.Unlock()

	go func() {
		// remove server when closed
		<-server.closedChan()
		r.activeServers.lock.Lock()
		r.activeServers.slc = removeServer(r.activeServers.slc, server)
		r.activeServers.lock.Unlock()
	}()

	return server
}
