package models

type Location struct {
	Longitude float64 `bson:"longitude"`
	Latitude  float64 `bson:"latitude"`
}

type Address struct {
	AddressLine1 string `bson:"address_line_1"`
	AddressLine2 string `bson:"address_line_2"`
	Zipcode      string `bson:"zipcode"`
	City         string `bson:"city"`
	State        string `bson:"state"`
}

type LocationAddress struct {
	Longitude float64 `bson:"longitude"`
	Latitude  float64 `bson:"latitude"`
	Address   Address `bson:"address"`
}
