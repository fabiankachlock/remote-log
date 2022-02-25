package remotelog

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
)

type tcpClient struct{}

func (client tcpClient) connect(host string, port int) error {
	addr := fmt.Sprintf("%s:%d", host, port)
	c, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	done := make(chan bool)
	fmt.Println("connected")

	go (func() {
		for {
			message, err := bufio.NewReader(c).ReadString('\n')
			if errors.Is(err, net.ErrClosed) || err == io.EOF {
				c.Close()
				fmt.Println("connection closed")
				return
			} else if err != nil {
				fmt.Println("Error:", err)
				done <- true
				return
			}
			fmt.Print("-> " + message)
		}
	})()

	<-done
	return nil
}
