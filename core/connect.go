package core

import (
	"github.com/ruixiaoedu/hiot/logger"
	"github.com/ruixiaoedu/hiot/model"
	"github.com/ruixiaoedu/hiot/repository"
	"go.uber.org/zap"
	"strings"
)

// OnClientConnected MQTT连接时通知
func (m *mqtt) OnClientConnected(clientID, ipaddress string) {
	pd := strings.Split(clientID, ":")
	if len(pd) != 2 {
		return
	}

	if err := repository.DB.Device().Online(&model.Device{
		ProductId: pd[0],
		DeviceId:  pd[1],
	}, ipaddress); err != nil {
		logger.Error("client online fail", zap.Error(err))
		return
	}
}

// OnClientDisconnected MQTT断开连接时通知
func (m *mqtt) OnClientDisconnected(clientID string) {
	pd := strings.Split(clientID, ":")
	if len(pd) != 2 {
		return
	}

	// 下线设备
	if err := repository.DB.Device().Offline(&model.Device{
		ProductId: pd[0],
		DeviceId:  pd[1],
	}); err != nil {
		logger.Error("device offline fail", zap.Error(err))
	}
}
