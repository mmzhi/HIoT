package manage

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ruixiaoedu/hiot/adapter"
	"github.com/ruixiaoedu/hiot/logger"
	"go.uber.org/zap"
)

// IManage HTTP接口管理
type IManage interface {
	// Run 运行
	Run()
}

// manage HTTP对象
type manage struct {
	gin    *gin.Engine
	engine adapter.Engine // 引擎
}

// NewManage 新建适配器
func NewManage(engine adapter.Engine) IManage {

	gin.SetMode(gin.ReleaseMode)

	m := manage{
		engine: engine,
		gin:    gin.New(),
	}
	m.gin.Use(RecoveryWithLogger())

	{
		// 第一版本接口
		authorized := m.gin.Group("/api/v1", m.BasicAuth())
		NewProductController(m).Routes(authorized)
		NewDeviceController(m).Routes(authorized)
	}

	return &m
}

// Run 运行
func (e *manage) Run() {
	err := e.gin.Run(fmt.Sprintf("0.0.0.0:%d", e.engine.Config().Manage.Port))
	if err != nil {
		logger.Fatal("http manage error", zap.Error(err))
	}
}
