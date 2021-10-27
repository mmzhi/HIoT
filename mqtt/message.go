package mqtt

import (
	"encoding/json"
)

type RequestPayload struct {
	Id string `json:"id"` // 消息ID
}

// NewRequestPayload 创建常规的
func NewRequestPayload(payload []byte) (*RequestPayload, error) {
	var requestPayload RequestPayload
	err := json.Unmarshal(payload, &requestPayload)
	return &requestPayload, err
}

// Success 返回成功信息
func (payload *RequestPayload) Success(data interface{}) *ResponsePayload {
	return &ResponsePayload{
		Id:      payload.Id,
		Code:    "0",
		Message: "OK",
		Data:    data,
	}
}

// Fail 返回失败信息
func (payload *RequestPayload) Fail(code string, message string, datas ...interface{}) *ResponsePayload {
	var data interface{} = nil
	if len(datas) > 0 {
		data = datas[0]
	}
	return &ResponsePayload{
		Id:      payload.Id,
		Code:    code,
		Message: message,
		Data:    data,
	}
}

type ResponsePayload struct {
	Id      string      `json:"id"`             // 消息ID
	Code    string      `json:"code"`           // 相应码
	Message string      `json:"message"`        // 相应描述
	Data    interface{} `json:"data,omitempty"` // 消息主体
}

// Payload 返回payload
func (payload *ResponsePayload) Payload() []byte {
	p, err := json.Marshal(payload)
	if err != nil {
		// TODO 日后增加捕捉异常再启动
		//logger.Panic("ResponsePayload Marshal Fail", zap.Error(err))
	}
	return p
}

// RequestMessage 请求消息接口
type RequestMessage interface {
	ClientId() string
	Topic() string
	Payload() []byte
}

// NewRequestMessage 创建请求消息
func NewRequestMessage(clientID string, topic string, payload []byte) RequestMessage {
	return &requestMessage{
		clientID: clientID,
		topic:    topic,
		payload:  payload,
	}
}

type requestMessage struct {
	clientID string
	topic    string
	payload  []byte
}

func (m *requestMessage) ClientId() string {
	return m.clientID
}

func (m *requestMessage) Topic() string {
	return m.topic
}

func (m *requestMessage) Payload() []byte {
	return m.payload
}

// ResponseMessage 应答消息接口
type ResponseMessage interface {
	Payload() []byte
	Qos() byte
}

// NewQos0ResponseMessage 创建Qos0应答消息
func NewQos0ResponseMessage(payload []byte) ResponseMessage {
	return &responseMessage{
		qos:     0,
		payload: payload,
	}
}

// NewQos1ResponseMessage 创建Qos1应答消息
func NewQos1ResponseMessage(payload []byte) ResponseMessage {
	return &responseMessage{
		qos:     1,
		payload: payload,
	}
}

// responseMessage 回复消息
type responseMessage struct {
	qos     byte
	payload []byte
}

func (m *responseMessage) Qos() byte {
	return m.qos
}

func (m *responseMessage) Payload() []byte {
	return m.payload
}
