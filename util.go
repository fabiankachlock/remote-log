package remotelog

func removeServer(servers []Server, targetServer Server) []Server {
	index := -1
	for i, s := range servers {
		if targetServer == s {
			index = i
			break
		}
	}
	return append(servers[:index], servers[index+1:]...)
}
