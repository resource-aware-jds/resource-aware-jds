package repository

import (
	"context"
	"errors"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	NodeRegistryCollection = "node-registry"
)

type nodeRegistry struct {
	database               *mongo.Database
	nodeRegistryCollection *mongo.Collection
}

type INodeRegistry interface {
	IsNodeAlreadyRegistered(ctx context.Context, keyHash string) (bool, error)
	RegisterWorkerNodeWithCertificate(ctx context.Context, ip string, port int32, certificate cert.TLSCertificate) error
	GetAllWorkerNode(ctx context.Context) ([]models.NodeEntry, error)
	GetNode(ctx context.Context, nodeID string) (*models.NodeEntry, error)
	UpdateNodeStatByID(ctx context.Context, nodeEntry models.NodeEntry) error
}

func ProvideControlPlane(database *mongo.Database) INodeRegistry {
	return &nodeRegistry{
		database:               database,
		nodeRegistryCollection: database.Collection(NodeRegistryCollection),
	}
}

func (c *nodeRegistry) IsNodeAlreadyRegistered(ctx context.Context, publicKeyHash string) (bool, error) {
	result := c.nodeRegistryCollection.FindOne(ctx, models.NodeEntry{
		PublicKeyHash: publicKeyHash,
	})

	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, result.Err()
	}

	return true, nil
}

func (c *nodeRegistry) RegisterWorkerNodeWithCertificate(ctx context.Context, ip string, port int32, certificate cert.TLSCertificate) error {
	parsePublicKey := certificate.GetPublicKey()

	publicKeyHash, err := parsePublicKey.GetSHA1Hash()
	if err != nil {
		return err
	}

	_, err = c.nodeRegistryCollection.InsertOne(ctx, models.NodeEntry{
		NodeID:        certificate.GetCertificate().Subject.SerialNumber,
		PublicKeyHash: publicKeyHash,
		IP:            ip,
		Port:          port,
	})
	return err
}

func (c *nodeRegistry) GetAllWorkerNode(ctx context.Context) ([]models.NodeEntry, error) {
	result, err := c.nodeRegistryCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var nodeEntries []models.NodeEntry
	err = result.All(ctx, &nodeEntries)
	return nodeEntries, err
}

func (c *nodeRegistry) GetNode(ctx context.Context, nodeID string) (*models.NodeEntry, error) {
	result := c.nodeRegistryCollection.FindOne(ctx, bson.M{
		"node_id": nodeID,
	})

	if result.Err() != nil {
		return nil, result.Err()
	}

	var nodeEntry models.NodeEntry
	err := result.Decode(&nodeEntry)
	return &nodeEntry, err
}

func (c *nodeRegistry) UpdateNodeStatByID(ctx context.Context, nodeEntry models.NodeEntry) error {
	_, err := c.nodeRegistryCollection.UpdateOne(
		ctx,
		bson.M{
			"node_id": nodeEntry.NodeID,
		},
		bson.M{
			"$set": bson.M{
				"ip":   nodeEntry.IP,
				"port": nodeEntry.Port,
			},
		},
	)
	return err
}
