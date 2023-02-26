package entities

import "time"

type Address struct {
	AddressLine1 string `bson:"address_line_1" json:"address_line_1"`
	AddressLine2 string `bson:"address_line_2" json:"address_line_2"`
	City         string `bson:"city" json:"city"`
	State        string `bson:"state" json:"state"`
	Zipcode      string `bson:"zipcode" json:"zipcode"`
}
type Location struct {
	ID          string  `bson:"id" json:"id"`
	Latitude    float64 `bson:"latitude" json:"latitude"`
	Longitude   float64 `bson:"longitude" json:"longitude"`
	Address     string  `bson:"address" json:"address"`
	DisplayName string  `bson:"display_name" json:"display_name"`
}

var ROLE_INTS = []int{0, 1, 2, 3}
var ROLES = map[string]int{
	"USER":     0,
	"REPORTER": 1,
	"ADMIN":    2,
	"SUPER":    3,
}
var ROLES_RV = map[int]string{
	0: "USER",
	1: "REPORTER",
	2: "ADMIN",
	3: "SUPER",
}

type User struct {
	ID              string `bson:"_id" json:"_id"`
	Name            string `bson:"name" json:"name"`
	Email           string `bson:"email" json:"emil"`
	Salt            string `bson:"salt" json:"salt"`
	Hash            []byte `bson:"hash" json:"-"`
	HashVersion     string `bson:"hashVersion" json:-"`
	Admin           bool   `bson:"admin" json:"-"`
	Role            int    `bson:"role" json:"role"`
	NewsletterOptIn bool   `bson:"newsletter_opt_in" json:"newsletter_opt_in`

	LastLoggedIn     time.Time `bson:"last_logged_in" json:"last_loged_in"`
	LastTokenRefresh time.Time `bson:"last_token_refresh" json:"-"`
	CreatedAt        time.Time `bson:"created_at" json:"created_at"`

	efaultLocation    Location   `bson:"default_location" json:"default_location"`
	FavoriteLocations []Location `bson:"favorite_locations" json:"favorite_locations"`
}
