package common

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func FiberLoadBalanceInsert() fiber.Handler {

	return func(c *fiber.Ctx) error {
		if viper.GetString("app.logLevel") == "DEBUG" {
			name, _ := os.Hostname()
			c.Append("handler", name)
		}
		return c.Next()
	}
}
