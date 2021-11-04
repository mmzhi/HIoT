package core

func (m *Core) OnMessagePublish(clientID, topic string, data []byte) {

	// 处理内部信息
	m.router.HandleMessage(clientID, topic, data)
}

// Publish 发布消息
func (m *Core) Publish(productId string, deviceId string, topic string, data []byte) error {
	return nil
}

// Rpc 发布同步消息
func (m *Core) Rpc(productId string, deviceId string, topic string, data []byte) ([]byte, error) {
	return nil, nil
}
