package sqlite

import (
	"github.com/fhmq/hmq/database"
	"github.com/fhmq/hmq/database/common"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	err := database.Register(string(database.SQLiteType), &builder{})
	if err != nil {
		return
	}
}

// builder 通用数据库创建生成器
type builder struct{}

// Build 创建一个数据库对象
func (b *builder) Build(dsn string, extend string) (database.IDatabase, error) {
	orm, err := gorm.Open(sqlite.Open("hiot.db"), &gorm.Config{
		Logger: database.NewGormLogger(),
	})
	if err != nil {
		return nil, err
	}
	return common.CreateDatabase(orm)
}
