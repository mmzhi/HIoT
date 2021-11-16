package main

import (
	"github.com/ruixiaoedu/hiot/adapter"
	"github.com/ruixiaoedu/hiot/config"
	"github.com/ruixiaoedu/hiot/core"
	"github.com/ruixiaoedu/hiot/logger"
	"github.com/ruixiaoedu/hiot/plugins/bridge"
	"github.com/ruixiaoedu/hiot/plugins/manage"
	"github.com/ruixiaoedu/hiot/repository"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"runtime"
)

type engine struct {
	core adapter.Core // 核心连接

	manage adapter.Manage // 管理

	bridge adapter.Bridge // 消息流转

	config *config.Config      // 配置
	db     repository.Database // 数据
}

func (e *engine) Core() adapter.Core {
	return e.core
}

func (e *engine) Manage() adapter.Manage {
	return e.manage
}

func (e *engine) Bridge() adapter.Bridge {
	return e.bridge
}

func (e *engine) Config() *config.Config {
	return e.config
}

func (e *engine) DB() repository.Database {
	return e.db
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	e := engine{}

	// 初始化 配置
	cfg, err := config.NewConfigure()
	if err != nil {
		logger.Fatal("configure broker config error: ", zap.Error(err))
	}
	e.config = cfg

	// 配置日志
	if cfg.Debug {
		logger.ConfigLogger(logger.Config{
			Debug: true,
		})
	}

	// 初始化 数据库
	db, err := repository.NewDatabase(cfg.Database)
	if err != nil {
		logger.Fatal("init repository error", zap.Error(err))
	}
	e.db = db

	// 初始化 MQTT
	c, err := core.NewCore(&e)
	if err != nil {
		logger.Fatal("New MQTT Broker error: ", zap.Error(err))
	}
	e.core = c

	// HTTP管理接口
	{
		m := manage.NewManage(&e)
		e.manage = m
		go m.Run()
	}

	// 桥接设计
	{
		b, err := bridge.NewBridge(&e)
		if err != nil {
			logger.Fatal("New Bridge error: ", zap.Error(err))
		}
		e.bridge = b
	}

	c.Run() // 启动MQTT服务
	s := waitForSignal()
	logger.Infof("signal received, broker closed. %s", s)
}

func waitForSignal() os.Signal {
	signalChan := make(chan os.Signal, 1)
	defer close(signalChan)
	signal.Notify(signalChan, os.Kill, os.Interrupt)
	s := <-signalChan
	signal.Stop(signalChan)
	return s
}
