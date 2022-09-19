package messaging

import (
	"github.com/streadway/amqp"
)

var (
	Conn *amqp.Connection
)

func Connect(url string) error {
	if Conn == nil {
		conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		if err != nil {
			return nil
		} else {
			Conn = conn
		}
	}
	return nil
}

func Close() error {
	if Conn == nil {
		return nil
	}
	return Conn.Close()
}

func OpenChannel() (*amqp.Channel, error) {
	ch, err := Conn.Channel()
	return ch, err
}
