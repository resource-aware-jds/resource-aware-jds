package service

import (
	"context"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/resource-aware-jds/resource-aware-jds/config"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/pool"
	"github.com/resource-aware-jds/resource-aware-jds/repository"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	ErrNodeAlreadyRegistered = errors.New("there is the other node registered with the inputted key")
)

type ControlPlane struct {
	nodeRegistryRepository repository.INodeRegistry
	caCertificate          cert.CACertificate
	config                 config.ControlPlaneConfigModel
	workerNodePool         pool.WorkerNode
}

type IControlPlane interface {
	RegisterWorker(ctx context.Context, ip string, port int32, nodePublicKey cert.KeyData) (certificate cert.TLSCertificate, err error)
	GetAllWorkerNodeFromRegistry(ctx context.Context) ([]models.NodeEntry, error)
	CheckInWorkerNode(ctx context.Context, ip string, port int32, cert []byte) error
}

func ProvideControlPlane(
	nodeRegistryRepository repository.INodeRegistry,
	caCertificate cert.CACertificate,
	config config.ControlPlaneConfigModel,
	workerNodePool pool.WorkerNode,
) IControlPlane {
	return &ControlPlane{
		nodeRegistryRepository: nodeRegistryRepository,
		caCertificate:          caCertificate,
		config:                 config,
		workerNodePool:         workerNodePool,
	}
}

func (s *ControlPlane) RegisterWorker(ctx context.Context, ip string, port int32, nodePublicKey cert.KeyData) (certificate cert.TLSCertificate, err error) {
	hashedPublicKey, err := nodePublicKey.GetSHA1Hash()
	if err != nil {
		return nil, err
	}

	isExists, err := s.nodeRegistryRepository.IsNodeAlreadyRegistered(ctx, hashedPublicKey)
	if err != nil {
		return nil, err
	}
	if isExists {
		return nil, ErrNodeAlreadyRegistered
	}

	// Sign the certificate.
	clientUUID := uuid.New()
	signedCertificate, err := s.caCertificate.CreateCertificateAndSign(
		pkix.Name{
			CommonName:   fmt.Sprintf("RAJDS Worker %s", clientUUID.String()),
			SerialNumber: clientUUID.String(),
		},
		nodePublicKey,
		365*24*time.Hour,
	)
	if err != nil {
		return nil, err
	}

	// Insert the certificate in the database.
	err = s.nodeRegistryRepository.RegisterWorkerNodeWithCertificate(ctx, ip, port, signedCertificate)
	if err != nil {
		return nil, err
	}

	// Response the certificate back.
	return signedCertificate, nil
}

func (s *ControlPlane) GetAllWorkerNodeFromRegistry(ctx context.Context) ([]models.NodeEntry, error) {
	return s.nodeRegistryRepository.GetAllWorkerNode(ctx)
}

func (s *ControlPlane) CheckInWorkerNode(ctx context.Context, ip string, port int32, rawPEMCertificateData []byte) error {
	// Load Certificate
	// Validate the certificate signature
	parsedCertificate, err := cert.LoadCertificate(rawPEMCertificateData)
	if err != nil {
		return err
	}

	if len(parsedCertificate) == 0 {
		return fmt.Errorf("no certificate to verify")
	}

	focusedCertificate := parsedCertificate[0]
	err = s.caCertificate.ValidateSignature(focusedCertificate)
	if err != nil {
		logrus.Error(err)
		return err
	}

	if time.Now().Before(focusedCertificate.NotBefore) || time.Now().After(focusedCertificate.NotAfter) {
		return fmt.Errorf("client certificate expired")
	}

	nodeEntry, err := s.nodeRegistryRepository.GetNode(ctx, cert.GetNodeIDFromCertificate(focusedCertificate))
	if err != nil {
		return err
	}

	nodeEntry.IP = ip
	nodeEntry.Port = port

	// Update Worker Node Stat
	err = s.nodeRegistryRepository.UpdateNodeStatByID(ctx, *nodeEntry)
	if err != nil {
		logrus.WithField("nodeID", nodeEntry.NodeID).Warnf("Failed to update node stat (%s)", err.Error())
	}

	s.workerNodePool.RemoveNodeFromPool(nodeEntry.NodeID)

	// Add Worker Node to the pool
	return s.workerNodePool.AddWorkerNode(ctx, *nodeEntry)
}
