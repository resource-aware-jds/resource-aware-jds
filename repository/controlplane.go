package repository

import (
	"context"
	"errors"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
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
	IsNodeAlreadyRegistered(ctx context.Context, keyHash string) (bool, error)
	RegisterWorkerNodeWithCertificate(ctx context.Context, ip string, port int32, certificate cert.TLSCertificate) error
}

func ProvideControlPlane(database *mongo.Database) IControlPlane {
	return &controlPlane{
		database:               database,
		nodeRegistryCollection: database.Collection(NodeRegistryCollection),
	}
}

func (c *controlPlane) IsNodeAlreadyRegistered(ctx context.Context, publicKeyHash string) (bool, error) {
	result := c.nodeRegistryCollection.FindOne(ctx, models.NodeEntry{
		PublicKeyHash: publicKeyHash,
	})

	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return true, nil
		}
		return false, result.Err()
	}

	return false, nil
}

func (c *controlPlane) RegisterWorkerNodeWithCertificate(ctx context.Context, ip string, port int32, certificate cert.TLSCertificate) error {
	parsePublicKey, err := certificate.GetPublicKey()
	if err != nil {
		return err
	}

	_, err = c.nodeRegistryCollection.InsertOne(ctx, models.NodeEntry{
		NodeID:        certificate.GetCertificate().Subject.SerialNumber,
		PublicKey:     parsePublicKey,
		Certificate:   certificate.GetCertificate().Raw,
		PublicKeyHash: parsePublicKey.GetSHA1Hash(),
		IP:            ip,
		Port:          port,
	})
	if err != nil {
		return err
	}

	return nil
}
