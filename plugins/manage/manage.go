package manage

import "github.com/ruixiaoedu/hiot/config"

type _config struct {
	Port     int    // 端口
	Username string // 用户名
	Password string // 密码
}

// IManage HTTP接口管理
type IManage interface {
	// Run 运行
	Run()
}

// NewManage 新建适配器
func NewManage() (IManage, error) {
	return &Engine{
		config: &_config{
			Port:     config.Config.Manage.Port,
			Username: config.Config.Manage.Username,
			Password: config.Config.Manage.Password,
		},
	}, nil
}
