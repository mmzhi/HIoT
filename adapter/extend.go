package adapter

import (
	"errors"
	"github.com/fhmq/hmq/database"
)

// 扩展功能，该插件为必选

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

var provider IBuilder

// IBuilder 构建器
type IBuilder interface {
	Build(database database.IDatabase) (IAdapter, error)
}

// NewAdapter 新建适配器
func NewAdapter(database database.IDatabase) (IAdapter, error) {
	if provider == nil {
		return nil, errors.New("not exists")
	}
	adapter, err := provider.Build(database)
	if err != nil {
		return nil, err
	}
	return adapter, nil
}

// Register adapter provider
func Register(builder IBuilder) error {
	if provider != nil {
		return errors.New("already exists")
	}
	provider = builder
	return nil
}

// IHandler broker要实现的接口
type IHandler interface {
	Publish(topic string, data []byte)
}
