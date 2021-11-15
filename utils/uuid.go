package utils

import (
	"github.com/segmentio/ksuid"
)

// GenUniqueId 获取唯一编码
func GenUniqueId() string {
	id, err := ksuid.NewRandom()
	if err != nil {
		return ""
	}
	return id.String()
}
