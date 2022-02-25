package remotelog

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/google/uuid"
)

type tcpServer struct {
	tcpConnections map[string]net.Conn
	listener       *net.Listener
	results        chan error
	done           chan bool
	closed         chan bool
	logger         *log.Logger
	id             string
}

func (s *tcpServer) Listen(options ServerOptions) error {
	addr := fmt.Sprintf("%s:%d", options.Host, options.Port)

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.listener = &l
	s.logger.Printf("started tcp server %s at %s\n", s.id, addr)
	go s.acceptConnections()
	return nil
}

func (s *tcpServer) Close() error {
	for _, c := range s.tcpConnections {
		c.Close()
	}
	s.logger.Printf("closed tcp server %s\n", s.id)
	return (*s.listener).Close()
}

func (s *tcpServer) acceptConnections() {
	for {
		c, err := (*s.listener).Accept()
		if err != nil {
			continue
		}

		id := uuid.NewString()
		s.tcpConnections[id] = c
		go s.handleConn(c, id)
	}
}

func (s *tcpServer) handleConn(c net.Conn, id string) {
	for {
		_, err := bufio.NewReader(c).ReadString('\n')
		if errors.Is(err, net.ErrClosed) || err == io.EOF {
			c.Close()
			delete(s.tcpConnections, id)
			return
		} else if err != nil {
			s.logger.Printf("Error: %s\n", err)
			continue
		}
	}
}

func (s *tcpServer) write(p []byte) {
	for _, conn := range s.tcpConnections {
		_, err := conn.Write(p)
		if err != nil {
			s.results <- err
			return
		}
	}
	s.done <- true
}

func (s *tcpServer) closedChan() <-chan bool {
	return s.closed
}
