package utils

import (
	"github.com/google/uuid"
)

// GenUniqueId 获取唯一编码
func GenUniqueId() string {
	id, err := uuid.NewRandom()
	if err != nil {
		return ""
	}
	return id.String()
}
