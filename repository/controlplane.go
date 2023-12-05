package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	NodeRegistryCollection = "node-registry"
)

type controlPlane struct {
	database               *mongo.Database
	nodeRegistryCollection *mongo.Collection
}

type IControlPlane interface {
	IsNodeAlreadyRegistered(ctx context.Context, nodePublicKeyBase64 string) (bool, error)
}

func ProvideControlPlane(database *mongo.Database) IControlPlane {
	return &controlPlane{
		database:               database,
		nodeRegistryCollection: database.Collection(NodeRegistryCollection),
	}
}

func (c *controlPlane) IsNodeAlreadyRegistered(ctx context.Context, nodePublicKeyBase64 string) (bool, error) {
	// TODO: Implement better encryption.
	result := c.nodeRegistryCollection.FindOne(ctx, bson.M{
		"key": nodePublicKeyBase64,
	})

	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return true, nil
		}
		return false, result.Err()
	}

	return false, nil
}
