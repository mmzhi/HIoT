package core

import (
	"github.com/ruixiaoedu/hiot/adapter"
	"github.com/ruixiaoedu/hiot/core/broker"
	log "github.com/ruixiaoedu/hiot/logger"
	"go.uber.org/zap"
	"sync"
)

// Core 用于处理broker与物联网之间的业务
type Core struct {
	// MQTT Borker
	broker *broker.Broker

	// 系统引擎
	engine adapter.Engine

	// 路由信息
	routeMutex sync.RWMutex
	routes     []route

	// RPC消息
	rpcLock    sync.Mutex             // 并发锁
	rpcChanMap map[string]chan []byte // 通用数据
}

// NewCore 创建一个mqtt服务
func NewCore(engine adapter.Engine) (*Core, error) {
	core := Core{
		engine:     engine,
		rpcChanMap: make(map[string]chan []byte),
	}

	var err error
	if core.broker, err = broker.NewBroker(engine.Config(), &core); err != nil {
		log.Error("new broker error", zap.Error(err))
		return nil, err
	}

	// 初始化路由
	core.initRouter()

	return &core, nil
}

func (m *Core) Run() {
	m.broker.Run()
}
