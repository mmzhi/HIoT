package impl

import (
	"github.com/fhmq/hmq/adapter"
	"github.com/fhmq/hmq/database"
)

func init() {
	err := adapter.Register(&builder{})
	if err != nil {
		return
	}
}

// builder 数据库创建生成器
type builder struct{}

// Build 创建扩展访问对象
func (b *builder) Build(database database.IDatabase) (adapter.IAdapter, error) {
	return struct {
		adapter.IAuthAdapter
		adapter.IConnectAdapter
		adapter.IMessageAdapter
	}{
		IAuthAdapter: &authAdapter{
			Database: database,
		},
		IConnectAdapter: &connectAdapter{},
		IMessageAdapter: &messageAdapter{},
	}, nil
}
