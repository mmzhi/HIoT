package manage

import (
	"github.com/fhmq/hmq/database"
)

// IManage HTTP接口管理
type IManage interface {
	// Run 运行
	Run()
}

// NewManage 新建适配器
func NewManage(database database.IDatabase) (IManage, error) {
	return &Engine{
		database: database,
	}, nil
}