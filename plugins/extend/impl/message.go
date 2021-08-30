package impl

type messageAdapter struct {
}

func (adapter *messageAdapter) OnClientSubscribe(clientID, username, topic string) {
	return
}

func (adapter *messageAdapter) OnClientUnsubscribe(clientID, username, topic string) {
	return
}

func (adapter *messageAdapter) OnMessagePublish(clientID, username, topic string, data []byte) {
	return
}
