package remotelog

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

type udpServer struct {
	clients  []*net.UDPAddr
	listener *net.UDPConn
	results  chan error
	done     chan bool
	closed   chan bool
	logger   *log.Logger
	id       string
}

func (s *udpServer) getId() string {
	return s.id
}

func (s *udpServer) Listen(options ServerOptions) error {
	addr := fmt.Sprintf("%s:%d", options.Host, options.Port)
	laddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	l, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return err
	}
	s.listener = l
	s.logger.Printf("started udp server %s at %s\n", s.id, addr)
	go s.connectClients()
	return nil
}

func (s *udpServer) Close() error {
	if s.clients == nil {
		// already closed
		return nil
	}
	s.closed <- true // remove server
	s.closed <- true // quit incoming connections

	close(s.done)
	close(s.results)
	close(s.closed)

	err := (*s.listener).Close()
	s.logger.Printf("closed udp server %s\n", s.id)

	s.clients = nil
	s.listener = nil
	return err
}

func (s *udpServer) connectClients() {
	buf := make([]byte, 1024)
	for {
		select {
		case <-s.closed:
			return
		default:
			n, addr, err := s.listener.ReadFromUDP(buf)

			if errors.Is(err, net.ErrClosed) || err == io.EOF {
				s.clients = removeUDPAddr(s.clients, addr)
				continue
			} else if err != nil {
				s.logger.Printf("Error: %s\n", err)
				continue
			}

			if string(buf[:n]) == "connect" {
				s.clients = append(s.clients, addr)
			} else if string(buf[:n]) == "disconnect" {
				s.clients = removeUDPAddr(s.clients, addr)
			}
		}
	}
}

func (s *udpServer) write(p []byte) {
	for _, addr := range s.clients {
		_, err := s.listener.WriteToUDP(p, addr)
		if err != nil {
			s.results <- err
			return
		}
	}
	s.done <- true
}

func (s *udpServer) closedChan() <-chan bool {
	return s.closed
}
