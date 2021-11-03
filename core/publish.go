package core

func (m *mqtt) OnMessagePublish(clientID, topic string, data []byte) {

	// 处理内部信息
	m.router.HandleMessage(clientID, topic, data)
}

// Publish 发布消息
func (m *mqtt) Publish(productId string, deviceId string, topic string, data []byte) {

}

// Rpc 发布同步消息
func (m *mqtt) Rpc(productId string, deviceId string, topic string, data []byte) ([]byte, error) {
	return nil, nil
}
