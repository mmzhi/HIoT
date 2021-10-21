package manage

import (
	"github.com/fhmq/hmq/database"
	"github.com/fhmq/hmq/logger"
	"go.uber.org/zap"
)

type Config struct {
	Port int // 端口
}

// IManage HTTP接口管理
type IManage interface {
	// Run 运行
	Run()
}

// NewManage 新建适配器
func NewManage(config *Config) (IManage, error) {

	db, err := database.Database()
	if err != nil {
		logger.Error("get database error", zap.Error(err))
		return nil, err
	}

	return &Engine{
		config:   config,
		database: db,
	}, nil
}
