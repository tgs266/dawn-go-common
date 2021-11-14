package messaging

import (
	"errors"
	"fmt"

	"github.com/streadway/amqp"
)

type Queue struct {
	Name         string
	Channel      *amqp.Channel
	Publisher    amqp.Queue
	Consumer     amqp.Queue
	HasPublisher bool
	HasConsumer  bool
}

var Queues = make(map[string]Queue)

func GetQueue(name string) (Queue, error) {
	if queue, ok := Queues[name]; ok {
		return queue, nil
	}
	return Queue{}, errors.New("Nope")
}

func MakeQueue(name string) (Queue, error) {
	ch, err := OpenChannel()
	if err != nil {
		return Queue{}, err
	}
	return Queue{
		Name:    name,
		Channel: ch,
	}, nil

}

func DeclarePublisherQueue(name string) {
	queue, err := GetQueue(name)
	if err != nil {
		queue, err = MakeQueue(name)
		if err != nil {
			fmt.Println(err)
		}
	}

	if queue.HasConsumer {
		return
	}

	q, err2 := queue.Channel.QueueDeclare(
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

	queue.Publisher = q
	queue.HasPublisher = true
	Queues[name] = queue
}

func DeclareConsumerQueue(name string) {

	queue, err := GetQueue(name)
	if err != nil {
		queue, err = MakeQueue(name)
		if err != nil {
			fmt.Println(err)
		}
	}

	if queue.HasConsumer {
		return
	}

	q, err2 := queue.Channel.QueueDeclare(
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

	queue.Consumer = q
	queue.HasConsumer = true
	Queues[name] = queue
}

func Publish(name string, data []byte) {
	queue, _ := GetQueue(name)

	err := queue.Channel.Publish(
		"",                   // exchange
		queue.Publisher.Name, // routing key
		false,                // mandatory
		false,                // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		},
	)
	if err != nil {
		fmt.Println("cant publish")
	}
}

func TestPublish(name string, test string) {
	body := test

	queue, _ := GetQueue(name)

	err := queue.Channel.Publish(
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
	queue, _ := GetQueue(name)

	msgs, err := queue.Channel.Consume(
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
