package common

import (
	"context"

	"github.com/spf13/viper"
	"github.com/tgs266/dawn-go-common/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DBSession struct {
	Client     *mongo.Client
	Ctx        context.Context
	ConnString string
	DB         *mongo.Database
	DBName     string
}

var (
	Conn       *mongo.Client
	Ctx               = context.Background()
	ConnString string = ""
	DB         *mongo.Database
	DBName     = ""
)

func CreateHealthDBSession() (*DBSession, error) {
	connString := viper.GetString("db.uri") + viper.GetString("db.database")

	session := &DBSession{
		Ctx:        context.Background(),
		ConnString: connString,
	}
	conn, err := mongo.Connect(Ctx, options.Client().
		ApplyURI(session.ConnString))
	if err != nil {
		return session, err
	}
	session.Client = conn
	return session, err
}

func CreateDBSession(dbName string) (*DBSession, error) {
	connString := viper.GetString("db.uri") + dbName

	session := &DBSession{
		Ctx:        context.Background(),
		DBName:     dbName,
		ConnString: connString,
	}
	err := session.Connect()

	return session, err
}

func (session *DBSession) Connect() error {
	conn, err := mongo.Connect(Ctx, options.Client().
		ApplyURI(session.ConnString))
	if err != nil {
		return err
	}
	err = conn.Ping(session.Ctx, nil)
	if err != nil {
		return err
	}

	db := conn.Database(session.DBName)

	session.Client = conn
	session.DB = db
	return nil
}

func (session *DBSession) Collection(colName string) *mongo.Collection {
	if session.DB != nil {
		return session.DB.Collection(colName)
	} else {
		if err := session.Connect(); err != nil {
			panic(errors.INTERNAL_SERVER_STANDARD_ERROR.PutDetail("reason", err.Error()))
		} else {
			return session.DB.Collection(colName)
		}
	}
}

func (session *DBSession) Ping() error {
	if session.Client != nil {
		return session.Client.Ping(session.Ctx, readpref.Primary())
	}
	return errors.INTERNAL_SERVER_STANDARD_ERROR
}

func (session *DBSession) Close() {
	if session.DB != nil {
		session.Client.Disconnect(session.Ctx)
	}
}
