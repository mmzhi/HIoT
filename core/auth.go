package core

import (
	"github.com/ruixiaoedu/hiot/core/broker"
	"github.com/ruixiaoedu/hiot/model"
	"strings"
)

// OnClientAuthenticate 检测连接是否授权
func (m *Core) OnClientAuthenticate(clientID, username, password string) bool {
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
	deviceDo, err := m.engine.DB().Device().Get(productId, deviceId)
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

// OnClientCheckAcl 检测连接是否授权
// 1、解析topic中的productId和deviceId
func (m *Core) OnClientCheckAcl(clientID, username, topic string, action broker.AccessType) bool {

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
	deviceDo, err := m.engine.DB().Device().Get(productId, deviceId)
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

	// 解析出topic的ProductId和DeviceId
	t := Topic(topic)
	topicType, topicProductId, topicDeviceId, _ := t.Parse()
	if topicType == TopicNoneType {
		// 未知的topic类型
		return false
	}

	if topicProductId != productId || topicDeviceId != deviceId {
		// 不一致，如果不为网关类型则返回false
		if deviceDo.ProductType != model.GatewayType {
			return false
		}

		subDeviceDo, err := m.engine.DB().Device().GetSubdevice(productId, deviceId, topicProductId, topicDeviceId)
		if err != nil {
			return false
		} else if subDeviceDo.State != model.OnlineState {
			// 子设备未上线，无法授权
			return false
		}
	}

	// TODO 判断topics是否符合

	return true
}
