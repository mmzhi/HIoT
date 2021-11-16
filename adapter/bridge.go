package adapter

type Bridge interface {

	// Push 设备发出的消息
	Push(topic string, data []byte) error
}
