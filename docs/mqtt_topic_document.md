
## 系统Topic

### 配置相关topic
1、获取设备配置信息（异步）
设备发布               sys/{productId}/{deviceId}/config/get
设备订阅               sys/{productId}/{deviceId}/config/get/reply

### 子设备相关topic
1、获取子设备列表（异步
设备发布               sys/{productId}/{deviceId}/subdevice/list
设备订阅               sys/{productId}/{deviceId}/subdevice/list/reply
2、子设备上线（异步）
设备发布               sys/{productId}/{deviceId}/subdevice/login
设备订阅               sys/{productId}/{deviceId}/subdevice/login/reply
3、子设备下线（异步）				
设备发布               sys/{productId}/{deviceId}/subdevice/logout
设备订阅               sys/{productId}/{deviceId}/subdevice/logout/reply
4、获取子设备配置信息（异步）		
设备发布               sys/{productId}/{deviceId}/subdevice/get_config
设备订阅               sys/{productId}/{deviceId}/subdevice/get_config/reply

## 自定义Topic
设备订阅或者设备发布     usr/{productId}/{deviceId}/{topics...}

## 自定义同步Topic
（系统topic实际已经是同步调用）
设备订阅               rpc/{productId}/{deviceId}/{messageId}/req/{topics...}
设备发布               rpc/{productId}/{deviceId}/{messageId}/resp/{topics...}