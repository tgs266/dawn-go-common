package messaging

import (
	"fmt"

	"github.com/streadway/amqp"
)

var (
	Conn *amqp.Connection
)

func Connect() error {
	for Conn == nil {
		conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		if err == nil {
			Conn = conn
		} else {
			fmt.Println(err)
		}
	}
	return nil
}

func OpenChannel() (*amqp.Channel, error) {
	ch, err := Conn.Channel()
	return ch, err
}
