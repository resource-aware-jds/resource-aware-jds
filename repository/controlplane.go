package repository

import "go.mongodb.org/mongo-driver/mongo"

type controlPlane struct {
	database *mongo.Database
}

type IControlPlane interface {
}

func ProvideControlPlane(database *mongo.Database) IControlPlane {
	return &controlPlane{
		database: database,
	}
}
