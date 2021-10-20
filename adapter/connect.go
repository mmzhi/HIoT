package adapter

type connectAdapter struct {
}

func (adapter *connectAdapter) OnClientConnected(clientID, username, ipaddress string) {
	return
}

func (adapter *connectAdapter) OnClientDisconnected(clientID, username string) {
	return
}
