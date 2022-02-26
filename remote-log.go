package remotelog

import (
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"

	"github.com/google/uuid"
)

type (
	RemoteLog struct {
		resultsChan   chan error
		doneChan      chan bool
		activeServers *sharedServerSlice
		Writer        io.Writer
		logger        *log.Logger
	}

	sharedServerSlice struct {
		lock sync.Mutex
		slc  []Server
	}

	ServerOptions struct {
		Host string
		Port int
	}

	InstanceOptions struct {
		EnableLogging bool
	}
)

func New(options InstanceOptions) RemoteLog {
	instance := RemoteLog{
		resultsChan: make(chan error, 1),
		doneChan:    make(chan bool, 1),
		activeServers: &sharedServerSlice{
			sync.Mutex{},
			[]Server{},
		},
	}

	if options.EnableLogging {
		instance.logger = log.New(os.Stdout, "", log.Ltime|log.Ldate)
	} else {
		instance.logger = log.New(ioutil.Discard, "", 0)
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
	getId() string
}

func (r RemoteLog) registerServer(s Server) {
	r.activeServers.lock.Lock()
	r.activeServers.slc = append(r.activeServers.slc, s)
	r.activeServers.lock.Unlock()

	go func() {
		// remove server when closed
		<-s.closedChan()
		r.activeServers.lock.Lock()
		r.activeServers.slc = removeServer(r.activeServers.slc, s)
		r.activeServers.lock.Unlock()
	}()
}

func (r RemoteLog) NewTcpServer() Server {
	server := &tcpServer{
		map[string]net.Conn{},
		nil,
		r.resultsChan,
		r.doneChan,
		make(chan bool),
		r.logger,
		uuid.NewString(),
	}
	r.registerServer(server)
	return server
}

func (r RemoteLog) NewUdpServer() Server {
	server := &udpServer{
		[]*net.UDPAddr{},
		nil,
		r.resultsChan,
		r.doneChan,
		make(chan bool),
		r.logger,
		uuid.NewString(),
	}
	r.registerServer(server)
	return server
}
