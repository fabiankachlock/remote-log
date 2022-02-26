package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	remotelog "github.com/fabiankachlock/remote-log"
)

func main() {
	args := os.Args[1:]
	if len(args) < 3 {
		fmt.Println("Error: missing arguments (ussage: <\"tcp\"|\"udp\"> <host> <port>)")
		return
	}

	host := args[1]
	port, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println(err)
		return
	}
	client := remotelog.NewClient(remotelog.InstanceOptions{
		EnableLogging: true,
	})

	var c remotelog.ConnectedClient
	if args[0] == "tcp" {
		c, err = client.ConnectTcp(remotelog.ServerOptions{
			Host: host,
			Port: port,
		})
	} else if args[0] == "udp" {
		c, err = client.ConnectUdp(remotelog.ServerOptions{
			Host: host,
			Port: port,
		})
	}

	if err != nil {
		fmt.Println(err)
	}

	<-time.After(time.Second * 7)
	c.Close()
}
