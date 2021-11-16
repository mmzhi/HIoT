package rabbitmq

import (
	"github.com/ruixiaoedu/hiot/adapter"
	"github.com/streadway/amqp"
	"strings"
)

type Client struct {
	pool     *Pool
	engine   adapter.Engine
	exchange string
}

func NewClient(engine adapter.Engine) (adapter.Bridge, error) {
	if engine.Config().Rabbitmq.Url == "" {
		return nil, nil
	}

	pool, err := NewPool(engine.Config().Rabbitmq.Url, Options{
		ConnectionNum: engine.Config().Rabbitmq.ConnectionNum,
		ChannelNum:    engine.Config().Rabbitmq.ChannelNum,
	})

	if err != nil {
		return nil, err
	}

	var exchange string
	if engine.Config().Rabbitmq.Exchange == "" {
		exchange = "hiot"
	} else {
		exchange = engine.Config().Rabbitmq.Exchange
	}

	return &Client{
		pool:     pool,
		engine:   engine,
		exchange: exchange,
	}, nil
}

func (c *Client) Push(topic string, data []byte) error {
	return c.pool.Publish(
		c.exchange,                           // exchange
		strings.Replace(topic, "/", ".", -1), // routing key
		false,                                // mandatory
		false,                                // immediate
		amqp.Publishing{
			Headers: map[string]interface{}{
				"topic": topic,
			},
			Body: data,
		})
}
