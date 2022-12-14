package common

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"github.com/tgs266/dawn-go-common/entities"
	"github.com/tgs266/dawn-go-common/errors"
	"github.com/tgs266/dawn-go-common/jwt"
	"github.com/tgs266/dawn-go-common/optional"
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

func ParseBody[T any](ctx DawnCtx) *optional.Optional[T] {
	var v T
	err := ctx.BodyParser(&v)
	if err != nil {
		err = errors.NewBadRequest(err)
	}
	return optional.New(v, err)
}

var UNAUTHORIZED_TO_USER_ID = &errors.DawnError{
	Name:        "UNAUTHORIZED_TO_USER_ID",
	Description: "user is not authorized to access this endpoint",
	Code:        403,
}

// enfore that a users request is actually for that user.
// will return without checking if the config is local
func (ctx DawnCtx) ValidateToUser(userId string) DawnCtx {
	if ConfigName == "local" {
		return ctx
	}
	if !(ctx.GetUserId() == userId || ctx.GetRole() >= 1) {
		panic(errors.NewForbidden(nil).SetDescription("request is trying to access a resource they don't have access to"))
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
