package repository

import (
	"github.com/ruixiaoedu/hiot/logger"
	"github.com/ruixiaoedu/hiot/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"math"
	"time"
)

// IDevice 设备数据库接口
type IDevice interface {

	// Add 添加设备
	Add(device *model.Device) error

	// Get 获取 Device
	Get(productId string, deviceId string) (*model.Device, error)

	// GetSubdevice 获取指定网关对象的子设备
	GetSubdevice(productId string, deviceId string, subProductId string, subDeviceId string) (*model.Device, error)

	// GetConfig 获取设备配置
	GetConfig(productId string, deviceId string) (*string, error)

	// List 获取 Device 列表
	List(page model.Page, device *model.Device) ([]model.Device, model.Page, error)

	// ListEnableSubdevice 获取可用的子设备列表（状态为未激活、上线、离线）
	ListEnableSubdevice(productId string, deviceId string) ([]model.Device, error)

	// Update 更新 Device
	Update(device *model.Device) error

	// UpdateState 更新设备状态
	UpdateState(productId string, deviceId string, state model.DeviceState) error

	// Online 设备上线
	Online(device *model.Device, ipaddress string) error

	// Offline 设备下线
	Offline(device *model.Device) error

	// UpdateConfig 更新设备配置
	UpdateConfig(productId string, deviceId string, config *string) error

	// UpdateGateway 更新网关
	UpdateGateway(productId string, deviceId string, gatewayProductId *string, gatewayDeviceId *string) error

	// UpdateSecret 更新密钥
	UpdateSecret(productId string, deviceId string, deviceSecret string) error

	// Delete 删除指定ID设备
	Delete(productId string, deviceId string) error
}

func NewDevice(orm *gorm.DB) IDevice {
	return &_device{orm}
}

type _device struct {
	*gorm.DB
}

func (db *_device) Add(device *model.Device) error {
	if tx := db.Create(device); tx.Error != nil {
		return Error(tx.Error)
	}
	return nil
}

func (db *_device) Get(productId string, deviceId string) (*model.Device, error) {
	var device model.Device
	if tx := db.Where("product_id = ? AND device_id = ?", productId, deviceId).First(&device); tx.Error != nil {
		return nil, Error(tx.Error)
	}
	return &device, nil
}

func (db *_device) GetSubdevice(productId string, deviceId string, subProductId string, subDeviceId string) (*model.Device, error) {
	var device model.Device
	if tx := db.Where("gateway_product_id = ? AND gateway_device_id = ? AND product_id = ? AND device_id = ? AND product_type = 3", productId, deviceId, subProductId, subDeviceId).First(&device); tx.Error != nil {
		return nil, Error(tx.Error)
	}
	return &device, nil
}

// GetConfig 获取设备配置
func (db *_device) GetConfig(productId string, deviceId string) (*string, error) {
	var device model.Device
	if tx := db.Select("config").Where("product_id = ? AND device_id = ?", productId, deviceId).First(&device); tx.Error != nil {
		return nil, Error(tx.Error)
	}
	return device.Config, nil
}

func (db *_device) List(page model.Page, device *model.Device) ([]model.Device, model.Page, error) {
	var devices []model.Device

	sql := db.Model(&model.Device{})
	if device.ProductId != "" {
		sql = sql.Where("product_id = ?", device.ProductId)
	}
	if tx := sql.Scopes(Paginate(&page)).Find(&devices); tx.Error != nil {
		return nil, page, Error(tx.Error)
	}
	var total int64
	if tx := db.Model(&model.Device{}).Count(&total); tx.Error != nil {
		return nil, page, Error(tx.Error)
	}
	page.Total = int(total)
	page.Pages = int(math.Ceil(float64(page.Total) / float64(page.Size)))
	return devices, page, nil
}

// ListEnableSubdevice 获取可用的子设备列表（状态为未激活、上线、离线）
func (db *_device) ListEnableSubdevice(productId string, deviceId string) ([]model.Device, error) {
	var subdevices []model.Device
	if tx := db.Model(&model.Device{}).Where(map[string]interface{}{
		"gateway_product_id": productId,
		"gateway_device_id":  deviceId,

		"state": []model.DeviceState{model.InactiveState, model.OnlineState, model.OfflineState},
	}).Find(&subdevices); tx.Error != nil {
		return nil, tx.Error
	}
	return subdevices, nil
}

func (db *_device) Update(device *model.Device) error {
	if tx := db.Model(device).Select("device_name").Updates(device); tx.Error != nil {
		return Error(tx.Error)
	}
	return nil
}

// UpdateState 更新 Device 状态
func (db *_device) UpdateState(productId string, deviceId string, state model.DeviceState) error {
	if tx := db.Model(&model.Device{
		ProductId: productId,
		DeviceId:  deviceId,
	}).Select("state").Updates(&model.Device{
		State: state,
	}); tx.Error != nil {
		return Error(tx.Error)
	}
	return nil
}

// Online 设备上线
func (db *_device) Online(device *model.Device, ipaddress string) error {
	if tx := db.Model(model.Device{
		ProductId: device.ProductId,
		DeviceId:  device.DeviceId,
	}).Updates(map[string]interface{}{
		"ip_address":  ipaddress,
		"state":       model.OnlineState,
		"online_time": time.Now(),
	}); tx.Error != nil {
		return Error(tx.Error)
	}
	return nil
}

// Offline 设备下线
func (db *_device) Offline(device *model.Device) error {
	device, err := db.Get(device.ProductId, device.DeviceId)
	if err != nil {
		return nil
	}

	// 下线时间
	offlineTime := time.Now()

	// 设备类型为网关，先下线其关联的在线的子设备
	if device.ProductType == model.GatewayType {
		if tx := db.DB.Model(&model.Device{}).Where(map[string]interface{}{
			"gateway_product_id": device.ProductId,
			"gateway_device_id":  device.DeviceId,
			"state":              model.OnlineState,
		}).Updates(map[string]interface{}{
			"state":        model.OfflineState,
			"offline_time": offlineTime,
		}); tx.Error != nil {
			logger.Error("subdevice offline fail", zap.Error(tx.Error))
		}
	}

	// 下线时间
	values := map[string]interface{}{
		"offline_time": offlineTime,
	}

	// 假设设备没有被禁用，便将状态设置为下线
	if device.State != model.DisabledState && device.State != model.InactiveDisabledState {
		values["state"] = model.OfflineState
	}

	if tx := db.DB.Model(&model.Device{
		ProductId: device.ProductId,
		DeviceId:  device.DeviceId,
	}).Updates(values); tx.Error != nil {
		return tx.Error
	}

	return nil
}

// UpdateConfig 更新设备配置
func (db *_device) UpdateConfig(productId string, deviceId string, config *string) error {
	if tx := db.Model(&model.Device{
		ProductId: productId,
		DeviceId:  deviceId,
	}).Select("config").Updates(&model.Device{
		Config: config,
	}); tx.Error != nil {
		return Error(tx.Error)
	}
	return nil
}

// UpdateGateway 更新网关
func (db *_device) UpdateGateway(productId string, deviceId string, gatewayProductId *string, gatewayDeviceId *string) error {
	if tx := db.Model(&model.Device{
		ProductId: productId,
		DeviceId:  deviceId,
	}).Select("gateway_product_id", "gateway_device_id").Updates(&model.Device{
		GatewayProductId: gatewayProductId,
		GatewayDeviceId:  gatewayDeviceId,
	}); tx.Error != nil {
		return Error(tx.Error)
	}
	return nil
}

// UpdateSecret 更新密钥
func (db *_device) UpdateSecret(productId string, deviceId string, deviceSecret string) error {
	if tx := db.Model(&model.Device{
		ProductId: productId,
		DeviceId:  deviceId,
	}).Select("device_secret").Updates(&model.Device{
		DeviceSecret: deviceSecret,
	}); tx.Error != nil {
		return Error(tx.Error)
	}
	return nil
}

// Delete 删除指定ID设备
func (db *_device) Delete(productId string, deviceId string) error {
	if tx := db.Model(&model.Device{
		ProductId: productId,
		DeviceId:  deviceId,
	}).Delete(&model.Device{}); tx.Error != nil {
		return Error(tx.Error)
	}
	return nil
}

// DeleteByProductId 删除指定产品所有设备
func (db *_device) DeleteByProductId(productId string) error {
	if tx := db.Where("product_id = ?", productId).Delete(&model.Device{}); tx.Error != nil {
		return Error(tx.Error)
	}
	return nil
}
