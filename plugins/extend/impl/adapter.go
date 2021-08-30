package impl

import (
	"github.com/fhmq/hmq/plugins/database"
	"github.com/fhmq/hmq/plugins/extend"
)

func init() {
	err := extend.Register(&builder{})
	if err != nil {
		return
	}
}

// builder 数据库创建生成器
type builder struct{}

// Build 创建一扩展访问对象
func (b *builder) Build(database database.IDatabase) (extend.IAdapter, error) {
	return struct {
		extend.IAuthAdapter
		extend.IConnectAdapter
		extend.IMessageAdapter
	}{
		IAuthAdapter:    &authAdapter{},
		IConnectAdapter: &connectAdapter{},
		IMessageAdapter: &messageAdapter{},
	}, nil
}
