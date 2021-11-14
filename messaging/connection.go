package messaging

import (
	"github.com/streadway/amqp"
)

var (
	Conn *amqp.Connection
)

func Connect() error {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	Conn = conn
	return err
}

func OpenChannel() (*amqp.Channel, error) {
	ch, err := Conn.Channel()
	return ch, err
}
