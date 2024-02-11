package mongo

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ProvideMongoConnection(config Config) (*mongo.Database, func(), error) {
	opts := options.Client().ApplyURI(config.ConnectionString)

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, func() {}, err
	}

	cleanup := func() {
		err := client.Disconnect(context.Background())
		if err != nil {
			logrus.Error("Mongo Disconnect Error: ", err)
		}
	}

	return client.Database(config.DatabaseName), cleanup, nil
}
