package remotelog

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/google/uuid"
)

type tcpServer struct {
	tcpConnections map[string]net.Conn
	listener       *net.Listener
	closed         chan bool
}

func newTcpServer() Server {
	server := &tcpServer{
		map[string]net.Conn{},
		nil,
		make(chan bool),
	}
	return server
}

func (s *tcpServer) Listen(host string, port int) error {
	addr := fmt.Sprintf("%s:%d", host, port)

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.listener = &l
	go s.acceptConnections()
	return nil
}

func (s *tcpServer) Close() error {
	for _, c := range s.tcpConnections {
		c.Close()
	}
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
			fmt.Println("Error:", err)
			continue
		}
	}
}

func (s *tcpServer) write(p []byte) {
	for _, conn := range s.tcpConnections {
		_, err := conn.Write(p)
		if err != nil {
			resultsChan <- err
			return
		}
	}
	doneChan <- true
}

func (s *tcpServer) closedChan() <-chan bool {
	return s.closed
}
