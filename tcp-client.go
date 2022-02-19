package remotelog

import (
	"bufio"
	"fmt"
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

	go (func() {
		for {
			message, err := bufio.NewReader(c).ReadString('\n')
			if err != nil {
				fmt.Println("Error", err)
				done <- true
				return
			}
			fmt.Print("-> " + message)
		}
	})()

	<-done
	return nil
}
