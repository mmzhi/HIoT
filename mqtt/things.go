package mqtt

import (
	"container/list"
	"github.com/eclipse/paho.mqtt.golang/packets"
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

type message struct {
	clientID string
	topic    string
	payload  []byte
}

type MessageHandler func(message message) message

type mqttRoute struct {
	topic    string
	callback MessageHandler
}

func (r *mqttRoute) match(topic string) bool {
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

type mqttRouter struct {
	sync.RWMutex
	routes *list.List
}

func newRouter() *mqttRouter {
	router := &mqttRouter{routes: list.New()}

	//router.addRoute("sys/+/+/config/get", "")
	//router.addRoute("sys/+/+/config/get", "")

	return router
}

func (r *mqttRouter) addRoute(topic string, callback MessageHandler) {
	r.Lock()
	defer r.Unlock()
	for e := r.routes.Front(); e != nil; e = e.Next() {
		if e.Value.(*mqttRoute).topic == topic {
			r := e.Value.(*mqttRoute)
			r.callback = callback
			return
		}
	}
	r.routes.PushBack(&mqttRoute{topic: topic, callback: callback})
}

func (r *mqttRouter) matchAndDispatch(message *packets.PublishPacket) {
	r.RLock()
	for e := r.routes.Front(); e != nil; e = e.Next() {
		if e.Value.(*mqttRoute).match(message.TopicName) {
			//hd := e.Value.(*mqttRoute).callback
			go func() {
				//hd(message)
			}()
		}
	}
	r.RUnlock()
}
