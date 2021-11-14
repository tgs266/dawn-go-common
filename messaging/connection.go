package messaging

import (
	"errors"

	"github.com/streadway/amqp"
)

var (
	Conn  *amqp.Connection
	Alive bool
)

func Connect() error {
	if Conn != nil || !Alive {
		conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		Conn = conn
		if err != nil {
			Alive = false
		} else {
			Alive = true
		}
		return err
	}
	return nil
}

func OpenChannel() (*amqp.Channel, error) {
	if !Alive {
		return nil, errors.New("nope")
	}
	ch, err := Conn.Channel()
	return ch, err
}
