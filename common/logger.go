package common

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gobwas/glob"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
)

var logLineCount int = 0
var logFileCount int = 1

type RequestLog struct {
	Date       string
	PID        string
	Level      string
	RequestId  string
	Error      *DawnError
	StatusCode string
	Method     string
	Path       string
}

type Request struct {
	Headers fasthttp.RequestHeader
}

type MessageLog struct {
	Date      string
	Level     string
	PID       string
	RequestId string
	Message   string
}

var LEVEL_FORMAT_STRING string = "%-5s"

func buildMessageLog(c *fiber.Ctx, message string) MessageLog {
	const layout = "2006-01-02 03:04:05"
	requestId := c.Locals("requestId")

	messageLog := MessageLog{
		Date:      time.Now().UTC().Format(layout),
		RequestId: fmt.Sprintf("%s", requestId),
		PID:       strconv.Itoa(os.Getpid()),
		Message:   message,
	}
	return messageLog
}

func cleanRequest(c *fiber.Ctx, r *fasthttp.Request) Request {
	headers := fasthttp.AcquireRequest().Header
	r.Header.CopyTo(&headers)
	fmt.Println(headers.String())
	return Request{
		Headers: headers,
	}
}

func BuildMessage(c *fiber.Ctx) RequestLog {
	const layout = "2006-01-02 03:04:05"
	requestId := c.Locals("requestId")

	message := RequestLog{
		Date:       time.Now().UTC().Format(layout),
		RequestId:  fmt.Sprintf("%s", requestId),
		Level:      "INFO",
		StatusCode: strconv.Itoa(c.Response().StatusCode()),
		Method:     c.Method(),
		Path:       c.Path(),
		PID:        strconv.Itoa(os.Getpid()),
	}
	return message
}

func writeToFile(message string, file string) {

	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(message + "\n"); err != nil {
		panic(err)
	}
}

func LogRequest(message RequestLog) {
	txtLogPath := ""
	jsonLogPath := ""

	if viper.GetViper().ConfigFileUsed() == "local" {
		txtLogPath = "text-log-" + strconv.Itoa(logFileCount) + ".log"
		jsonLogPath = "json-log-" + strconv.Itoa(logFileCount) + ".log"
	} else {
		hostname, _ := os.Hostname()
		txtLogPath = "/var/log/" + hostname + "/text-log-" + strconv.Itoa(logFileCount) + ".log"
		jsonLogPath = "/var/log/" + hostname + "/json-log-" + strconv.Itoa(logFileCount) + ".log"
	}

	if message.Error != nil {
		message.Level = "ERROR"
	}

	tempLogString, _ := json.MarshalIndent(message, "", "  ")
	jsonLogString := string(tempLogString)
	txtLogString := fmt.Sprintf("[%s] %s %s %s %s - %s %s", fmt.Sprintf(LEVEL_FORMAT_STRING, message.Level), message.Date, message.PID, message.RequestId, message.StatusCode, message.Method, message.Path)
	if message.Error != nil {
		txtLogString += " - " + message.Error.Error()
	}
	if viper.GetString("app.logType") == "json" {
		fmt.Println(jsonLogString)
	} else {
		fmt.Println(txtLogString)
	}
	writeToFile(jsonLogString, jsonLogPath)
	writeToFile(txtLogString, txtLogPath)

}

func FiberLogger() fiber.Handler {

	return func(c *fiber.Ctx) error {
		errHandler := c.App().Config().ErrorHandler
		chainErr := c.Next()

		message := BuildMessage(c)

		if chainErr != nil {
			dawnError := ErrorHandler(c, chainErr)
			message.Error = dawnError
		}

		if set := viper.IsSet("logging.ignore"); set {
			ignores := viper.GetStringSlice("logging.ignore")
			matched := false
			for _, str := range ignores {
				g := glob.MustCompile(str)
				matched = g.Match(c.Path()) || matched
			}
			if !matched {
				LogRequest(message)
			}
		} else {
			LogRequest(message)
		}

		if chainErr != nil {
			if err := errHandler(c, chainErr); err != nil {
				_ = c.SendStatus(fiber.StatusInternalServerError)
			}
		}

		return nil
	}
}

func ErrorHandler(ctx *fiber.Ctx, err error) *DawnError {
	var returnError *DawnError
	if e, ok := err.(*DawnError); ok {
		returnError = e
	} else {
		returnError = Build(err)
	}

	return returnError
}

/// LOG LEVELS

func stringToLevel(str string) int {
	switch str {
	case "TRACE":
		return 1
	case "DEBUG":
		return 2
	case "INFO":
		return 3
	}
	return 1
}

func TRACE(c *fiber.Ctx, message string) {
	if stringToLevel("TRACE") >= stringToLevel(viper.GetString("app.logLevel")) {
		_log(c, "TRACE", message)
	}
}
func DEBUG(c *fiber.Ctx, message string) {
	if stringToLevel("DEBUG") >= stringToLevel(viper.GetString("app.logLevel")) {
		_log(c, "DEBUG", message)
	}
}
func INFO(c *fiber.Ctx, message string) {
	if stringToLevel("INFO") >= stringToLevel(viper.GetString("app.logLevel")) {
		_log(c, "INFO", message)
	}
}

func _log(c *fiber.Ctx, level, message string) {
	lg := buildMessageLog(c, message)
	lg.Level = level
	logString := ""
	if viper.GetString("app.logType") == "json" {
		tempLogString, _ := json.MarshalIndent(lg, "", "  ")
		logString = string(tempLogString)
	} else {
		logString = fmt.Sprintf("[%s] %s %s %s %s", fmt.Sprintf(LEVEL_FORMAT_STRING, lg.Level), lg.Date, lg.PID, lg.RequestId, lg.Message)
	}
	fmt.Println(logString)
}
