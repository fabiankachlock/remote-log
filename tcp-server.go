package remotelog

import (
	"bufio"
	"fmt"
	"net"

	"github.com/google/uuid"
)

type tcpServer struct {
	tcpConnections map[string]net.Conn
}

func (s tcpServer) start(host string, port int) error {
	addr := fmt.Sprintf("%s:%d", host, port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			return err
		}

		id := uuid.NewString()
		s.tcpConnections[id] = c
		go s.handleConn(c, id)
	}
}

func (s tcpServer) handleConn(c net.Conn, id string) {
	for {
		_, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println("Error", err)
			delete(s.tcpConnections, id)
			return
		}
	}
}

func (s tcpServer) write(p []byte) {
	for _, conn := range s.tcpConnections {
		_, err := conn.Write(p)
		if err != nil {
			resultsChan <- err
			return
		}
	}
	doneChan <- true
}
