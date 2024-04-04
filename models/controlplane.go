package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NodeEntry struct {
	ID            *primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	NodeID        string              `bson:"node_id" json:"nodeID"`
	PublicKeyHash string              `bson:"public_key_hash" json:"publicKeyHash"`
	IP            string              `bson:"ip" json:"ip"`
	Port          int32               `bson:"port" json:"port"`
}
