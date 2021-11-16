package metrics

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/spf13/viper"
)

var totalRequests = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total_path",
		Help: "Number of get requests.",
	},
	[]string{"path"},
)

func New() fiber.Handler {

	return func(c *fiber.Ctx) error {

		if c.Path() == viper.GetString("server.context-path")+"/metrics" {
			return c.Next()
		}

		totalRequests.WithLabelValues(c.Path()).Inc()

		return c.Next()
	}
}
