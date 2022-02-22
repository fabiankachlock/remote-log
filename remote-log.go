package remotelog

var (
	messagesChan                  = make(chan []byte, 64)
	resultsChan                   = make(chan error, 1)
	doneChan                      = make(chan bool, 1)
	client        remoteLogClient = tcpClient{}
	activeServers                 = []Server{}
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
		<-server.closedChan()
		index := -1
		for i, s := range activeServers {
			if server == s {
				index = i
				break
			}
		}
		activeServers = append(activeServers[:index], activeServers[index+1:]...)

	}()

	return server
}
