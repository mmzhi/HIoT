package broker

// Things 提供物联网与broker的接口
type Things interface {
	// OnClientAuthenticate 连接是否授权
	OnClientAuthenticate(clientID, username, password string) bool
	// OnClientCheckAcl 客户端订阅发送行为是否授权
	OnClientCheckAcl(clientID, username, topic string, action AccessType) bool
	// OnClientConnected 客户端连接
	OnClientConnected(clientID, ipaddress string)
	// OnClientDisconnected 客户端断开连接
	OnClientDisconnected(clientID string)
	// OnClientSubscribe 客户端订阅Topic
	OnClientSubscribe(clientID, topic string)
	// OnClientUnsubscribe 客户端取消订阅Topic
	OnClientUnsubscribe(clientID, topic string)
	// OnMessagePublish 客户端发布消息
	OnMessagePublish(clientID, topic string, data []byte)
}
