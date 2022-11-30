package prometheus

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
)

type CustomCounter struct {
	counter  *prometheus.CounterVec
	function func(ctx *fiber.Ctx, counter *prometheus.CounterVec, args ...string)
}

type CustomHistogram struct {
	histogram *prometheus.HistogramVec
	function  func(ctx *fiber.Ctx, histogram *prometheus.HistogramVec, value float64, args ...string)
}

type MiddlewareCounter struct {
	counter  *prometheus.CounterVec
	function func(ctx *fiber.Ctx, counter *prometheus.CounterVec, statusCode string)
}

// cant get status code when calling trigger
func (c CustomCounter) Trigger(ctx *fiber.Ctx, args ...string) {
	c.function(ctx, c.counter, args...)
}

// cant get status code when calling trigger
func (c CustomHistogram) Trigger(ctx *fiber.Ctx, value float64, args ...string) {
	c.function(ctx, c.histogram, value, args...)
}
