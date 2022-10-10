package main

import (
	"fmt"

	"github.com/tgs266/dawn-go-common/errors"
	"github.com/tgs266/dawn-go-common/jwt"
)

func a() *errors.DawnError {
	return errors.NewInternal(nil)
}

func main() {
	// x := entities.User{
	// 	Name:  "t",
	// 	Email: "asdf",
	// 	ID:    "asdf",
	// }

	// t, _, _ := jwt.IssueJWT(x)
	// fmt.Println(t)

	y := jwt.ExtractRefreshClaims("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNiMGJhMmI1LTRiYTQtNGIyNi04YzhiLWU2MDczZGZkZDk1YyIsImV4cCI6MTY2NTk4MTY3Mn0._PpORtwAhFN8e0HPhoKApJf17I8OXNJXFA6byFDfGAA")
	fmt.Println(y)

	// err := a()
	// fmt.Println(err.StackTrace)
}
