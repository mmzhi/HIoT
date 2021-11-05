package adapter

import (
	"github.com/ruixiaoedu/hiot/config"
	"github.com/ruixiaoedu/hiot/repository"
)

// Engine 引擎
type Engine interface {

	// Core 核心接口
	Core() Core

	// Manage 管理接口
	Manage() Manage

	// Config 获取配置
	Config() *config.Config

	// DB 获取数据库示例
	DB() repository.Database
}
