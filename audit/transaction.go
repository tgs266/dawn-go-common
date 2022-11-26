package audit

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"github.com/tgs266/dawn-go-common/common"
)

type AuditTransaction struct {
	collectionName string
	record         *AuditRecord
	session        *common.DBSession
}

type AuditRecord struct {
	// autogenerated
	ID string `json:"id" bson:"_id"`
	// autogenerated
	Timestamp  time.Time `json:"timestamp" bson:"timestamp"`
	Collection string    `json:"collection" bson:"collection"`

	EntityID       string      `json:"entityId" bson:"entityId"`
	Action         AuditAction `json:"type" bson:"type"`
	Actor          string      `json:"actor" bson:"actor"`
	ModifiedFields []string    `json:"modifiedFields,omitempty" bson:"modifiedFields,omitempty"`
}

func contains(arr []string, v string) bool {
	for _, a := range arr {
		if a == v {
			return true
		}
	}
	return false
}

func (tx *AuditTransaction) Action(v AuditAction) *AuditTransaction {
	tx.record.Action = v
	return tx
}

func (tx *AuditTransaction) Actor(v string) *AuditTransaction {
	tx.record.Actor = v
	return tx
}

func (tx *AuditTransaction) ModifiedField(v string) *AuditTransaction {
	if !contains(tx.record.ModifiedFields, v) {
		tx.record.ModifiedFields = append(tx.record.ModifiedFields, v)
	}
	return tx
}

func (tx *AuditTransaction) EntityID(v string) *AuditTransaction {
	tx.record.EntityID = v
	return tx
}

func (tx *AuditTransaction) validate() {
	if tx.record.Actor == "" {
		panic("audit record actor field must be filled")
	}
	if tx.record.Action == "" {
		panic("audit record action field must be filled")
	}
}

func (tx *AuditTransaction) Store() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				if viper.GetBool("logging.audit.warn") {
					fmt.Printf("WARNING: audit logging failed. Error: %v\n", r)
				}
			}
		}()
		tx.validate()
		_, err := tx.session.DB.Collection("audit").InsertOne(context.TODO(), tx.record)
		if err != nil {
			panic(err)
		}
	}()
}
