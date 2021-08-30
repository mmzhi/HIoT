package database

import "errors"

// IDatabase 数据库接口
type IDatabase interface {
	Product() IProduct
	Device() IDevice
}

// IProduct 产品接口
type IProduct interface {

	// Add 添加产品
	Add(product Product) error

	// Get 获取 product
	Get(productId string) (Product, error)

	// List 获取 product 列表
	List(page int, limit int) ([]Product, error)

	// Update 更新 product
	Update(product Product) error

	// Delete 删除指定ID产品
	Delete(productId string) error
}

// IDevice 设备接口
type IDevice interface {

	// Add 添加设备
	Add(device Device) error

	// Get 获取 Device
	Get(productId string, deviceId string) (Device, error)

	// List 获取 Device 列表
	List(page int, limit int) ([]Device, error)

	// Update 更新 Device
	Update(device Device) error

	// Delete 删除指定ID设备
	Delete(productId string, deviceId string) error
}

type Type string

const (
	MySQLType Type = "mysql"
)

// IBuilder 构建器
type IBuilder interface {
	Build(dsn string, extend string) (IDatabase, error)
}

var providers = make(map[string]IBuilder)

// Register database provider
func Register(name string, i IBuilder) error {
	if name == "" || i == nil {
		return errors.New("invalid args")
	}

	if _, dup := providers[name]; dup {
		return errors.New("already exists")
	}

	providers[name] = i

	return nil
}

// UnRegister 反注册
func UnRegister(name string) {
	delete(providers, name)
}

// NewDatabase 新建数据库对象
func NewDatabase(name string, dsn string, extend string) (IDatabase, error) {
	if name == "" {
		name = string(MySQLType)
	}
	builder, ok := providers[name]
	if !ok {
		return nil, errors.New("not exists")
	}
	return builder.Build(dsn, extend)
}
