package common

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"github.com/tgs266/dawn-go-common/entities"
	"github.com/tgs266/dawn-go-common/errors"
	"github.com/tgs266/dawn-go-common/jwt"
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

func (ctx DawnCtx) Deadline() (deadline time.Time, ok bool) {
	return ctx.FiberCtx.Context().Deadline()
}

func (ctx DawnCtx) Value(key interface{}) interface{} {
	return ctx.FiberCtx.Context().Value(key)
}

func (ctx DawnCtx) Err() error {
	return ctx.FiberCtx.Context().Err()
}

func (ctx DawnCtx) Done() <-chan struct{} {
	return ctx.FiberCtx.Context().Done()
}

func (ctx DawnCtx) BodyParser(out interface{}) error {
	return ctx.FiberCtx.BodyParser(out)
}

var UNAUTHORIZED_TO_USER_ID = &errors.DawnError{
	Name:        "UNAUTHORIZED_TO_USER_ID",
	Description: "user is not authorized to access this endpoint",
	Code:        403,
}

func (ctx DawnCtx) ValidateToUser(userId string) DawnCtx {
	if viper.GetBool("app.auth") {
		admin, _ := strconv.ParseBool(string(ctx.FiberCtx.Request().Header.Peek("admin")))
		if string(ctx.FiberCtx.Request().Header.Peek("user_id")) != userId && !admin {
			panic(UNAUTHORIZED_TO_USER_ID)
		}
	}
	return ctx
}

func (ctx DawnCtx) GetJWT() string {
	return string(ctx.FiberCtx.Request().Header.Peek("Authorization"))
}

func (ctx DawnCtx) GetRole() int {
	token := ctx.GetJWT()
	claims := jwt.ExtractClaims(token)
	return claims.Role
}

func (ctx DawnCtx) GetUserId() string {
	token := ctx.GetJWT()
	claims := jwt.ExtractClaims(token)
	return claims.ID
}

func (ctx DawnCtx) ValidateToAdmin() DawnCtx {
	if viper.GetBool("app.auth") {
		role := ctx.GetRole()
		if role < entities.ROLES["ADMIN"] {
			panic(UNAUTHORIZED_TO_USER_ID)
		}
	}
	return ctx
}

func (ctx DawnCtx) ValidateToSuper() DawnCtx {
	if viper.GetBool("app.auth") {
		role := ctx.GetRole()
		if role != entities.ROLES["SUPER"] {
			panic(UNAUTHORIZED_TO_USER_ID)
		}
	}
	return ctx
}
