package remotelog

var (
	resultsChan   = make(chan error, 1)
	doneChan      = make(chan bool, 1)
	activeServers = []Server{}
)

type Server interface {
	Listen(host string, port int) error
	Close() error
	write(p []byte)
	closedChan() <-chan bool
}

func NewTcp() Server {
	server := newTcpServer()
	activeServers = append(activeServers, server)

	go func() {
		// remove server when closed
		<-server.closedChan()
		activeServers = removeServer(activeServers, server)
	}()

	return server
}
