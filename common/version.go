package common

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func RegisterVersionAtPath(app *fiber.App, path string) {
	api := app.Group(path)
	api.Get("version", Version)
}

func RegisterVersion(app *fiber.App) {
	api := app.Group(viper.GetString("server.context-path"))
	api.Get("version", Version)
}

func Version(c *fiber.Ctx) error {

	version := os.Getenv("VERSION")
	if len(version) == 0 {
		version = "unknown"
	}

	return c.Status(fiber.StatusOK).JSON(
		map[string]string{
			"version": version,
		},
	)
}
