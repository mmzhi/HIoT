package core

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/ruixiaoedu/hiot/logger"
	"go.uber.org/zap"
)

// OnMessagePublish 处理接收的信息
func (m *Core) OnMessagePublish(clientID, topic string, data []byte) {

	// 处理内部信息
	m.HandleMessage(clientID, topic, data)
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

	return nil, nil
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
