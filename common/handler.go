package common

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tgs266/dawn-go-common/errors"
)

// handles converting errors to be output
func DawnErrorHandler(ctx *fiber.Ctx, err error) error {
	var outputErr *errors.DawnError
	if v, ok := err.(*errors.DawnError); ok {
		outputErr = v
	} else {
		outputErr = errors.NewInternal(err)
	}

	message := outputErr.BuildStandardError(ctx)

	err = ctx.Status(outputErr.Code).JSON(message)

	return nil
}
