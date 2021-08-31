package impl

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

type manage struct {
}

// Response 应答结构体
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// fail 错误的应答
func fail(code int, message string) *Response {
	if message == "" {
		message = "Error"
	}
	if code == 0 {
		code = -1
	}
	return &Response{
		Code:    code,
		Message: message,
	}
}

// success 成功的应答
func success(data interface{}) *Response {
	return &Response{
		Code:    0,
		Message: "OK",
		Data:    data,
	}
}
