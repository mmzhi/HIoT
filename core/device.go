package core

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/fhmq/hmq/database"
	"github.com/fhmq/hmq/model"
	"regexp"
	"strconv"
	"strings"
)

type Device struct {
	model.Device
}

// valid 验证密码是否正确
func (d *Device) valid(password string) bool {

	if d.ProductId == "" || d.DeviceId == "" || d.DeviceSecret == "" {
		// 密钥不存在
		return false
	}

	var passwords = strings.Split(password, "|")
	if len(passwords) != 4 {
		// 假如不为 ${nonce}|${timestamp}|${signMethod} 格式，返回错误
		return false
	}
	var nonce, _timestamp, signMethod, pwd = passwords[0], passwords[1], passwords[2], passwords[3]

	if match, err := regexp.MatchString(`^[A-Za-z0-9]{4,8}$`, nonce); err != nil {
		return false
	} else if match == false {
		return false
	}

	timestamp, err := strconv.ParseInt(_timestamp, 10, 64)
	if err != nil {
		return false
	}

	// 校验密码
	var p = fmt.Sprintf("clientid=%s\nnonce=%s\ntimestamp=%d",
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

// deviceController 处理设备相关消息业务
type deviceController struct {
	*mqtt
}

// getConfig 获取设备配置
func (m *deviceController) getConfig(message RequestMessage) ResponseMessage {
	// 获取目标 productId 和 deviceId
	productId, deviceId, err := parseSysTopic(message.Topic())
	if err != nil {
		return nil // 不作处理
	}
	payload, err := NewRequestPayload(message.Payload())
	if err != nil {
		return nil // 不作处理
	}

	config, err := database.Database.Device().GetConfig(productId, deviceId)
	if err != nil {
		return nil // TODO 暂时不作处理
	}

	return NewQos0ResponseMessage(payload.Success(struct {
		Config *string `json:"config"`
	}{
		Config: config,
	}).Payload())
}
