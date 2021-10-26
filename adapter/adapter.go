package adapter

// 扩展功能

// IMqtt mqtt 接口
type IMqtt interface {
	// Kick 剔除设备
	Kick(productId string, deviceId string)
}
