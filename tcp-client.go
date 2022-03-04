package remotelog

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

type tcpClient struct {
	id     string
	conn   net.Conn
	logger *log.Logger
}

func (client tcpClient) connect(options ServerOptions) error {
	addr := fmt.Sprintf("%s:%d", options.Host, options.Port)
	client.logger.SetPrefix("[" + addr + "] ")
	c, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	done := make(chan bool)
	client.logger.Printf("connected to %s\n", addr)

	go (func() {
		for {
			message, err := bufio.NewReader(c).ReadString('\n')
			if errors.Is(err, net.ErrClosed) || err == io.EOF {
				c.Close()
				client.logger.Println("connection closed")
				return
			} else if err != nil {
				client.logger.Printf("Error: %s", err)
				done <- true
				return
			}
			client.logger.Print("-> " + message)
		}
	})()

	<-done
	return nil
}

func (client tcpClient) Close() error {
	return client.conn.Close()
}
