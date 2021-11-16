package metrics

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
)

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of get requests.",
	},
	[]string{"path"},
)

func New() fiber.Handler {

	return func(c *fiber.Ctx) error {

		totalRequests.WithLabelValues(c.Path()).Inc()

		return c.Next()
	}
}
