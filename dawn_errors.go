package common

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type BaseError interface {
	Error() string
}

type DawnError struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	LogDetails  string            `json:"log_details"`
	Code        int               `json:"code"`
	Details     map[string]string `json:"details"`
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
	for k, v := range err.Details {
		details[k] = v
	}
	return StandardError{Source: viper.GetString("app.name"), ErrorCode: err.Name, Description: err.Description, Details: details}
}

func (err *DawnError) AddLogDetails(logDetails string) *DawnError {
	err.LogDetails = logDetails
	return err
}

func (err *DawnError) PutDetail(key string, value string) *DawnError {
	err.Details[key] = value
	return err
}

func Build(err error) *DawnError {
	return &DawnError{
		Name:        "INTERNAL_SERVER_ERROR",
		Description: err.Error(),
		Code:        500,
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
