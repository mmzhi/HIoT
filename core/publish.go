package core

import (
	"errors"
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/ruixiaoedu/hiot/logger"
	"github.com/ruixiaoedu/hiot/utils"
	"go.uber.org/zap"
	"strings"
	"time"
)

// OnMessagePublish 处理接收的信息
func (m *Core) OnMessagePublish(clientID, topic string, data []byte) {
	if strings.HasPrefix(topic, "rpc/") {
		m.RpcReply(topic, data)
	} else if strings.HasPrefix(topic, "rpc/") {

	} else {
		// 处理内部信息
		m.HandleMessage(clientID, topic, data)
	}
}

// Publish 发布消息
func (m *Core) Publish(topic string, qos byte, data []byte) error {

	// 创建返回的消息包
	packet := packets.NewControlPacket(packets.Publish).(*packets.PublishPacket)

	packet.TopicName = topic
	packet.Qos = qos
	packet.Payload = data

	// 发送消息
	m.broker.PublishMessage(packet)

	return nil
}

// Rpc 发布同步消息
// rpc/{messageId}/usr/{productId}/{deviceId}/{topics...}
func (m *Core) Rpc(topic string, qos byte, data []byte) ([]byte, error) {

	id := utils.GenUniqueId()

	// 处理 chan
	ch := make(chan []byte)
	m.rpcLock.Lock()
	m.rpcChanMap[id] = ch
	m.rpcLock.Unlock()
	defer func() {
		m.rpcLock.Lock()
		delete(m.rpcChanMap, id)
		m.rpcLock.Unlock()
	}()

	// 发送信息
	packet := packets.NewControlPacket(packets.Publish).(*packets.PublishPacket)

	packet.TopicName = "rpc/" + id + "/" + topic
	packet.Qos = qos
	packet.Payload = data

	// 发送消息
	m.broker.PublishMessage(packet)

	logger.Infof("消息发送成功，消息ID: %s，topic: %s", id, topic)
	for {
		select {
		// 等待接收列表
		case response := <-ch:
			logger.Infof("接收回调成功，消息ID：%s", id)
			return response, nil
		// 10秒超时时间
		case <-time.After(10 * time.Second):
			logger.Info("消息接收超时")
			return nil, errors.New("timeout")
		}
	}
}

// RpcReply 处理RPC设备返回的消息
func (m *Core) RpcReply(topic string, data []byte) {

	logger.Debugf("收到MQTT回调信息，Topic：%s", topic)

	params := rpcTopicRegexp.FindStringSubmatch(topic)
	if len(params) != 4 {
		return
	}

	m.rpcLock.Lock()
	if ch, ok := m.rpcChanMap[params[1]]; ok {
		ch <- data
	}
	m.rpcLock.Unlock()
}

// HandleMessage 处理MQTT的消息
func (m *Core) HandleMessage(clientID, topic string, data []byte) {
	m.routeMutex.RLock()
	defer m.routeMutex.RUnlock()
	for _, e := range m.routes {
		if e.match(topic) {
			logger.Debug("match topic", zap.String("route", e.topic), zap.String("topic", topic))
			go func() {
				msg := e.callback(&requestMessage{
					clientID: clientID,
					topic:    topic,
					payload:  data,
				})
				if msg != nil {
					// 创建返回的消息包
					packet := packets.NewControlPacket(packets.Publish).(*packets.PublishPacket)
					packet.TopicName = topic + "/reply"
					packet.Qos = msg.Qos()
					packet.Payload = msg.Payload()

					// 发送消息
					m.broker.PublishMessage(packet)
				}
			}()
			return
		}
	}
}
