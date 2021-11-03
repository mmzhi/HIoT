package core

// Core core 接口
type Core interface {
	// Start 启动客户端
	Start()

	// Kick 剔除设备
	Kick(productId string, deviceId string)

	// Publish 发布消息
	Publish(productId string, deviceId string, topic string, data []byte) error

	// Rpc 发布同步消息
	Rpc(productId string, deviceId string, topic string, data []byte) ([]byte, error)
}
