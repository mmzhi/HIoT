package mqtt

import (
	"github.com/fhmq/hmq/config"
	log "github.com/fhmq/hmq/logger"
	"github.com/fhmq/hmq/mqtt/broker"
	"go.uber.org/zap"
)

// mqtt 用于处理broker与物联网之间的业务
type mqtt struct {
	broker *broker.Broker
}

type MQTT interface {
	// Start 启动客户端
	Start()
}

// NewMqtt 创建一个mqtt服务
func NewMqtt(cfg *config.Config) (MQTT, error) {
	var (
		err error
		m   mqtt
	)

	if m.broker, err = broker.NewBroker(&m, cfg); err != nil {
		log.Error("new broker error", zap.Error(err))
		return nil, err
	}

	return &m, nil
}

func (m *mqtt) Start() {
	m.broker.Start()
}
