package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	conn, err := NewPool("amqp://guest:guest@172.16.0.30:5672/")
	if err != nil {
		log.Fatalf("%s: %s", "Failed to connect to RabbitMQ", err)
	}

	//ch, err := conn.Channel()
	//if err != nil {
	//	log.Fatalf("%s: %s", "Failed to open a Channel", err)
	//}
	//defer ch.Close()

	body := "Hello World!"
	err = conn.Publish(
		"test",                 // exchange
		"tr.4335.ytyt343.greg", // routing key
		false,                  // mandatory
		false,                  // immediate
		amqp.Publishing{
			Headers: map[string]interface{}{
				"topic":     "tr/4335/ytyt343/greg",
				"productId": "123456",
				"deviceId":  "76342239423",
			},
			Body: []byte(body),
		})
	if err != nil {
		log.Fatalf("%s: %s", "Failed to open a Channel", err)
	}

	time.Sleep(time.Hour)
}
