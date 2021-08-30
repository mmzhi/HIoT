package impl

import "github.com/fhmq/hmq/plugins/extend"

type authAdapter struct{}

func (adapter *authAdapter) OnClientAuthenticate(clientID, username, password string) bool {
	return true
}

func (adapter *authAdapter) OnClientCheckAcl(clientID, username, topic string, accessType extend.AccessType) bool {
	return true
}
