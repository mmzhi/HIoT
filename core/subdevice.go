package core

import (
	"encoding/json"
	"github.com/ruixiaoedu/hiot/logger"
	"github.com/ruixiaoedu/hiot/model"
	"github.com/ruixiaoedu/hiot/repository"
	"go.uber.org/zap"
)

// subdeviceController 处理子设备相关消息业务
type subdeviceController struct {
	*mqtt
}

type SubdeviceGetListItemResponse struct {
	ProductId  string `json:"productId"`
	DeviceId   string `json:"deviceId"`
	DeviceName string `json:"deviceName"`
}
type SubdeviceGetListResponse struct {
	List []SubdeviceGetListItemResponse `json:"list"`
}

// getConfig 获取子设备列表
func (m *subdeviceController) getList(message RequestMessage) ResponseMessage {
	// 获取目标 productId 和 deviceId
	productId, deviceId, err := parseSysTopic(message.Topic())
	if err != nil {
		return nil // 不作处理
	}
	payload, err := NewRequestPayload(message.Payload())
	if err != nil {
		return nil // 不作处理
	}

	// 检验是否网关设备，日后在topic订阅时便删除
	device, err := repository.DB.Device().Get(productId, deviceId)
	if err != nil {
		return nil // TODO 暂时不作处理
	} else if device.ProductType != model.GatewayType {
		return nil // TODO 暂时不作处理
	}

	var subdevices []model.Device
	if subdevices, err = repository.DB.Device().ListEnableSubdevice(productId, deviceId); err != nil {
		return nil // TODO 暂时不作处理
	}

	respItems := make([]SubdeviceGetListItemResponse, 0)
	for _, v := range subdevices {
		respItems = append(respItems, SubdeviceGetListItemResponse{
			ProductId:  v.ProductId,
			DeviceId:   v.DeviceId,
			DeviceName: v.DeviceName,
		})
	}

	return NewQos0ResponseMessage(payload.Success(SubdeviceGetListResponse{
		List: respItems,
	}).Payload())
}

// login 子设备上线
func (m *subdeviceController) login(message RequestMessage) ResponseMessage {
	// 获取目标 productId 和 deviceId
	productId, deviceId, err := parseSysTopic(message.Topic())
	if err != nil {
		return nil // 不作处理
	}
	var payload struct {
		RequestPayload
		Data struct {
			ProductId string `json:"productId"`
			DeviceId  string `json:"deviceId"`
		} `json:"data"`
	}
	err = json.Unmarshal(message.Payload(), &payload)
	if err != nil {
		return nil // 不作处理
	}

	// 检验是否网关设备，日后在topic订阅时便删除
	device, err := repository.DB.Device().Get(productId, deviceId)
	if err != nil {
		return nil // TODO 暂时不作处理
	} else if device.ProductType != model.GatewayType {
		return nil // TODO 暂时不作处理
	}

	subdevice, err := repository.DB.Device().
		GetSubdevice(productId, deviceId, payload.Data.ProductId, payload.Data.DeviceId)
	if err != nil {
		return nil // TODO 暂时不作处理
	} else if subdevice.State == model.DisabledState || subdevice.State == model.InactiveDisabledState {
		// 子设备处于禁用状态
		return nil // TODO 暂时不作处理
	}

	if err := repository.DB.Device().Online(subdevice, *device.IpAddress); err != nil {
		logger.Error("subdevice online fail", zap.Error(err))
		return nil // TODO 暂时不作处理
	}

	return NewQos0ResponseMessage(payload.Success(nil).Payload())
}

// logout 子设备下线
func (m *subdeviceController) logout(message RequestMessage) ResponseMessage {
	// 获取目标 productId 和 deviceId
	productId, deviceId, err := parseSysTopic(message.Topic())
	if err != nil {
		return nil // 不作处理
	}
	var payload struct {
		RequestPayload
		Data struct {
			ProductId string `json:"productId"`
			DeviceId  string `json:"deviceId"`
		} `json:"data"`
	}
	err = json.Unmarshal(message.Payload(), &payload)
	if err != nil {
		return nil // 不作处理
	}

	subdevice, err := repository.DB.Device().
		GetSubdevice(productId, deviceId, payload.Data.ProductId, payload.Data.DeviceId)
	if err != nil {
		return nil // TODO 暂时不作处理
	} else if subdevice.State == model.DisabledState || subdevice.State == model.InactiveDisabledState {
		// 子设备处于禁用状态，不可操作
		return nil // TODO 暂时不作处理
	}

	if err := repository.DB.Device().Offline(subdevice); err != nil {
		logger.Error("subdevice offline fail", zap.Error(err))
		return nil // TODO 暂时不作处理
	}

	return NewQos0ResponseMessage(payload.Success(nil).Payload())
}

// getConfig 获取子设备配置
func (m *subdeviceController) getConfig(message RequestMessage) ResponseMessage {
	// 获取目标 productId 和 deviceId
	productId, deviceId, err := parseSysTopic(message.Topic())
	if err != nil {
		return nil // 不作处理
	}
	var payload struct {
		RequestPayload
		Data struct {
			ProductId string `json:"productId"`
			DeviceId  string `json:"deviceId"`
		} `json:"data"`
	}
	err = json.Unmarshal(message.Payload(), &payload)
	if err != nil {
		return nil // 不作处理
	}

	subdevice, err := repository.DB.Device().
		GetSubdevice(productId, deviceId, payload.Data.ProductId, payload.Data.DeviceId)
	if err != nil {
		return nil // TODO 暂时不作处理
	} else if subdevice.State == model.DisabledState || subdevice.State == model.InactiveDisabledState {
		// 子设备处于禁用状态，不可操作
		return nil // TODO 暂时不作处理
	}

	config, err := repository.DB.Device().GetConfig(payload.Data.ProductId, payload.Data.DeviceId)
	if err != nil {
		return nil // TODO 暂时不作处理
	}

	return NewQos0ResponseMessage(payload.Success(struct {
		Config *string `json:"config"`
	}{
		Config: config,
	}).Payload())
}
