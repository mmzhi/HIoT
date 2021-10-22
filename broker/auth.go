package broker

import (
	"fmt"
	"github.com/fhmq/hmq/model"
	"regexp"
	"strings"
)

// AccessType acl type
type AccessType int

const (
	PUB AccessType = 1 // 发布
	SUB AccessType = 2 // 订阅
)

var (
	sysTopicRegexp  = regexp.MustCompile(`^sys/([\d\w-_]{1,32})/([\d\w-_]{1,32})/`)
	userTopicRegexp = regexp.MustCompile(`^user/([\d\w-_]{1,32})/([\d\w-_]{1,32})/`)
)

// CheckTopicAuth 检测连接是否授权
func (b *Broker) CheckTopicAuth(clientID, username, topic string, action AccessType) bool {
	//if b.adapter != nil {
	//
	//
	//	return b.adapter.OnClientCheckAcl(clientID, username, topic, action)
	//}

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

	// TODO 以上是内置方法，日后在处理

	if clientID != username {
		// 不符合clientID和username的命名规则
		return false
	}

	// 解析产品Id和设备Id
	var ids = strings.Split(username, ":")
	if len(ids) != 2 {
		return false
	}

	var productId, deviceId = ids[0], ids[1]
	if productId == "" || deviceId == "" {
		return false
	}

	// 从数据库获取数据
	deviceDo, err := b.database.Device().Get(productId, deviceId)
	if err != nil {
		return false
	}

	// 设备禁用
	if deviceDo.State == model.DisabledState {
		// 设备已被禁用，无法授权
		return false
	}

	if deviceDo.ProductType == model.SubDeviceType {
		// 子设备无法直接授权ACL
		return false
	}

	// 假如是设备或者网关，对于自身的topic处理
	if deviceDo.ProductType == model.DeviceType || deviceDo.ProductType == model.GatewayType {

		// 符合系统topic，返回
		if strings.HasPrefix(topic, fmt.Sprintf("sys/%s/%s/", deviceDo.ProductId, deviceDo.DeviceId)) {
			return true
		}

		// 符合用户自定义，返回true
		if strings.HasPrefix(topic, fmt.Sprintf("user/%s/%s/", deviceDo.ProductId, deviceDo.DeviceId)) {
			return true
		}
	}

	// 对于是网关类型，判断是否符合其子设备的topic
	if deviceDo.ProductType == model.GatewayType {
		var params []string
		// 是否符合指定topic格式
		if params = sysTopicRegexp.FindStringSubmatch(topic); len(params) == 2 {
			// 不处理
		} else if params = userTopicRegexp.FindStringSubmatch(topic); len(params) == 2 {
			// 不处理
		} else {
			// 都不符合，退出处理
			goto out1
		}

		var subProductId, subDeviceId = params[0], params[1]

		subDeviceDo, err := b.database.Device().GetSubdevice(productId, deviceId, subProductId, subDeviceId)
		if err != nil {
			return false
		}

		if subDeviceDo.State == model.DisabledState {
			// 子设备已被禁用，无法授权
			return false
		}

		return true
	}
out1:

	return false
}

// CheckConnectAuth 检测连接是否授权
func (b *Broker) CheckConnectAuth(clientID, username, password string) bool {
	if clientID != username {
		// 不符合clientID和username的命名规则
		return false
	}

	// 解析产品Id和设备Id
	var ids = strings.Split(username, ":")
	if len(ids) != 2 {
		return false
	}

	var productId, deviceId = ids[0], ids[1]
	if productId == "" || deviceId == "" {
		return false
	}

	// 从数据库获取数据
	deviceDo, err := b.database.Device().Get(productId, deviceId)
	if err != nil {
		return false
	}

	// 设备禁用
	if deviceDo.State == model.DisabledState {
		// 设备已被禁用，无法授权
		return false
	}

	if deviceDo.ProductType != model.DeviceType && deviceDo.ProductType != model.GatewayType {
		// 仅允许设备和网关授权
		return false
	}

	deviceBo := Device{
		*deviceDo,
	}

	// 校验密码
	return deviceBo.valid(password)
}
