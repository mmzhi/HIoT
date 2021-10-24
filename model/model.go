package model

import (
	"time"
)

// AccessType acl type
type AccessType int

const (
	PubAccessType AccessType = 1 // 发布
	SubAccessType AccessType = 2 // 订阅
)

// ProductType 产品类型
type ProductType int

const (
	DeviceType    ProductType = 1 // 设备
	GatewayType   ProductType = 2 // 网关
	SubDeviceType ProductType = 3 // 子设备
)

// DeviceState 设备状态
type DeviceState int

const (
	InactiveState         DeviceState = 0 // 未激活
	OnlineState           DeviceState = 1 // 在线
	OfflineState          DeviceState = 2 // 离线
	DisabledState         DeviceState = 3 // 禁用
	InactiveDisabledState DeviceState = 4 // 未激活且禁用
)

// Product 产品的表结构声明
type Product struct {
	ProductId string `gorm:"primaryKey"`

	ProductType ProductType
	ProductName string

	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName gorm获取表名
func (*Product) TableName() string {
	return "product"
}

// Device 设备的表结构声明
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

// TableName gorm获取表名
func (*Device) TableName() string {
	return "device"
}

// Page 分页
type Page struct {
	Total   int // 总数
	Size    int // 每页大小
	Current int // 页码
	Pages   int // 总页数
}
