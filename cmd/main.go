package main

import (
	"github.com/fhmq/hmq/config"
	"github.com/fhmq/hmq/core"
	"github.com/fhmq/hmq/database"
	"github.com/fhmq/hmq/logger"
	"github.com/fhmq/hmq/plugins/manage"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 初始化 配置
	cfg, err := config.Configure()
	if err != nil {
		logger.Fatal("configure broker config error: ", zap.Error(err))
	}

	// 配置日志
	if cfg.Debug {
		logger.ConfigLogger(logger.Config{
			Debug: true,
		})
	}

	// 初始化 数据库
	err = database.InitDatabase(cfg.Database)
	if err != nil {
		logger.Fatal("init database error", zap.Error(err))
	}

	// 初始化 MQTT
	m, err := core.NewCore()
	if err != nil {
		logger.Fatal("New MQTT Broker error: ", zap.Error(err))
	}

	// HTTP管理接口
	{
		m, err := manage.NewManage()
		if err != nil {
			log.Fatal("new manage fail", zap.Error(err))
			return
		}
		go m.Run()
	}

	m.Start() // 启动MQTT服务
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
