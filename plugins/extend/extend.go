package extend

// AccessType acl type
type AccessType int

const (
	AccessPublish   AccessType = 1
	AccessSubscribe AccessType = 2
)

// IAuthAdapter 授权适配器接口
type IAuthAdapter interface {
	// OnClientAuthenticate 授权请求
	OnClientAuthenticate(clientID, username, password string) bool
	// OnClientCheckAcl ACL请求
	OnClientCheckAcl(clientID, username, topic string, accessType AccessType) bool
}

// IConnectAdapter 连接适配器接口
type IConnectAdapter interface {
	// OnClientConnected 客户端连接
	OnClientConnected(clientID, username, ipaddress string)
	// OnClientDisconnected 客户端断开连接
	OnClientDisconnected(clientID, username string)
}

// IMessageAdapter 消息适配器接口
type IMessageAdapter interface {
	// OnClientSubscribe 客户端订阅Topic
	OnClientSubscribe(clientID, username, topic string)
	// OnClientUnsubscribe 客户端取消订阅Topic
	OnClientUnsubscribe(clientID, username, topic string)
	// OnMessagePublish 客户端发布消息
	OnMessagePublish(clientID, username, topic string, data []byte)
}

type IAdapter interface {
	IConnectAdapter
	IAuthAdapter
	IMessageAdapter
}

type IHandler interface {
	Publish(topic string, data []byte)
}
