package mysql

import (
	"github.com/fhmq/hmq/plugins/database"
	"gorm.io/gorm"
)

type _device struct {
	db *gorm.DB
}

func (db *_device) Add(device database.Device) error {
	return nil
}

func (db *_device) Get(productId string, deviceId string) (database.Device, error) {
	return database.Device{}, nil
}

func (db *_device) List(page int, limit int) ([]database.Device, error) {
	return nil, nil
}

func (db *_device) Update(device database.Device) error {
	return nil
}

func (db *_device) Delete(productId string, deviceId string) error {
	return nil
}
