package database

import "time"

type ProductType int

const (
	DeviceType    ProductType = 1
	GatewayType   ProductType = 2
	SubDeviceType ProductType = 3
)

// Product Table of product
type Product struct {

	ProductId		string			`gorm:"primaryKey"`

	ProductType ProductType
	ProductName string

	CreatedAt		time.Time
	UpdatedAt		time.Time

}


// Device Table of device
type Device struct {

	ProductId			string		`gorm:"primaryKey"`
	DeviceId			string		`gorm:"primaryKey"`

	ProductType  ProductType
	DeviceName   string
	DeviceSecret string

	GatewayProductId	*string
	GatewayDeviceId		*string

	FirmwareVersion		*string
	IpAddress			*string
	State				int8

	OnlineTime		*time.Time
	OfflineTime		*time.Time

	Config			*string

	CreatedAt		time.Time
	UpdatedAt		time.Time
}