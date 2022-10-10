package main

import (
	"fmt"

	"github.com/tgs266/dawn-go-common/errors"
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

	// y := jwt.ExtractClaims(t)
	// fmt.Println(y)

	err := a()
	fmt.Println(err.StackTrace)
}
