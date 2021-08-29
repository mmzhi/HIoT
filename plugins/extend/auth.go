package extend

type authAdapter struct{}

func (adapter *authAdapter) OnClientAuthenticate(clientID, username, password string) bool {
	return true
}

func (adapter *authAdapter) OnClientCheckAcl(clientID, username, topic string, accessType AccessType) bool {
	return true
}
