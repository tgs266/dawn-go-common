package common

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"github.com/tgs266/dawn-go-common/errors"
)

func DawnErrorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	stackTrace := ""
	if ctx.Locals("stack_trace") != nil {
		stackTrace = fmt.Sprint(ctx.Locals("stack_trace"))
	}

	message := errors.StandardError{Source: viper.GetString("app.name"), ErrorCode: "INTERNAL_SERVER",
		Description: "Internal Server Error Occurred", Details: map[string]string{"RequestId": ""}}

	if err != nil {
		if e, ok := err.(*errors.DawnError); ok {
			code = e.Code
			message = err.(*errors.DawnError).BuildStandardError(ctx)
		} else {
			err = errors.NewInternal(err)
		}

	} else {
		err = errors.NewUnknown()
	}

	err.(*errors.DawnError).StackTrace = stackTrace

	logMessage := BuildMessage(ctx)
	logMessage.Error = err.(*errors.DawnError)
	logMessage.Level = "ERROR"
	logMessage.StatusCode = strconv.Itoa(code)

	LogRequest(logMessage)

	if code == 500 {
		message = errors.NewInternal(nil).BuildStandardError(ctx)
	}

	errors.ErrorCount.WithLabelValues(logMessage.StatusCode, logMessage.Method, ctx.Route().Path).
		Inc()

	err = ctx.Status(code).JSON(message)

	return nil
}
