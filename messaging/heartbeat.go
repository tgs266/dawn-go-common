package messaging

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

type Heartbeat struct {
	Status   string `json:"status"`
	DBStatus string `json:"dbstatus"`
	HostName string `json:"hostname"`
}

func PublishHeartbeat(heartbeat Heartbeat) {
	b, _ := json.Marshal(heartbeat)
	Publish("heartbeat", b)
}

func GetHeartbeatMessageStream() <-chan amqp.Delivery {
	Connect()
	DeclareConsumerQueue("heartbeat")
	return CreateMessageConsumer("heartbeat")
}
