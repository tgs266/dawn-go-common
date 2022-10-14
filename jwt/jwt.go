package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	SharedEntities "github.com/tgs266/dawn-go-common/entities"
	DawnErrors "github.com/tgs266/dawn-go-common/errors"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

type Claims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	ID    string `json:"id"`
	Role  int    `json:"role"`
	jwt.RegisteredClaims
}
type RefreshClaims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

func IssueJWT(user SharedEntities.User) (string, time.Time, error) {
	var err error

	expTime := viper.GetInt("JWT.expiration")
	expiration := time.Now().Add(time.Minute * time.Duration(expTime))
	atClaims := Claims{
		Name:  user.Name,
		Email: user.Email,
		ID:    user.ID,
		Role:  user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			Issuer:    "dawn",
		},
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(viper.GetString("JWT.ACCESS_SECRET")))
	if err != nil {
		return "", expiration, err
	}
	return token, expiration, nil
}

func IssueRefreshToken(user SharedEntities.User) (string, time.Time, error) {
	var err error

	expTime := viper.GetInt("JWT.refresh_expiration")
	expiration := time.Now().Add(time.Minute * time.Duration(expTime))
	atClaims := RefreshClaims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
		},
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(viper.GetString("JWT.ACCESS_SECRET")))
	if err != nil {
		return "", expiration, err
	}
	return token, expiration, nil
}

func ExtractClaims(token string) *Claims {
	out, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(("invalid signing method"))
		}
		return []byte(viper.GetString("JWT.ACCESS_SECRET")), nil
	})
	if err != nil {
		panic(DawnErrors.NewUnauthorized(err).AddLogDetails(err.Error()))
	}

	if claims, ok := out.Claims.(*Claims); ok && out.Valid {
		return claims
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		panic(DawnErrors.NewUnauthorizedInvalid(err))
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		panic(DawnErrors.NewUnauthorized(err))
	} else {
		panic(DawnErrors.NewUnauthorizedInvalid(err))
	}
}

func ExtractClaimsNoError(token string) *Claims {
	out, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, nil
		}
		return []byte(viper.GetString("JWT.ACCESS_SECRET")), nil
	})

	if err != nil {
		return nil
	}

	if claims, ok := out.Claims.(*Claims); ok && out.Valid {
		return claims
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		return nil
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		return nil
	} else {
		return nil
	}
}

func ExtractRefreshClaims(token string) *RefreshClaims {

	out, err := jwt.ParseWithClaims(token, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(("invalid signing method"))
		}
		return []byte(viper.GetString("JWT.ACCESS_SECRET")), nil
	})

	if err != nil {
		panic(DawnErrors.NewUnauthorized(err).AddLogDetails(err.Error()).PutDetail("token", token))
	}

	if claims, ok := out.Claims.(*RefreshClaims); ok && out.Valid {
		return claims
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		panic(DawnErrors.NewUnauthorizedInvalid(err))
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		panic(DawnErrors.NewUnauthorized(err))
	} else {
		panic(DawnErrors.NewUnauthorizedInvalid(err))
	}
}

func ValidateTokenToUser(c *fiber.Ctx, userId string) *Claims {

	token := string(c.Request().Header.Peek("Authorization"))
	claims := ExtractClaims(token)
	if viper.GetBool("app.auth") {
		if userId != claims.ID && (claims.Role < SharedEntities.ROLES["ADMIN"]) {
			panic(DawnErrors.NewUnauthorized(nil).PutDetail("reason", "user is authenticated but not authorized"))
		}
		return claims
	}
	return claims
}

func ValidateToken(c *fiber.Ctx) *Claims {
	token := string(c.Request().Header.Peek("Authorization"))
	return ValidateTokenNoCtx(token)
}

func ValidateTokenNoCtx(token string) *Claims {
	claims := ExtractClaims(token)
	return claims
}
