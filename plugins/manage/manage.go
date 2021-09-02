package manage

import (
	"errors"
	"github.com/fhmq/hmq/database"
)

// IManage HTTP接口管理
type IManage interface {
	// Run 运行
	Run()
}

var provider IBuilder

// IBuilder 构建器
type IBuilder interface {
	Build(database database.IDatabase) (IManage, error)
}

// NewManage 新建适配器
func NewManage(database database.IDatabase) (IManage, error) {
	if provider == nil {
		return nil, errors.New("not exists")
	}
	adapter, err := provider.Build(database)
	if err != nil {
		return nil, err
	}
	return adapter, nil
}

// Register manage provider
func Register(builder IBuilder) error {
	if provider != nil {
		return errors.New("already exists")
	}
	provider = builder
	return nil
}

// IHandler broker要实现的接口
type IHandler interface {
	Publish(topic string, data []byte)
}
