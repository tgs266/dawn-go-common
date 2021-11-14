package common

import (
	"context"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Conn       *mongo.Client
	Ctx               = context.Background()
	ConnString string = ""
	DB         *mongo.Database
)

func CreateDBSession(dbName string) error {
	ConnString = viper.GetString("db.uri") + viper.GetString("db.database")
	var err error
	Conn, err = mongo.Connect(Ctx, options.Client().
		ApplyURI(ConnString))
	if err != nil {
		return err
	}
	err = Conn.Ping(Ctx, nil)
	if err != nil {
		return err
	}

	DB = Conn.Database(dbName)

	return nil
}
