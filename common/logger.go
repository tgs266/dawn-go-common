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
	"github.com/mileusna/useragent"
	"github.com/spf13/viper"
	"github.com/tgs266/dawn-go-common/errors"
	"github.com/tgs266/dawn-go-common/jwt"
)

var logLineCount int = 0
var logFileCount int = 1

type RequestLog struct {
	Proxy       bool              `json:"proxy"`
	UseCache    bool              `json:"useCache"`
	CacheStatus string            `json:"cacheStatus"`
	ServiceName string            `json:"serviceName"`
	Date        string            `json:"date"`
	Level       string            `json:"level"`
	RequestId   string            `json:"requestId"`
	Error       *errors.DawnError `json:"error"`
	StatusCode  string            `json:"statusCode"`
	Method      string            `json:"method"`
	Path        string            `json:"path"`
	UserID      string            `json:"userId"`
	Duration    float64           `json:"duration"`
	Message     string            `json:"message"`
	Request     Request           `json:"request"`
	UserAgent   UserAgent         `json:"userAgent"`
	IPs         []string          `json:"ips"`
	Event       *Event            `json:"event"`
}

type Request struct {
	QueryParams map[string]string `json:"queryParams"`
	Headers     map[string]string `json:"headers"`
	Cookies     map[string]string `json:"cookies"`
}

type Event struct {
	ID         string                 `json:"id" bson:"id"`
	UserId     string                 `json:"userId" bson:"userId"`
	PagePath   string                 `json:"pagePath" bson:"pagePath"`
	Category   string                 `json:"category" bson:"category"`
	Action     string                 `json:"action" bson:"action"`
	Parameters map[string]interface{} `json:"parameters" bson:"parameters"`
	CreatedAt  time.Time              `json:"createdAt" bson:"createdAt"`
}

type UserAgent struct {
	Bot       bool   `json:"bot"`
	Tablet    bool   `json:"tablet"`
	Mobile    bool   `json:"mobile"`
	Desktop   bool   `json:"desktop"`
	Device    string `json:"device"`
	OS        string `json:"os"`
	OSVersion string `json:"osVersion"`
	Name      string `json:"name"`
	Version   string `json:"version"`
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

func ParseUserAgent(c *fiber.Ctx) UserAgent {
	userAgent := string(c.Request().Header.UserAgent())
	ua := useragent.Parse(userAgent)

	return UserAgent{
		Name:      ua.Name,
		Version:   ua.Version,
		OS:        ua.OS,
		OSVersion: ua.OSVersion,
		Mobile:    ua.Mobile,
		Tablet:    ua.Tablet,
		Desktop:   ua.Desktop,
		Bot:       ua.Bot,
	}
}

func BuildMessage(c *fiber.Ctx) *RequestLog {
	requestId := c.Locals("requestId")
	start := c.Locals("start")
	event := c.Locals("event")
	proxyInterface := c.Locals("proxy")
	ucInterface := c.Locals("useCache")
	cStatus := c.Locals("cacheStatus")
	durationFloat := -1.0

	proxy := false
	useCache := false
	cacheStatus := ""
	if proxyInterface != nil {
		proxy = proxyInterface.(bool)
	}
	if ucInterface != nil {
		useCache = ucInterface.(bool)
	}
	if cStatus != nil {
		cacheStatus = cStatus.(string)
	}

	cookies := map[string]string{}
	reqHeaders := map[string]string{}
	queryParams := map[string]string{}
	c.Request().Header.VisitAll(func(k, v []byte) {
		reqHeaders[string(k)] = string(v)
	})
	c.Request().Header.VisitAllCookie(func(k, v []byte) {
		cookies[string(k)] = string(v)
	})

	c.Request().URI().QueryArgs().VisitAll(func(k, v []byte) {
		queryParams[string(k)] = string(v)
	})

	resHeaders := map[string]string{}
	c.Response().Header.VisitAll(func(k, v []byte) {
		resHeaders[string(k)] = string(v)
	})

	if start != nil {
		durationFloat = float64(time.Since(start.(time.Time)).Nanoseconds()) / 1000000
	}

	var actualEvent *Event
	if event != nil {
		actualEvent = event.(*Event)
	}

	var userId string
	if (DawnCtx{FiberCtx: c}.GetJWT() != "") {
		claims := jwt.ExtractClaimsNoError(DawnCtx{FiberCtx: c}.GetJWT())
		if claims != nil {
			userId = claims.ID
		}
	} else {
		userId = ""
	}

	message := &RequestLog{
		Proxy:       proxy,
		UseCache:    useCache,
		CacheStatus: cacheStatus,
		ServiceName: viper.GetString("app.name"),
		Date:        time.Now().Format(time.RFC3339),
		RequestId:   fmt.Sprintf("%s", requestId),
		Level:       "INFO",
		StatusCode:  strconv.Itoa(c.Response().StatusCode()),
		Method:      c.Method(),
		Path:        c.Path(),
		UserID:      userId,
		Duration:    durationFloat,
		IPs:         c.IPs(),
		Request: Request{
			Headers:     reqHeaders,
			QueryParams: queryParams,
			Cookies:     cookies,
		},
		UserAgent: ParseUserAgent(c),
		Event:     actualEvent,
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
		nextErr := c.Next()

		message := BuildMessage(c)

		if nextErr != nil {
			dawnError := ErrorConverter(c, nextErr)
			message.Error = dawnError
			if err := errHandler(c, nextErr); err != nil {
				_ = c.SendStatus(fiber.StatusInternalServerError)
			}
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

		return nil
	}
}

// converts any errors to dawn errors
func ErrorConverter(ctx *fiber.Ctx, err error) *errors.DawnError {
	var returnError *errors.DawnError
	if e, ok := err.(*errors.DawnError); ok {
		returnError = e
	} else {
		returnError = errors.NewInternal(err)
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
