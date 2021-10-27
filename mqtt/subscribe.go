package mqtt

import (
	"github.com/fhmq/hmq/database"
	"github.com/fhmq/hmq/model"
)

func (m *mqtt) OnClientSubscribe(clientID, topic string) {

}

func (m *mqtt) OnClientUnsubscribe(clientID, topic string) {

}

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

	device, err := database.Database().Device().Get(productId, deviceId)
	if err != nil {
		return nil // TODO 暂时不作处理
	} else if device.ProductType != model.GatewayType {
		return nil // TODO 暂时不作处理
	}

	var subdevices []model.Device
	if tx := database.Database().Orm().Model(&model.Device{}).Where(map[string]interface{}{
		"gateway_product_id": productId,
		"gateway_device_id":  deviceId,
	}).Not(map[string]interface{}{ // 提出禁用的子设备
		"state": []model.DeviceState{model.DisabledState, model.InactiveDisabledState},
	}).Find(&subdevices); tx.Error != nil {
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
