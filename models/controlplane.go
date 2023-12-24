package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NodeEntry struct {
	ID            *primitive.ObjectID `bson:"_id,omitempty"`
	NodeID        string              `bson:"nodeId"`
	PublicKeyHash string              `bson:"public_key_hash"`
	IP            string              `bson:"ip"`
	Port          int32               `bson:"port"`
}
