package common

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gobwas/glob"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
)

var logLineCount int = 0
var logFileCount int = 1

type RequestLog struct {
	ServiceName string     `json:"serviceName"`
	Date        string     `json:"date"`
	Level       string     `json:"level"`
	RequestId   string     `json:"requestId"`
	Error       *DawnError `json:"error"`
	StatusCode  string     `json:"statusCode"`
	Method      string     `json:"method"`
	Path        string     `json:"path"`
	UserID      string     `json:"userId"`
	Duration    float64    `json:"duration"`
	Message     string     `json:"message"`
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
	requestId := c.Locals("requestId")

	messageLog := MessageLog{
		Date:      time.Now().UTC().Format(time.RFC3339),
		RequestId: fmt.Sprintf("%s", requestId),
		PID:       strconv.Itoa(os.Getpid()),
		Message:   message,
	}
	return messageLog
}

func BuildMessage(c *fiber.Ctx) *RequestLog {
	requestId := c.Locals("requestId")
	start := c.Locals("start")
	durationFloat := -1.0

	reqHeaders := map[string]string{}
	c.Request().Header.VisitAll(func(k, v []byte) {
		reqHeaders[string(k)] = string(v)
	})

	resHeaders := map[string]string{}
	c.Response().Header.VisitAll(func(k, v []byte) {
		resHeaders[string(k)] = string(v)
	})

	if start != nil {
		durationFloat = float64(time.Since(start.(time.Time)).Nanoseconds()) / 1000000
	}

	message := &RequestLog{
		ServiceName: viper.GetString("app.name"),
		Date:        time.Now().Format(time.RFC3339),
		RequestId:   fmt.Sprintf("%s", requestId),
		Level:       "INFO",
		StatusCode:  strconv.Itoa(c.Response().StatusCode()),
		Method:      c.Method(),
		Path:        c.OriginalURL(),
		UserID:      string(c.Request().Header.Peek("user_id")),
		Duration:    durationFloat,
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

func ClearLogFolder() {
	logFolder := ""

	if viper.GetViper().ConfigFileUsed() == "local" {
		logFolder = ""
	} else {
		hostname, _ := os.Hostname()
		logFolder = "/var/log/" + hostname + "/"
	}

	files, err := filepath.Glob(logFolder + "*-log-*.log")
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}
}

func LogRequest(message *RequestLog) {
	txtLogPath := ""
	jsonLogPath := ""
	logFolder := ""

	if strings.Contains(viper.GetViper().ConfigFileUsed(), "local") {
		logFolder = ""
		txtLogPath = "text-log-" + strconv.Itoa(logFileCount) + ".log"
		jsonLogPath = "json-log-" + strconv.Itoa(logFileCount) + ".log"
	} else {
		hostname, _ := os.Hostname()
		logFolder = "/var/log/" + hostname + "/"
		txtLogPath = logFolder + "text-log-" + strconv.Itoa(logFileCount) + ".log"
		jsonLogPath = logFolder + "json-log-" + strconv.Itoa(logFileCount) + ".log"
	}

	if message.Error != nil {
		message.Level = "ERROR"
	}

	if _, err := os.Stat(txtLogPath); os.IsNotExist(err) {
		os.MkdirAll(logFolder, 0700)
	}

	if _, err := os.Stat(jsonLogPath); os.IsNotExist(err) {
		os.MkdirAll(logFolder, 0700)
	}

	message.Message = fmt.Sprintf("[%s] %s %s", message.Method, message.StatusCode, message.Path)

	if message.Error != nil {
		message.Message += " - " + message.Error.Error()
	} else {
	}

	tempLogString, _ := json.Marshal(message)
	jsonLogString := string(tempLogString)
	txtLogString := fmt.Sprintf("[%s] %s %s %s - %s %s", fmt.Sprintf(LEVEL_FORMAT_STRING, message.Level), message.Date, message.RequestId, message.StatusCode, message.Method, message.Path)

	if message.Error != nil {
		txtLogString += " - " + message.Error.Error()
	}

	if viper.GetString("app.logType") == "json" {
		fmt.Println(jsonLogString)
	} else {
		fmt.Println(txtLogString)
	}

}

func FiberLogger() fiber.Handler {

	return func(c *fiber.Ctx) error {
		errHandler := c.App().Config().ErrorHandler
		now := time.Now()
		c.Locals("start", now)
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
