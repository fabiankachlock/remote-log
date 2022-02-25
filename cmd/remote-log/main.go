package main

import (
	"fmt"
	"os"
	"strconv"

	remotelog "github.com/fabiankachlock/remote-log"
)

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Println("Error: missing arguments (ussage: <host> <port>)")
		return
	}

	host := args[0]
	port, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	client := remotelog.NewClient()
	err = client.ConnectTcp(remotelog.ServerOptions{
		Host: host,
		Port: port,
	})
	if err != nil {
		fmt.Println(err)
	}
}
