package common

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/spf13/viper"
)

var ErrorCount *prometheus.CounterVec

type BaseError interface {
	Error() string
}

type DawnError struct {
	StackTrace  string            `json:"stack_trace"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	LogDetails  string            `json:"log_details"`
	Code        int               `json:"code"`
	Details     map[string]string `json:"details"`
	ServiceName string            `json:"service_name"`
}

type StandardError struct {
	Source      string            `json:"source"`
	ErrorCode   string            `json:"errorCode"`
	Description string            `json:"description"`
	Details     map[string]string `json:"details"`
}

func (err *DawnError) Error() string {
	str := err.Name + ": " + err.Description
	if err.LogDetails != "" {
		str += " - " + err.LogDetails
	}
	return str
}

func (err *DawnError) BuildStandardError(ctx *fiber.Ctx) StandardError {
	requestId := ctx.Locals("requestId")
	details := map[string]string{"RequestId": fmt.Sprintf("%s", requestId)}

	serviceName := err.ServiceName
	if serviceName == "" {
		serviceName = viper.GetString("app.name")
	}

	for k, v := range err.Details {
		details[k] = v
	}
	return StandardError{Source: serviceName, ErrorCode: err.Name, Description: err.Description, Details: details}
}

func (err *DawnError) AddLogDetails(logDetails string) *DawnError {
	err.LogDetails = logDetails
	return err
}

func (err *DawnError) PutDetail(key string, value string) *DawnError {
	if err.Details == nil || len(err.Details) == 0 {
		err.Details = map[string]string{key: value}
	} else {
		err.Details[key] = value
	}
	return err
}

func (err *DawnError) ChangeServiceName(name string) *DawnError {
	err.ServiceName = name
	return err
}

func Build(err error) *DawnError {
	return &DawnError{
		Name:        "INTERNAL_SERVER_ERROR",
		Description: err.Error(),
		Code:        500,
		StackTrace:  "",
	}
}

func (err *DawnError) LogJson(c *fiber.Ctx) {
	jsonErrBytes, _ := json.Marshal(err)
	fmt.Println(string(jsonErrBytes))
}

func (err *DawnError) LogString(c *fiber.Ctx) {
	requestId := c.Locals("requestId")
	output := strconv.Itoa(os.Getpid()) + " " + fmt.Sprintf("%s", requestId) + " " + strconv.Itoa(err.Code) + " - " + c.Method() + " " + c.Route().Path + " - " + err.Error()
	if err.LogDetails != "" {
		output += " - " + err.LogDetails
	}
	fmt.Println(output)
}

func RegisterDawnPrometheus() {
	constLabels := make(prometheus.Labels)
	constLabels["service"] = viper.GetString("app.name")

	ErrorCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:        prometheus.BuildFQName("http", "", "requests_total_error"),
			Help:        "Count all http requests by status code, method and path.",
			ConstLabels: constLabels,
		},
		[]string{"status_code", "method", "path"},
	)
}

var INTERNAL_SERVER_STANDARD_ERROR = &DawnError{
	Name:        "INTERNAL_SERVER_ERROR",
	Description: "Unkown internal server error occurred",
	Code:        500,
}

func DawnErrorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	stackTrace := ""
	if ctx.Locals("stack_trace") != nil {
		stackTrace = fmt.Sprint(ctx.Locals("stack_trace"))
	}

	message := StandardError{Source: viper.GetString("app.name"), ErrorCode: "INTERNAL_SERVER",
		Description: "Internal Server Error Occurred", Details: map[string]string{"RequestId": ""}}

	if e, ok := err.(*DawnError); ok {
		code = e.Code
		message = err.(*DawnError).BuildStandardError(ctx)
	} else {
		err = Build(err)
	}

	err.(*DawnError).StackTrace = stackTrace

	logMessage := BuildMessage(ctx)
	logMessage.Error = err.(*DawnError)
	logMessage.Level = "ERROR"
	logMessage.StatusCode = strconv.Itoa(code)

	LogRequest(logMessage)

	if code == 500 {
		message = INTERNAL_SERVER_STANDARD_ERROR.BuildStandardError(ctx)
	}

	ErrorCount.WithLabelValues(logMessage.StatusCode, logMessage.Method, ctx.Route().Path).
		Inc()

	err = ctx.Status(code).JSON(message)

	return nil
}
