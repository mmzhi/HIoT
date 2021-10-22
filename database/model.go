package database

import (
	"gorm.io/gorm"
	"time"
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
	InactiveState DeviceState = 0 // 未激活
	OnlineState   DeviceState = 1 // 在线
	OfflineState  DeviceState = 2 // 离线
	DisabledState DeviceState = 3 // 禁用
)

// Product 产品的表结构声明
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

// IPage 分页接口
type IPage interface {
	GetCurrent() int
	GetSize() int
}

// Paginate 分页方法
func Paginate(page *Page) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page.Current <= 0 {
			page.Current = 1
		}
		switch {
		case page.Size > 10000:
			page.Size = 10000
		case page.Size <= 0:
			page.Size = 10
		}
		offset := (page.Current - 1) * page.Size
		return db.Offset(offset).Limit(page.Size)
	}
}
