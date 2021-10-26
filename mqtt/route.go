package mqtt

// 该文件实现类似mqtt.route中的功能

import (
	"strings"
	"sync"
)

//
// 设备topic
// 1、获取设备配置信息			设备调用				sys/{productId}/{deviceId}/config/get
//
// 子设备topic
// 1、获取子设备列表			设备调用				sys/{productId}/{deviceId}/subdevice/list
// 2、子设备上线				设备调用				sys/{productId}/{deviceId}/subdevice/login
// 3、子设备下线				设备调用				sys/{productId}/{deviceId}/subdevice/logout
// 4、获取子设备配置信息		设备调用				sys/{productId}/{deviceId}/subdevice/config/get
//

type Message interface {
	ClientID() string
	Topic() string
	Payload() []byte
}

type message struct {
	clientID string
	topic    string
	payload  []byte
}

func (m *message) ClientID() string {
	return m.clientID
}

func (m *message) Topic() string {
	return m.topic
}

func (m *message) Payload() []byte {
	return m.payload
}

type MessageHandler func(message Message) (Message, error)

// route 路由
type route struct {
	topic    string
	callback MessageHandler
}

// match 是否和匹配
func (r *route) match(topic string) bool {
	return r.topic == topic || routeIncludesTopic(r.topic, topic)
}

func routeIncludesTopic(route, topic string) bool {
	return match(strings.Split(route, "/"), strings.Split(topic, "/"))
}

func match(route []string, topic []string) bool {
	if len(route) == 0 {
		return len(topic) == 0
	}

	if len(topic) == 0 {
		return route[0] == "#"
	}

	if route[0] == "#" {
		return true
	}

	if (route[0] == "+") || (route[0] == topic[0]) {
		return match(route[1:], topic[1:])
	}
	return false
}

// router 路由管理
type router struct {
	sync.RWMutex
	routes []route
}

// Router 初始化路由
func Router() *router {
	router := &router{
		routes: []route{},
	}

	//router.addRoute("sys/+/+/config/get", "")
	//router.addRoute("sys/+/+/config/get", "")

	return router
}

// OnMessagePublish 处理MQTT的消息
func (r *router) OnMessagePublish(clientID, username, topic string, data []byte) {
	r.RLock()
	for _, e := range r.routes {
		if e.match(topic) {
			go func() {
				e.callback(&message{
					clientID: clientID,
					topic:    topic,
					payload:  data,
				})
			}()
		}
	}
	r.RUnlock()
}
