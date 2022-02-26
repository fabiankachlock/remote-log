package remotelog

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

type udpClient struct {
	id     string
	conn   *net.UDPConn
	logger *log.Logger
}

func (client udpClient) connect(options ServerOptions) error {
	addr := fmt.Sprintf("%s:%d", options.Host, options.Port)
	raddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	c, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return err
	}

	client.conn = c
	_, err = c.Write([]byte("connect"))
	if err != nil {
		return err
	}

	done := make(chan bool)
	go (func() {
		buf := make([]byte, 1024)
		for {
			n, _, err := client.conn.ReadFromUDP(buf)
			if errors.Is(err, net.ErrClosed) || err == io.EOF {
				client.logger.Println("connection closed")
				client.Close()
				done <- true
			} else if err != nil {
				client.logger.Printf("Error: %s", err)
				continue
			}
			message := string(buf[:n])
			fmt.Print("-> " + message)
		}
	})()
	<-done
	return nil
}

func (client udpClient) Close() error {
	client.conn.Write([]byte("disconnect"))
	return client.conn.Close()
}
