package messaging

import (
	"fmt"

	"github.com/streadway/amqp"
)

var (
	Conn *amqp.Connection
)

func Connect() error {
	fmt.Println("trying to connect")
	if Conn != nil {
		conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		Conn = conn
		if err != nil {
			fmt.Println("connection failed")
			fmt.Println(err)
		} else {
			fmt.Println("success")
		}
		return err
	}
	return nil
}

func OpenChannel() (*amqp.Channel, error) {
	ch, err := Conn.Channel()
	return ch, err
}
