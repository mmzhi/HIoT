package core

// Kick 踢掉客户端
func (m *mqtt) Kick(productId string, deviceId string) {
	m.broker.Kick(productId + ":" + deviceId)
}
