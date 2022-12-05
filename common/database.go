package common

import (
	"context"

	errs "errors"

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

// handles find one errors
// will return the proper dawn error for different errors found based on FindOne singleResult documentation
func SingleResultHandler(res *mongo.SingleResult, v interface{}) error {
	if res.Err() != nil {
		return res.Decode(v)
	} else if errs.Is(res.Err(), mongo.ErrNoDocuments) {
		return errors.NewNotFound(res.Err())
	} else {
		return errors.NewInternal(res.Err())
	}
}

// handles update one errors
func UpdateResultHandler(res *mongo.UpdateResult, err error) error {
	if err != nil {
		return errors.NewInternal(err)
	} else if res.MatchedCount == 0 {
		return errors.NewNotFound(nil)
	} else {
		return nil
	}
}

// handles update one errors
func DeleteResultHandler(res *mongo.DeleteResult, err error) error {
	if err != nil {
		return errors.NewInternal(err)
	} else if res.DeletedCount == 0 {
		return errors.NewNotFound(nil)
	} else {
		return nil
	}
}

// handles update one errors
func InsertOneResultHandler(res *mongo.InsertOneResult, err error) error {
	if err != nil {
		return errors.NewInternal(err)
	} else {
		return nil
	}
}
