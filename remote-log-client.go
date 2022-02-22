package remotelog

type remoteLogClient interface {
	connect(host string, port int) error
}

func Connect(host string, port int) error {
	return client.connect(host, port)
}
