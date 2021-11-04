package manage

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
// 8、获取设备配置					GET		/api/v1/device/{productId}/{deviceId}/config
// 9、修改设备配置					POST	/api/v1/device/{productId}/{deviceId}/config
// 10、重置设备					POST	/api/v1/device/{productId}/{deviceId}/reset
// 11、获取子设备与网关的拓扑关系	GET		/api/v1/device/{productId}/{deviceId}/topology
// 12、修改子设备与网关的拓扑关系	POST	/api/v1/device/{productId}/{deviceId}/topology
// 13、解除子设备与网关的拓扑关系	DELETE	/api/v1/device/{productId}/{deviceId}/topology

// 消息通信
// 1、向指定设备发送异步消息		POST	/api/v1/message/publish
// 2、rpc向设备发送同步消息			POST	/api/v1/message/rpc
