package broker

import (
	"github.com/fhmq/hmq/adapter"
	"strings"
)

const (
	SUB = adapter.AccessSubscribe // 订阅
	PUB = adapter.AccessPublish   // 发布
)

func (b *Broker) CheckTopicAuth(clientID, username, topic string, action adapter.AccessType) bool {
	if b.adapter != nil {
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

		return b.adapter.OnClientCheckAcl(clientID, username, topic, action)
	}

	return true

}

func (b *Broker) CheckConnectAuth(clientID, username, password string) bool {
	if b.adapter != nil {
		return b.adapter.OnClientAuthenticate(clientID, username, password)
	}

	return true

}
