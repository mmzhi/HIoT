package core

// 该文件实现类似mqtt.route中的功能

import (
	"strings"
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

// newRouter 初始化路由
func (m *Core) initRouter() {

	deviceCtl := deviceController{m}
	subdeviceCtl := subdeviceController{m}

	m.routes = []route{
		{"sys/+/+/config/get", deviceCtl.getConfig},

		{"sys/+/+/subdevice/list", subdeviceCtl.getList},
		{"sys/+/+/subdevice/login", subdeviceCtl.login},
		{"sys/+/+/subdevice/logout", subdeviceCtl.logout},
		{"sys/+/+/subdevice/config/get", subdeviceCtl.getConfig},
	}
}
