package mqtt

func (m *mqtt) OnMessagePublish(clientID, topic string, data []byte) {

	// 处理内部信息
	m.router.HandleMessage(clientID, topic, data)
}

func (m *mqtt) PublishMessage(clientID, topic string, data []byte) {

}
