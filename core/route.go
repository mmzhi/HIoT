package core

// 该文件实现类似mqtt.route中的功能

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/ruixiaoedu/hiot/logger"
	"go.uber.org/zap"
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

// MessageHandler 消息处理
type MessageHandler func(message RequestMessage) ResponseMessage

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
	*Core
	sync.RWMutex
	routes []route
}

// newRouter 初始化路由
func newRouter(m *Core) *router {
	deviceCtl := deviceController{m}
	subdeviceCtl := subdeviceController{m}
	router := &router{
		Core: m,
		routes: []route{
			{"sys/+/+/config/get", deviceCtl.getConfig},

			{"sys/+/+/subdevice/list", subdeviceCtl.getList},
			{"sys/+/+/subdevice/login", subdeviceCtl.login},
			{"sys/+/+/subdevice/logout", subdeviceCtl.logout},
			{"sys/+/+/subdevice/config/get", subdeviceCtl.getConfig},
		},
	}
	return router
}

// HandleMessage 处理MQTT的消息
func (r *router) HandleMessage(clientID, topic string, data []byte) {
	r.RLock()
	defer r.RUnlock()
	for _, e := range r.routes {
		if e.match(topic) {
			logger.Debug("match topic", zap.String("route", e.topic), zap.String("topic", topic))
			go func() {
				m := e.callback(&requestMessage{
					clientID: clientID,
					topic:    topic,
					payload:  data,
				})
				if m != nil {
					// 创建返回的消息包
					packet := packets.NewControlPacket(packets.Publish).(*packets.PublishPacket)
					packet.TopicName = topic + "/reply"
					packet.Qos = m.Qos()
					packet.Payload = m.Payload()

					// 发送消息
					r.broker.PublishMessage(packet)
				}
			}()
			return
		}
	}
}
