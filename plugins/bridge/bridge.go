package bridge

import (
	"github.com/ruixiaoedu/hiot/adapter"
	"github.com/ruixiaoedu/hiot/plugins/bridge/rabbitmq"
)

// bridge
type bridge struct {
	bridges []adapter.Bridge
}

func NewBridge(engine adapter.Engine) (adapter.Bridge, error) {
	var bd bridge
	// 初始化 rabbit
	if rabbit, err := rabbitmq.NewClient(engine); err != nil {
		return nil, err
	} else if rabbit != nil {
		bd.bridges = append(bd.bridges, rabbit)
	}
	return &bd, nil
}

// Push 处理设备发布的内容
func (b bridge) Push(topic string, data []byte) error {
	for i := 0; i < len(b.bridges); i++ {
		_ = b.bridges[i].Push(topic, data)
	}
	return nil
}
