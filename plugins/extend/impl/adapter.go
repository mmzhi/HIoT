package impl

import (
	"github.com/fhmq/hmq/plugins/database"
	"github.com/fhmq/hmq/plugins/extend"
)

// NewAdapter 创建一个新的适配器
func NewAdapter(database database.IDatabase) extend.IAdapter {
	return struct {
		extend.IAuthAdapter
		extend.IConnectAdapter
		extend.IMessageAdapter
	}{
		IAuthAdapter:    &authAdapter{},
		IConnectAdapter: &connectAdapter{},
		IMessageAdapter: &messageAdapter{},
	}
}
