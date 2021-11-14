package common

import (
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func FiberLoadBalanceInsert() fiber.Handler {

	return func(c *fiber.Ctx) error {
		if viper.GetString("app.logLevel") == "DEBUG" {
			c.Append("handler", strconv.Itoa(os.Getpid()))
		}
		return c.Next()
	}
}
