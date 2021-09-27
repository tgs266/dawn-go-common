package common

import (
	"github.com/gofiber/fiber/v2"
)

type DawnCtx struct {
	FiberCtx *fiber.Ctx
}

func (ctx DawnCtx) INFO(message string) {
	INFO(ctx.FiberCtx, message)
}

func (ctx DawnCtx) DEBUG(message string) {
	DEBUG(ctx.FiberCtx, message)
}

func (ctx DawnCtx) TRACE(message string) {
	TRACE(ctx.FiberCtx, message)
}

func BuildCtx(c *fiber.Ctx) DawnCtx {
	return DawnCtx{
		FiberCtx: c,
	}
}

func (ctx DawnCtx) BodyParser(out interface{}) error {
	return ctx.FiberCtx.BodyParser(out)
}
