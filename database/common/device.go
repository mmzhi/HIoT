package common

import (
	"github.com/fhmq/hmq/database"
	"github.com/fhmq/hmq/model"
	"gorm.io/gorm"
	"math"
)

type _device struct {
	orm *gorm.DB
}

func (db *_device) Add(device *model.Device) error {
	if tx := db.orm.Create(device); tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (db *_device) Get(productId string, deviceId string) (*model.Device, error) {
	var device model.Device
	if tx := db.orm.Where("product_id = ? AND device_id = ?", productId, deviceId).First(&device); tx.Error != nil {
		return nil, tx.Error
	}
	return &device, nil
}

func (db *_device) GetSubdevice(productId string, deviceId string, subProductId string, subDeviceId string) (*model.Device, error) {
	var device model.Device
	if tx := db.orm.Where("gateway_product_id = ? AND gateway_device_id = ? AND product_id = ? AND device_id = ? AND product_type = 3", productId, deviceId, subProductId, subDeviceId).First(&device); tx.Error != nil {
		return nil, tx.Error
	}
	return &device, nil
}

// GetConfig 获取设备配置
func (db *_device) GetConfig(productId string, deviceId string) (*string, error) {
	var device model.Device
	if tx := db.orm.Select("config").Where("product_id = ? AND device_id = ?", productId, deviceId).First(&device); tx.Error != nil {
		return nil, tx.Error
	}
	return device.Config, nil
}

func (db *_device) List(page model.Page, device *model.Device) ([]model.Device, model.Page, error) {
	var devices []model.Device
	if tx := db.orm.Model(&model.Device{}).Where("product_id = ?", device.ProductId).Scopes(database.Paginate(&page)).Find(&devices); tx.Error != nil {
		return nil, page, tx.Error
	}
	var total int64
	if tx := db.orm.Model(&model.Device{}).Count(&total); tx.Error != nil {
		return nil, page, tx.Error
	}
	page.Total = int(total)
	page.Pages = int(math.Ceil(float64(page.Total) / float64(page.Size)))
	return devices, page, nil
}

func (db *_device) Update(device *model.Device) error {
	if tx := db.orm.Model(device).Select("device_name").Updates(device); tx.Error != nil {
		return tx.Error
	}
	return nil
}

// UpdateState 更新 Device 状态
func (db *_device) UpdateState(productId string, deviceId string, state model.DeviceState) error {
	if tx := db.orm.Model(&model.Device{}).Select("state").Updates(&model.Device{
		ProductId: productId,
		DeviceId:  deviceId,
		State:     state,
	}); tx.Error != nil {
		return tx.Error
	}
	return nil
}

// UpdateConfig 更新设备配置
func (db *_device) UpdateConfig(productId string, deviceId string, config *string) error {
	if tx := db.orm.Model(&model.Device{}).Select("config").Updates(&model.Device{
		ProductId: productId,
		DeviceId:  deviceId,
		Config:    config,
	}); tx.Error != nil {
		return tx.Error
	}
	return nil
}

// Delete 删除指定ID设备
func (db *_device) Delete(productId string, deviceId string) error {
	if tx := db.orm.Where("product_id = ? AND device_id = ?", productId, deviceId).Delete(&model.Device{}); tx.Error != nil {
		return tx.Error
	}
	return nil
}

// DeleteByProductId 删除指定产品所有设备
func (db *_device) DeleteByProductId(productId string) error {
	if tx := db.orm.Where("product_id = ?", productId).Delete(&model.Device{}); tx.Error != nil {
		return tx.Error
	}
	return nil
}
