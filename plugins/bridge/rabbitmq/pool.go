package rabbitmq

import (
	"github.com/ruixiaoedu/hiot/logger"
	"github.com/streadway/amqp"
	"time"
)

// Pool 连接池
type Pool struct {
	url      string // amqp地址
	option   Options
	channels chan *Channel
}

// Options 选项
type Options struct {
	ConnectionNum int // 连接数
	ChannelNum    int // 每个连接的channel数量
}

// Connection 连接
type Connection struct {
	*amqp.Connection
}

// Channel 渠道
type Channel struct {
	*amqp.Channel
}

const (
	waitTime = 3 * time.Second
)

func NewPool(url string, opts ...Options) (*Pool, error) {
	var opt Options
	if len(opts) > 0 {
		opt = opts[0]
	}

	if opt.ConnectionNum <= 0 {
		opt.ConnectionNum = 5
	}

	if opt.ChannelNum <= 0 {
		opt.ChannelNum = 10
	}

	var pool = Pool{
		url:      url,
		option:   opt,
		channels: make(chan *Channel, opt.ConnectionNum*opt.ChannelNum),
	}

	for i := 0; i < opt.ConnectionNum; i++ {
		connect, err := NewConnection(url)
		if err != nil {
			return nil, err
		}

		for j := 0; j < opt.ChannelNum; j++ {
			ch, err := connect.Channel()
			if err != nil {
				return nil, err
			}
			pool.channels <- ch
		}
	}

	return &pool, nil
}

func NewConnection(url string) (*Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	connection := &Connection{
		Connection: conn,
	}

	go func() {
		for {
			reason, ok := <-connection.Connection.NotifyClose(make(chan *amqp.Error))
			// exit this goroutine if closed by developer
			if !ok {
				logger.Debug("connection closed")
				break
			}
			logger.Debugf("connection closed, reason: %v", reason)

			// reconnect if not closed by developer
			for {
				// wait 1s for reconnect
				time.Sleep(waitTime)

				conn, err := amqp.Dial(url)
				if err == nil {
					connection.Connection = conn
					logger.Debugf("reconnect success")
					break
				}

				logger.Debugf("reconnect failed, err: %v", err)
			}
		}
	}()

	return connection, nil
}

// Channel wrap amqp.Connection.Channel, get a auto reconnect channel
func (c *Connection) Channel() (*Channel, error) {

	ch, err := c.Connection.Channel()
	if err != nil {
		return nil, err
	}

	channel := &Channel{
		Channel: ch,
	}

	go func() {
		for {
			reason, ok := <-channel.Channel.NotifyClose(make(chan *amqp.Error))
			if !ok {
				logger.Debug("channel closed")
				break
			}
			logger.Debugf("channel closed, reason: %v", reason)

			for {
				time.Sleep(waitTime)
				ch, err := c.Connection.Channel()
				if err == nil {
					logger.Debug("channel recreate success")
					channel.Channel = ch
					break
				}
				logger.Debugf("channel recreate failed, err: %v", err)
			}
		}
	}()
	return channel, nil
}

func (p *Pool) Channel() *Channel {
	return <-p.channels
}

func (p *Pool) ReturnChannel(ch *Channel) {
	p.channels <- ch
}

func (p *Pool) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	ch := p.Channel()
	defer p.ReturnChannel(ch)

	return ch.Publish(exchange, key, mandatory, immediate, msg)
}
