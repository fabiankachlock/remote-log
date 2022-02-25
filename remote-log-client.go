package remotelog

func ConnectTcp(host string, port int) error {
	return tcpClient{}.connect(host, port)
}
