package repository

import (
	"errors"
	"github.com/fhmq/hmq/model"
	"github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
	"strings"
)

// Error 对异常重新包装
func Error(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.ErrDataNotExist
	} else if sqliteError, ok := err.(sqlite3.Error); ok {
		if strings.HasPrefix(sqliteError.Error(), "UNIQUE constraint failed") {
			return model.ErrDuplicateData
		}
	}
	return model.ErrDatabase
}
