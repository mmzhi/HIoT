package mqtt

import (
	"github.com/fhmq/hmq/database"
	"github.com/fhmq/hmq/logger"
	"github.com/fhmq/hmq/model"
	"go.uber.org/zap"
	"strings"
	"time"
)

func (m *mqtt) OnClientConnected(clientID, ipaddress string) {
	pd := strings.Split(clientID, ":")
	if len(pd) != 2 {
		return
	}
	// 修改状态为登录，修改上线是见为当前时间
	if tx := database.Database().Orm().Model(model.Device{
		ProductId: pd[0],
		DeviceId:  pd[1],
	}).Updates(map[string]interface{}{
		"ipaddress":  ipaddress,
		"State":      model.OnlineState,
		"OnlineTime": time.Now(),
	}); tx.Error != nil {
		logger.Error("client online fail", zap.Error(tx.Error))
		return
	}
}

func (m *mqtt) OnClientDisconnected(clientID string) {
	pd := strings.Split(clientID, ":")
	if len(pd) != 2 {
		return
	}

	// 下线设备
	m.offline(pd[0], pd[1])
}

// offline 下线设备
func (m *mqtt) offline(productId string, deviceId string) {
	device, err := database.Database().Device().Get(productId, deviceId)
	if err != nil {
		logger.Error("get client fail", zap.Error(err))
	}

	// 设备类型为网关，需同时下线其关联的在线的子设备
	if device.ProductType == model.GatewayType {
		if tx := database.Database().Orm().Model(&model.Device{}).Where(map[string]interface{}{
			"GatewayProductId": productId,
			"GatewayDeviceId":  deviceId,
			"State":            model.OnlineState,
		}).Updates(map[string]interface{}{
			"State":       model.OfflineState,
			"OfflineTime": time.Now(),
		}); tx.Error != nil {
			logger.Error("subdevice offline fail", zap.Error(tx.Error))
		}
	}

	if device.State == model.DisabledState || device.State == model.InactiveDisabledState {
		// 设备为禁用状态，只更新下线时间
		if tx := database.Database().Orm().Model(&model.Device{
			ProductId: productId,
			DeviceId:  deviceId,
		}).Updates(map[string]interface{}{
			"OfflineTime": time.Now(),
		}); tx.Error != nil {
			logger.Error("device offline fail", zap.Error(tx.Error))
		}
	} else {
		if tx := database.Database().Orm().Model(&model.Device{
			ProductId: productId,
			DeviceId:  deviceId,
		}).Updates(map[string]interface{}{
			"State":       model.OfflineState,
			"OfflineTime": time.Now(),
		}); tx.Error != nil {
			logger.Error("device offline fail", zap.Error(tx.Error))
		}
	}
}
