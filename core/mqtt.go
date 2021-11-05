package core

import (
	"github.com/ruixiaoedu/hiot/adapter"
	"github.com/ruixiaoedu/hiot/core/broker"
	log "github.com/ruixiaoedu/hiot/logger"
	"go.uber.org/zap"
)

// Core 用于处理broker与物联网之间的业务
type Core struct {
	broker *broker.Broker

	engine adapter.Engine
	router *router
}

// NewCore 创建一个mqtt服务
func NewCore(engine adapter.Engine) (*Core, error) {
	core := Core{
		engine: engine,
	}

	var err error
	if core.broker, err = broker.NewBroker(engine.Config(), &core); err != nil {
		log.Error("new broker error", zap.Error(err))
		return nil, err
	}
	core.router = newRouter(&core)

	return &core, nil
}

func (m *Core) Run() {
	m.broker.Run()
}
