package messaging

import (
	"fmt"

	"github.com/streadway/amqp"
)

type Queue struct {
	Name      string
	Publisher amqp.Queue
	Consumer  amqp.Queue
}

var Queues = make(map[string]Queue)

func DeclarePublisherQueue(name string) {
	ch, err := OpenChannel()
	if err != nil {
		fmt.Println("cant open channel")
	}
	q, err2 := ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err2 != nil {
		fmt.Println("cant open queue")
	}

	if queue, ok := Queues["foo"]; ok {
		queue.Publisher = q
		Queues[name] = queue
	} else {
		queue := Queue{
			Name:      name,
			Publisher: q,
		}
		Queues[name] = queue
	}
}

func DeclareConsumerQueue(name string) {
	ch, err := OpenChannel()
	if err != nil {
		fmt.Println("cant open channel")
	}
	q, err2 := ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err2 != nil {
		fmt.Println("cant open queue")
	}
	if queue, ok := Queues["foo"]; ok {
		queue.Consumer = q
		Queues[name] = queue
	} else {
		queue := Queue{
			Name:     name,
			Consumer: q,
		}
		Queues[name] = queue
	}
}

func TestPublish(name string, test string) {
	body := test
	ch, err := OpenChannel()
	if err != nil {
		fmt.Println("cant open channel")
	}

	queue := Queues[name]

	err = ch.Publish(
		"",                   // exchange
		queue.Publisher.Name, // routing key
		false,                // mandatory
		false,                // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		fmt.Println("cant publish")
	}
}

func CreateMessageConsumer(name string) <-chan amqp.Delivery {
	ch, err := OpenChannel()
	if err != nil {
		fmt.Println("cant open channel")
	}

	queue := Queues[name]

	msgs, err := ch.Consume(
		queue.Consumer.Name, // queue
		"",                  // consumer
		true,                // auto-ack
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,                 // args
	)
	if err != nil {
		fmt.Println("cant publish")
	}
	return msgs
}
