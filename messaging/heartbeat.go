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

func GetHeartbeatMessageStream() <-chan amqp.Delivery {
	Connect()
	DeclareConsumerQueue("heartbeat")
	return CreateMessageConsumer("heartbeat")
}

func TellAllToSendHeartbeats() {
	Connect()
	DeclarePublisherQueue("send_heartbeat")
	Publish("send_heartbeat", []byte("send"))
}
