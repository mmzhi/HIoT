package database

import "time"

type ProductType int

const (
	DeviceType    ProductType = 1 // 设备
	GatewayType   ProductType = 2 // 网关
	SubDeviceType ProductType = 3 // 子设备
)

type DeviceState int

const (
	InactiveState DeviceState = 0 // 未激活
	OnlineState   DeviceState = 1 // 在线
	OfflineState  DeviceState = 2 // 离线
	DisabledState DeviceState = 3 // 禁用
)

// Product Table of product
type Product struct {
	ProductId string `gorm:"primaryKey"`

	ProductType ProductType
	ProductName string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (*Product) TableName() string {
	return "product"
}

// Device Table of device
type Device struct {
	ProductId string `gorm:"primaryKey"`
	DeviceId  string `gorm:"primaryKey"`

	ProductType  ProductType
	DeviceName   string
	DeviceSecret string

	GatewayProductId *string
	GatewayDeviceId  *string

	FirmwareVersion *string
	IpAddress       *string
	State           DeviceState

	OnlineTime  *time.Time
	OfflineTime *time.Time

	Config *string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (*Device) TableName() string {
	return "device"
}
