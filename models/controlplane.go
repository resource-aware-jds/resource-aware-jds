package models

import (
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NodeEntry struct {
	ID            *primitive.ObjectID `bson:"_id,omitempty"`
	NodeID        string              `bson:"nodeId"`
	PublicKey     cert.RAJDSPublicKey `bson:"public_key"`
	PublicKeyHash string              `bson:"public_key_hash"`
	Certificate   []byte              `bson:"certificate"`
}
