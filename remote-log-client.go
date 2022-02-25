package remotelog

type remotelogClient struct{}

func NewClient() remotelogClient {
	return remotelogClient{}
}

func (c remotelogClient) ConnectTcp(options ServerOptions) error {
	return tcpClient{}.connect(options)
}
