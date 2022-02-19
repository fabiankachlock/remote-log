package remotelog

import "net"

var (
	messagesChan                 = make(chan []byte, 64)
	resultsChan                  = make(chan error, 1)
	doneChan                     = make(chan bool, 1)
	server       remoteLogServer = tcpServer{map[string]net.Conn{}}
	client       remoteLogClient = tcpClient{}
)

type remoteLogServer interface {
	write(p []byte)
	start(host string, port int) error
}

type remoteLogClient interface {
	connect(host string, port int) error
}

func Start(host string, port int) error {
	go func() {
		for msg := range messagesChan {
			server.write(msg)
		}
	}()
	return server.start(host, port)
}

func Connect(host string, port int) error {
	return client.connect(host, port)
}
