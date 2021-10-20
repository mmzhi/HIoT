package adapter

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/fhmq/hmq/database"
	"regexp"
	"strings"
)

type device struct {
	database.Device
}

// valid 验证密码是否正确
func (d *device) valid(password string) bool {

	if d.ProductId == "" || d.DeviceId == "" || d.DeviceSecret == "" {
		// 密钥不存在
		return false
	}

	var passwords = strings.Split(password, "|")
	if len(passwords) != 4 {
		// 假如不为 ${nonce}|${timestamp}|${signMethod} 格式，返回错误
		return false
	}
	var nonce, timestamp, signMethod, pwd = passwords[0], passwords[1], passwords[2], passwords[3]

	if match, err := regexp.MatchString(`^[A-Za-z0-9]{4,8}$`, nonce); err != nil {
		return false
	} else if match == false {
		return false
	}

	// 校验密码
	var p = fmt.Sprintf("clientid=%s\nnonce=%s\ntimestamp=%s",
		d.ProductId+":"+d.DeviceId, nonce, timestamp)
	switch signMethod {
	case "HmacSHA256":
		h := hmac.New(sha256.New, []byte(d.DeviceSecret))
		h.Write([]byte(p))
		return hex.EncodeToString(h.Sum(nil)) == pwd
	case "HmacSM3":
		// SM3 在日后支持计划
		return false
	default:
		// 未知的签名方式，返回错误
		return false
	}
}
