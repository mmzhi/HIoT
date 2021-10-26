package broker

import (
	"strings"
)

// AccessType acl type
type AccessType int

const (
	PUB AccessType = 1 // 发布
	SUB AccessType = 2 // 订阅
)

// CheckTopicAuth 检测连接是否授权
func (b *Broker) CheckTopicAuth(clientID, username, topic string, action AccessType) bool {

	if strings.HasPrefix(topic, "$SYS/broker/connection/clients/") {
		return true
	}

	if strings.HasPrefix(topic, "$share/") && action == SUB {
		substr := groupCompile.FindStringSubmatch(topic)
		if len(substr) != 3 {
			return false
		}
		topic = substr[2]
	}

	// TODO 以上是内置方法，日后再处理

	return b.things.OnClientCheckAcl(clientID, username, topic, action)
}

// CheckConnectAuth 检测连接是否授权
func (b *Broker) CheckConnectAuth(clientID, username, password string) bool {
	return b.things.OnClientAuthenticate(clientID, username, password)
}
