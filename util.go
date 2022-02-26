package remotelog

import (
	"net"
)

func removeServer(servers []Server, targetServer Server) []Server {
	index := -1
	for i, s := range servers {
		if targetServer.getId() == s.getId() {
			index = i
			break
		}
	}
	if index >= 0 {
		return append(servers[:index], servers[index+1:]...)
	}
	return servers
}

func removeUDPAddr(list []*net.UDPAddr, targetAddr *net.UDPAddr) []*net.UDPAddr {
	index := -1
	for i, s := range list {
		if targetAddr == s {
			index = i
			break
		}
	}
	if index >= 0 {
		return append(list[:index], list[index+1:]...)
	}
	return list
}
