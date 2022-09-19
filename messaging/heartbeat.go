package messaging

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

type Heartbeat struct {
	Status      string `json:"status"`
	DBStatus    string `json:"dbstatus"`
	HostName    string `json:"hostname"`
	RequireAuth bool   `json:"require_auth"`
	ContextPath string `json:"context_path"`
}

func DecodeHeartbeat(data []byte) *Heartbeat {
	heartbeat := &Heartbeat{}
	json.Unmarshal(data, heartbeat)
	return heartbeat
}

func EncodeHeartbeat(heartbeat Heartbeat) []byte {
	b, _ := json.Marshal(heartbeat)
	return b
}

func PublishHeartbeat(heartbeat Heartbeat) {
	Publish("heartbeat", EncodeHeartbeat(heartbeat))
}

func GetHeartbeatMessageStream(url string) <-chan amqp.Delivery {
	Connect(url)
	DeclareConsumerQueue("heartbeat")
	return CreateMessageConsumer("heartbeat")
}

func TellAllToSendHeartbeats(url string) {
	Connect(url)
	DeclarePublisherQueue("send_heartbeat")
	q, _ := GetQueue("send_heartbeat")
	q.Channel.ExchangeDeclare(
		"send_heartbeat_exchange", // name
		"fanout",                  // type
		true,                      // durable
		false,                     // auto-deleted
		false,                     // internal
		false,                     // no-wait
		nil,                       // arguments
	)
	PublishOverExchange("send_heartbeat", "send_heartbeat_exchange", []byte("send"))
}
