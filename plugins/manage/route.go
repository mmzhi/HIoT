package manage

import (
	"fmt"
	"github.com/fhmq/hmq/adapter"
	"github.com/fhmq/hmq/database"
	"github.com/fhmq/hmq/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HTTP接口一览
//
// 产品管理
// 1、添加产品					POST 	/api/v1/product
// 2、修改产品信息					POST 	/api/v1/product/{productId}
// 3、获取产品列表					GET 	/api/v1/product/{productId}
// 4、获取指定产品信息				GET		/api/v1/product
// 5、删除产品					DELETE	/api/v1/product/{productId}

// 设备管理
// 1、添加设备					POST	/api/v1/device/{productId}
// 2、修改设备信息					POST	/api/v1/device/{productId}/{deviceId}
// 3、启用设备					POST	/api/v1/device/{productId}/{deviceId}/enable
// 4、停用设备					POST	/api/v1/device/{productId}/{deviceId}/disable
// 5、获取设备列表					GET		/api/v1/device
// 6、获取指定设备信息				GET		/api/v1/device/{productId}/{deviceId}
// 7、删除设备					DELETE	/api/v1/device/{productId}/{deviceId}
// 8、修改设备配置					POST	/api/v1/device/{productId}/{deviceId}/config
// 9、重置设备					POST	/api/v1/device/{productId}/{deviceId}/reset
// 10、修改子设备与网关的拓扑关系	POST	/api/v1/device/{productId}/{deviceId}/topology
// 11、删除子设备与网关的拓扑关系	DELETE	/api/v1/device/{productId}/{deviceId}/topology

// 消息通信
// 1、向指定设备发送异步消息		POST	/api/v1/message/publish
// 2、rrpc向设备发送同步消息		POST	/api/v1/message/rrpc

// Engine HTTP对象
type Engine struct {
	*gin.Engine
	config   *Config            // 配置
	database database.IDatabase // 数据库功能
	handler  adapter.IHandler   // broker扩展方法
}

// Run 运行
func (e *Engine) Run() {
	gin.SetMode(gin.ReleaseMode)

	e.Engine = gin.New()
	e.Engine.Use(RecoveryWithLogger())

	NewProductController(e).Run()
	NewDeviceController(e).Run()

	err := e.Engine.Run(fmt.Sprintf("0.0.0.0:%d", e.config.Port))

	if err != nil {
		logger.Fatal("http manage error", zap.Error(err))
	}
}
