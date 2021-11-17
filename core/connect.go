package core

import (
	"encoding/json"
	"github.com/ruixiaoedu/hiot/logger"
	"github.com/ruixiaoedu/hiot/model"
	"github.com/ruixiaoedu/hiot/utils"
	"go.uber.org/zap"
	"strings"
	"time"
)

// ClientStatusTrigger MQTT上下线通知
type ClientStatusTrigger struct {
	State     model.DeviceState `json:"state"`     // 状态
	ProductId string            `json:"productId"` // 产品id
	DeviceId  string            `json:"deviceId"`  // 设备id
	IpAddress *string           `json:"ipAddress"` // 客户端ip
	Time      string            `json:"time"`      // 发生时间（如果是上线，则是上线时间，下线则是下线时间）
}

// OnClientConnected MQTT连接时通知
func (m *Core) OnClientConnected(clientID, ipaddress string) {
	pd := strings.Split(clientID, ":")
	if len(pd) != 2 {
		return
	}

	onlineTime := time.Now()

	if err := m.engine.DB().Device().Online(&model.Device{
		ProductId:  pd[0],
		DeviceId:   pd[1],
		OnlineTime: &onlineTime,
	}, ipaddress); err != nil {
		logger.Error("client online fail", zap.Error(err))
		return
	}

	logger.Debug("client online", zap.String("clientID", clientID))

	// 发送桥接的通知
	clientStatusTrigger := struct {
		RequestPayload
		Data ClientStatusTrigger `json:"data"`
	}{
		RequestPayload: RequestPayload{
			Id: utils.GenUniqueId(),
		},
		Data: ClientStatusTrigger{
			State:     model.OnlineState,
			ProductId: pd[0],
			DeviceId:  pd[1],
			Time:      onlineTime.String(),
			IpAddress: &ipaddress,
		}}
	bs, _ := json.Marshal(&clientStatusTrigger)
	m.engine.Bridge().Push("trg/"+pd[0]+"/"+pd[1]+"/mqtt/state", bs)
}

// OnClientDisconnected MQTT断开连接时通知
func (m *Core) OnClientDisconnected(clientID string) {
	pd := strings.Split(clientID, ":")
	if len(pd) != 2 {
		return
	}

	offlineTime := time.Now()

	// 下线设备
	if err := m.engine.DB().Device().Offline(&model.Device{
		ProductId:   pd[0],
		DeviceId:    pd[1],
		OfflineTime: &offlineTime,
	}); err != nil {
		logger.Error("device offline fail", zap.Error(err))
	}

	logger.Debug("client offline", zap.String("clientID", clientID))

	// 发送桥接的通知
	clientStatusTrigger := struct {
		RequestPayload
		Data ClientStatusTrigger `json:"data"`
	}{
		RequestPayload: RequestPayload{
			Id: utils.GenUniqueId(),
		},
		Data: ClientStatusTrigger{
			State:     model.OfflineState,
			ProductId: pd[0],
			DeviceId:  pd[1],
			Time:      offlineTime.String(),
		}}
	bs, _ := json.Marshal(&clientStatusTrigger)
	m.engine.Bridge().Push("trg/"+pd[0]+"/"+pd[1]+"/mqtt/state", bs)
}
