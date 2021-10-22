package common

import (
	"github.com/fhmq/hmq/database"
	"gorm.io/gorm"
)

type _device struct {
	orm *gorm.DB
}

func (db *_device) Add(device *database.Device) error {
	if tx := db.orm.Create(device); tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (db *_device) Get(productId string, deviceId string) (*database.Device, error) {
	var device database.Device
	if tx := db.orm.Where("product_id = ? AND device_id = ?", productId, deviceId).First(&device); tx.Error != nil {
		return nil, tx.Error
	}
	return &device, nil
}

func (db *_device) GetSubdevice(productId string, deviceId string, subProductId string, subDeviceId string) (*database.Device, error) {
	var device database.Device
	if tx := db.orm.Where("gateway_product_id = ? AND gateway_device_id = ? AND product_id = ? AND device_id = ? AND product_type = 3", productId, deviceId, subProductId, subDeviceId).First(&device); tx.Error != nil {
		return nil, tx.Error
	}
	return &device, nil
}

func (db *_device) List(page database.Page) ([]database.Device, database.Page, error) {
	var devices []database.Device
	//if tx := db.orm.Offset(page - 1*limit).Limit(limit).Find(&devices); tx.Error != nil {
	//	return nil, tx.Error
	//}
	return devices, page, nil
}

func (db *_device) Update(device *database.Device) error {
	if tx := db.orm.Save(device); tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (db *_device) Delete(productId string, deviceId string) error {
	if tx := db.orm.Delete(&database.Device{
		ProductId: productId,
		DeviceId:  deviceId,
	}); tx.Error != nil {
		return tx.Error
	}
	return nil
}
