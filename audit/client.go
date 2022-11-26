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

func (a *AuditClient) Audit(collectionName string) *AuditTransaction {
	return &AuditTransaction{
		collectionName: collectionName,
		session:        a.session,
		record: &AuditRecord{
			ID:         uuid.NewString(),
			Collection: collectionName,
			Timestamp:  time.Now(),
		},
	}
}
