package main

import (
	"fmt"

	"github.com/tgs266/dawn-go-common/entities"
	"github.com/tgs266/dawn-go-common/jwt"
)

func main() {
	x := entities.User{
		Name:  "t",
		Email: "asdf",
		ID:    "asdf",
	}

	t, _, _ := jwt.IssueJWT(x)
	fmt.Println(t)

	y := jwt.ExtractClaims(t)
	fmt.Println(y)
}
