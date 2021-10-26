package manage

type Config struct {
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
func NewManage(config *Config) (IManage, error) {
	return &Engine{
		config: config,
	}, nil
}
