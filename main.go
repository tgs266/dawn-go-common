package main

import (
	"fmt"

	"github.com/tgs266/dawn-go-common/entities"
	"github.com/tgs266/dawn-go-common/errors"
)

func a() *errors.DawnError {
	return errors.NewInternal(nil)
}

func main() {
	fmt.Println(entities.ROLES)
	// x := entities.User{
	// 	Name:  "t",
	// 	Email: "asdf",
	// 	ID:    "asdf",
	// }

	// t, _, _ := jwt.IssueJWT(x)
	// fmt.Println(t)

	// err := a()
	// fmt.Println(err.StackTrace)
}
