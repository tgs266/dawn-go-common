package common

import (
	"context"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"gitlab.cs.umd.edu/dawn/dawn-go-common/messaging"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type HealthStruct struct {
	Status   string `json:"status"`
	DBStatus string `json:"dbstatus"`
	// Details  map[string]string `json:"details"`
}

// func (err *HealthStruct) AddLogDetails(key string, value string) *HealthStruct {
// 	err.Details[key] = value
// 	return err
// }

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
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		if err := Conn.Ping(ctx, readpref.Primary()); err != nil {
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

func StartHeartbeatMessenger() {
	ticker := time.NewTicker(30 * time.Second)
	quit := make(chan struct{})

	hostname, _ := os.Hostname()

	messaging.Connect()
	messaging.DeclarePublisherQueue("heartbeat")

	go func() {
		for {
			select {
			case <-ticker.C:
				healthStruct := GetHealthStruct()
				heartBeat := messaging.Heartbeat{
					Status:   healthStruct.Status,
					DBStatus: healthStruct.DBStatus,
					HostName: hostname,
				}
				messaging.PublishHeartbeat(heartBeat)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
