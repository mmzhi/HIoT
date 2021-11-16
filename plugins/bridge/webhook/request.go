package webhook

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type Request struct {
	Event     string `json:"event"`     // 消息事件类型
	EventTime string `json:"eventTime"` // 平台事件生成时间
}

// MessageRequest 消息类型的报文
type MessageRequest struct {
	Request
	Topic     string `json:"topic"`     // 消息的topic
	ProductId string `json:"productId"` // 产品ID
	DeviceId  string `json:"deviceId"`  // 设备ID
	Payload   string `json:"payload"`   // 消息主体，base64编码
}

type WebHook struct {
	httpclient *http.Client
}

// Send 发送消息
func (w *WebHook) Send(request interface{}) error {
	data, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/iot/hiot/hook", bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := w.httpclient.Do(req)

	if err != nil {
		return err
	} else if resp.StatusCode != 200 {
		return errors.New("http request not 200")
	}

	return nil
}
