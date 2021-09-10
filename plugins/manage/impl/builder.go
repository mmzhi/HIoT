package impl

import (
	"github.com/fhmq/hmq/database"
	"github.com/fhmq/hmq/plugins/manage"
)

func init() {
	err := manage.Register(&builder{})
	if err != nil {
		return
	}
}

// builder HTTP管理接口创建生成器
type builder struct{}

// Build 创建扩展访问对象
func (b *builder) Build(database database.IDatabase) (manage.IManage, error) {
	return &Engine{
		database: database,
	}, nil
}
