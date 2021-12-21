package entities

import "time"

type Address struct {
	AddressLine1 string `bson:"address_line_1"`
	AddressLine2 string `bson:"address_line_2"`
	City         string `bson:"city"`
	State        string `bson:"state"`
	Zipcode      string `bson:"zipcode"`
}
type Location struct {
	ID          string  `bson:"id"`
	Latitude    float64 `bson:"latitude"`
	Longitude   float64 `bson:"longitude"`
	Address     string  `bson:"address"`
	DisplayName string  `bson:"display_name"`
}

var ROLE_INTS = []int{0, 1, 2}
var ROLES = map[string]int{
	"USER":  0,
	"ADMIN": 1,
	"SUPER": 2,
}
var ROLES_REV = map[int]string{
	0: "USER",
	1: "ADMIN",
	2: "SUPER",
}

type User struct {
	ID    string `bson:"_id"`
	Name  string `bson:"name"`
	Email string `bson:"email"`
	Salt  string `bson:"salt"`
	Hash  []byte `bson:"hash"`
	Admin bool   `bson:"admin"`
	Role  int `bson:"role"`

	LastLoggedIn     time.Time `bson:"last_logged_in"`
	LastTokenRefresh time.Time `bson:"last_token_refresh"`

	DefaultLocation   Location   `bson:"default_location"`
	FavoriteLocations []Location `bson:"favorite_locations"`
}
