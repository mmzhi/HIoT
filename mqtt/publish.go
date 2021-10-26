package mqtt

func (m *mqtt) OnMessagePublish(clientID, topic string, data []byte) {
	m.router.HandleMessage(clientID, topic, data)
}
