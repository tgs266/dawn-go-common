package prometheus

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Client struct {
	constLabels prometheus.Labels

	requestTotal    *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	requestInFlight *prometheus.GaugeVec
	responseSize    *prometheus.HistogramVec
	customCounters  []MiddlewareCounter
}

func New(service string) *Client {

	constLabels := make(prometheus.Labels)
	if service != "" {
		constLabels["service"] = service
	}
	hostname, err := os.Hostname()
	if err == nil {
		constLabels["hostname"] = hostname
	}

	counter := promauto.With(prometheus.DefaultRegisterer).NewCounterVec(
		prometheus.CounterOpts{
			Name:        prometheus.BuildFQName("http", "", "requests_total"),
			Help:        "Count all http requests by status code, method and path.",
			ConstLabels: constLabels,
		},
		[]string{"status_code", "method", "path"},
	)

	histogram := promauto.With(prometheus.DefaultRegisterer).NewHistogramVec(prometheus.HistogramOpts{
		Name:        prometheus.BuildFQName("http", "", "request_duration_seconds"),
		Help:        "Duration of all HTTP requests by status code, method and path.",
		ConstLabels: constLabels,
		Buckets: []float64{
			0.000000001, // 1ns
			0.000000002,
			0.000000005,
			0.00000001, // 10ns
			0.00000002,
			0.00000005,
			0.0000001, // 100ns
			0.0000002,
			0.0000005,
			0.000001, // 1µs
			0.000002,
			0.000005,
			0.00001, // 10µs
			0.00002,
			0.00005,
			0.0001, // 100µs
			0.0002,
			0.0005,
			0.001, // 1ms
			0.002,
			0.005,
			0.01, // 10ms
			0.02,
			0.05,
			0.1, // 100 ms
			0.2,
			0.5,
			1.0, // 1s
			2.0,
			5.0,
			10.0, // 10s
			15.0,
			20.0,
			30.0,
		},
	},
		[]string{"status_code", "method", "path"},
	)

	responseSizeBuckets := promauto.With(prometheus.DefaultRegisterer).NewHistogramVec(prometheus.HistogramOpts{
		Name:        prometheus.BuildFQName("http", "", "response_size_bytes"),
		Help:        "Response size of all HTTP requests by status code, method and path.",
		ConstLabels: constLabels,
		Buckets: []float64{
			5.0,
			10.0,
			25.0,
			50.0,
			100.0,
			250.0,
			500.0,
			1000.0, // 1kb
			2500.0,
			5000.0,
			10000.0,
			25000.0,
			50000.0,
			100000.0,
			150000.0,
			200000.0,
			250000.0,
			500000.0,
			1000000.0, // 1mb
			5000000.0,
			10000000.0,
		},
	},
		[]string{"status_code", "method", "path"},
	)

	gauge := promauto.With(prometheus.DefaultRegisterer).NewGaugeVec(prometheus.GaugeOpts{
		Name:        prometheus.BuildFQName("http", "", "requests_in_progress_total"),
		Help:        "All the requests in progress",
		ConstLabels: constLabels,
	}, []string{"method"})

	return &Client{
		constLabels: constLabels,

		requestTotal:    counter,
		requestDuration: histogram,
		requestInFlight: gauge,
		responseSize:    responseSizeBuckets,

		customCounters: []MiddlewareCounter{},
	}
}

func (c *Client) CreateCustomCounter(name, help string, labelNames []string, function func(ctx *fiber.Ctx, counter *prometheus.CounterVec, args ...string)) CustomCounter {
	counter := promauto.With(prometheus.DefaultRegisterer).NewCounterVec(
		prometheus.CounterOpts{
			Name:        name,
			Help:        help,
			ConstLabels: c.constLabels,
		},
		labelNames,
	)
	cc := CustomCounter{
		counter:  counter,
		function: function,
	}
	return cc
}

func (c *Client) CreateCustomHistogram(name, help string, labelNames []string, buckets []float64, function func(ctx *fiber.Ctx, counter *prometheus.HistogramVec, value float64, args ...string)) CustomHistogram {
	histogram := promauto.With(prometheus.DefaultRegisterer).NewHistogramVec(
		prometheus.HistogramOpts{
			Name:        name,
			Help:        help,
			ConstLabels: c.constLabels,
			Buckets:     buckets,
		},
		labelNames,
	)
	cc := CustomHistogram{
		histogram: histogram,
		function:  function,
	}
	return cc
}

func (c *Client) AddMiddlewareCustomCounter(name, help string, labelNames []string, function func(ctx *fiber.Ctx, counter *prometheus.CounterVec, statusCode string)) MiddlewareCounter {
	counter := promauto.With(prometheus.DefaultRegisterer).NewCounterVec(
		prometheus.CounterOpts{
			Name:        name,
			Help:        help,
			ConstLabels: c.constLabels,
		},
		labelNames,
	)
	cc := MiddlewareCounter{
		counter:  counter,
		function: function,
	}
	c.customCounters = append(c.customCounters, cc)
	return cc
}

func (c *Client) Middleware(ctx *fiber.Ctx) error {

	start := time.Now()
	method := ctx.Route().Method

	if ctx.Route().Path == "/metrics" {
		return ctx.Next()
	}

	c.requestInFlight.WithLabelValues(method).Inc()
	defer func() {
		c.requestInFlight.WithLabelValues(method).Dec()
	}()

	err := ctx.Next()
	status := fiber.StatusInternalServerError
	if err != nil {
		if e, ok := err.(*fiber.Error); ok {
			// Get correct error code from fiber.Error type
			status = e.Code
		}
	} else {
		status = ctx.Response().StatusCode()
	}

	path := ctx.Route().Path
	statusCode := strconv.Itoa(status)
	c.requestTotal.WithLabelValues(statusCode, method, path).Inc()
	elapsed := float64(time.Since(start).Nanoseconds()) / 1e9
	c.requestDuration.WithLabelValues(statusCode, method, path).Observe(elapsed)

	c.responseSize.WithLabelValues(statusCode, method, path).Observe(float64(len(ctx.Response().Body())))

	for _, counter := range c.customCounters {
		counter.function(ctx, counter.counter, statusCode)
	}

	return err
}

// hosting on prom endpoint
func version(w http.ResponseWriter, r *http.Request) {
	version := os.Getenv("VERSION")
	if len(version) == 0 {
		version = "unknown"
	}

	json.NewEncoder(w).Encode(map[string]string{
		"version": version,
	})
}

func (c *Client) Serve() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/version", version)
	http.ListenAndServe(":9216", nil)
}
