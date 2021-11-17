
## 系统Topic
格式：sys/{productId}/{deviceId}/{topics...}

### 配置相关topic
1、获取设备配置信息
设备发布               sys/{productId}/{deviceId}/config/get
设备订阅               sys/{productId}/{deviceId}/config/get/reply

### 子设备相关topic
1、获取子设备列表
设备发布               sys/{productId}/{deviceId}/subdevice/list
设备订阅               sys/{productId}/{deviceId}/subdevice/list/reply
2、子设备上线
设备发布               sys/{productId}/{deviceId}/subdevice/login
设备订阅               sys/{productId}/{deviceId}/subdevice/login/reply
3、子设备下线
设备发布               sys/{productId}/{deviceId}/subdevice/logout
设备订阅               sys/{productId}/{deviceId}/subdevice/logout/reply
4、获取子设备配置信息
设备发布               sys/{productId}/{deviceId}/subdevice/get_config
设备订阅               sys/{productId}/{deviceId}/subdevice/get_config/reply

## 自定义Topic
1、异步消息
设备订阅或者设备发布     usr/{productId}/{deviceId}/{topics...}
2、同步消息（与异步消息不同之处是前面增加rpc/{messageId}）
设备订阅               rpc/{messageId}/usr/{productId}/{deviceId}/{topics...}
设备发布               rpc/{messageId}/usr/{productId}/{deviceId}/{topics...}


## 触发Topic（trigger）
1、设备上下线状态       trg/{productId}/{deviceId}/mqtt/state