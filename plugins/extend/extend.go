package extend

// AccessType acl type
type AccessType int

const (
	AccessPublish  		AccessType = 1
	AccessSubscribe 	AccessType = 2
)

type IAuthAdapter interface {
	OnClientAuthenticate(clientID, username, password string) bool
	OnClientCheckAcl(clientID, username, topic string, accessType AccessType) bool
}

type IConnectAdapter interface {
	OnClientConnected(clientID, username, ipaddress string)
	OnClientDisconnected(clientID, username string)
}

type IMessageAdapter interface {
	OnClientSubscribe(clientID, username, topic string)
	OnClientUnsubscribe(clientID, username, topic string)

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