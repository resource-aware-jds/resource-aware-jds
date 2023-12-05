package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type NodeEntry struct {
	ID        *primitive.ObjectID `bson:"_id,omitempty"`
	PublicKey string              `bson:"public_key"`
}
