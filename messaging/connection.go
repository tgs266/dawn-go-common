package messaging

import (
	"fmt"

	"github.com/streadway/amqp"
)

func Connect() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	fmt.Println(conn)
	fmt.Println(err)
}
