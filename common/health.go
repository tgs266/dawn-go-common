package common

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"github.com/tgs266/dawn-go-common/messaging"
)

var HealthSession *DBSession

type HealthStruct struct {
	Status   string `json:"status"`
	DBStatus string `json:"dbstatus"`
	// Details  map[string]string `json:"details"`
}

func Health(c *fiber.Ctx) error {
	ctx := BuildCtx(c)
	return c.Status(fiber.StatusOK).JSON(
		HealthService(ctx),
	)
}

func GetHealthStruct() HealthStruct {
	status := "available"
	dbstatus := ""
	if viper.IsSet("db.uri") {
		dbstatus = "up"
		if HealthSession == nil {
			HealthSession, _ = CreateHealthDBSession()
		}
		if err := HealthSession.Ping(); err != nil {
			status = "unavailable"
			dbstatus = "down"
		} else {
			dbstatus = "up"
		}
	}
	return HealthStruct{
		Status:   status,
		DBStatus: dbstatus,
	}
}

func HealthService(c DawnCtx) *HealthStruct {
	val := GetHealthStruct()
	return &val
}

func RegisterHealth(app *fiber.App) {
	api := app.Group(viper.GetString("server.context-path"))
	api.Get("health", Health)
}

var LastHeartbeat = messaging.Heartbeat{}

func RegisterHeartbeatPublisher() {
	messaging.Connect(viper.GetString("app.messaging-uri"))
	messaging.DeclarePublisherQueue("heartbeat")

	PublishHeartbeat()
	StartHeartbeatMessenger()
}

func PublishHeartbeat() {
	hostname, _ := os.Hostname()

	healthStruct := GetHealthStruct()
	if healthStruct.Status != LastHeartbeat.Status {
		heartBeat := messaging.Heartbeat{
			Status:      healthStruct.Status,
			DBStatus:    healthStruct.DBStatus,
			HostName:    hostname,
			RequireAuth: viper.GetBool("security.auth"),
			ContextPath: viper.GetString("server.context-path"),
		}
		messaging.PublishHeartbeat(heartBeat)
		LastHeartbeat = heartBeat
	}
}

func ForcePublishHeartbeat() {
	hostname, _ := os.Hostname()

	healthStruct := GetHealthStruct()
	heartBeat := messaging.Heartbeat{
		Status:      healthStruct.Status,
		DBStatus:    healthStruct.DBStatus,
		HostName:    hostname,
		RequireAuth: viper.GetBool("security.auth"),
		ContextPath: viper.GetString("server.context-path"),
	}
	messaging.PublishHeartbeat(heartBeat)
	LastHeartbeat = heartBeat
}

func StartTellAllConsumer() {
	hostname, _ := os.Hostname()

	messaging.Connect(viper.GetString("app.messaging-uri"))
	messaging.DeclareConsumerQueue("send_heartbeat-" + hostname)
	q, _ := messaging.GetQueue("send_heartbeat-" + hostname)
	q.Bind("send_heartbeat_exchange")
	msgs := messaging.CreateMessageConsumer("send_heartbeat-" + hostname)
	go func() {
		for range msgs {
			ForcePublishHeartbeat()
		}
	}()
}

func StartHeartbeatMessenger() {
	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})

	messaging.Connect(viper.GetString("app.messaging-uri"))
	messaging.DeclarePublisherQueue("heartbeat")

	StartTellAllConsumer()

	go func() {
		for {
			select {
			case <-ticker.C:
				ForcePublishHeartbeat()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func SendDeadHeartbeat() {
	if viper.GetBool("app.messaging") {
		hostname, _ := os.Hostname()

		healthStruct := GetHealthStruct()
		heartBeat := messaging.Heartbeat{
			Status:      "stopped",
			DBStatus:    healthStruct.DBStatus,
			HostName:    hostname,
			RequireAuth: viper.GetBool("security.auth"),
			ContextPath: viper.GetString("server.context-path"),
		}
		messaging.PublishHeartbeat(heartBeat)
		LastHeartbeat = heartBeat
		messaging.Close()
	}
}
