package audit

import (
	"time"

	"github.com/google/uuid"
	"github.com/tgs266/dawn-go-common/common"
)

type AuditClient struct {
	databaseName string
	session      *common.DBSession
}

var Client *AuditClient

func Init(name string) {
	Client = New(name)
}

func New(name string) *AuditClient {
	session, err := common.CreateDBSession(name)
	if err != nil {
		panic(err)
	}
	return &AuditClient{
		databaseName: name,
		session:      session,
	}
}

func (a *AuditClient) Audit(collectionName string, f func(tx *AuditTransaction)) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	tx := &AuditTransaction{
		collectionName: collectionName,
		session:        a.session,
		record: &AuditRecord{
			ID:         uuid.NewString(),
			Collection: collectionName,
			Timestamp:  time.Now(),
		},
	}
	f(tx)
}
