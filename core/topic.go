package core

import (
	"regexp"
	"strings"
)

var (
	// 系统的topic正则，格式：sys/{productId}/{deviceId}/{topics...}
	sysTopicRegexp = regexp.MustCompile(`^sys/([\d\w-_]{1,32})/([\d\w-_]{1,32})/`)

	// 自定义的topic正则，格式：usr/{productId}/{deviceId}/{topics...}
	usrTopicRegexp = regexp.MustCompile(`^usr/([\d\w-_]{1,32})/([\d\w-_]{1,32})/`)

	// 自定义rpc请求的topic正则，格式：rpc/{messageId}/{productId}/{deviceId}/{topics...}
	// 设备只能订阅，设备订阅其中messageId必须为通配符
	rpcReqTopicRegexp = regexp.MustCompile(`^rpc/+/([\d\w-_]{1,32})/([\d\w-_]{1,32})/`)

	// 自定义rpc请求的topic正则，格式：rpc/{messageId}/{productId}/{deviceId}/{topics...}
	// 设备只能发布
	rpcRspTopicRegexp = regexp.MustCompile(`^rpc/([\d\w]{1,32})/([\d\w-_]{1,32})/([\d\w-_]{1,32})/`)
)

type TopicType byte

var (
	TopicNoneType TopicType = 0 // 非支持的topic
	TopicSysType  TopicType = 1 // 系统的topic
	TopicUsrType  TopicType = 2 // 自定义的topic
	TopicRpcType  TopicType = 3 // rpc的topic
)

type Topic string

// Parse 解析出productId、deviceId和topics
func (t Topic) Parse() (topicType TopicType, productId string, deviceId string, topics string) {

	if topicN := strings.SplitN(string(t), "/", 2); len(topicN) == 2 {
		switch topicN[0] {
		case "sys":
			if topicN := strings.SplitN(topicN[1], "/", 3); len(topicN) == 3 {
				return TopicSysType, topicN[0], topicN[1], topicN[2]
			}
		case "usr":
			if topicN := strings.SplitN(topicN[1], "/", 3); len(topicN) == 3 {
				return TopicUsrType, topicN[0], topicN[1], topicN[2]
			}
		case "rpc":
			if topicN := strings.SplitN(topicN[1], "/", 4); len(topicN) == 4 {
				return TopicRpcType, topicN[1], topicN[2], topicN[3]
			}
		}
	}

	return TopicNoneType, "", "", ""
}
