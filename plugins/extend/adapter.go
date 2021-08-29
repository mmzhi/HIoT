package extend

import "github.com/fhmq/hmq/plugins/database"

// NewAdapter 创建一个新的适配器
func NewAdapter(database database.Database) IAdapter {
	return struct {
		IAuthAdapter
		IConnectAdapter
		IMessageAdapter
	}{
		IAuthAdapter:    &authAdapter{},
		IConnectAdapter: &connectAdapter{},
		IMessageAdapter: &messageAdapter{},
	}
}
