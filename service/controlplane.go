package service

import (
	"context"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
	"github.com/resource-aware-jds/resource-aware-jds/repository"
	"time"
)

var (
	ErrNodeAlreadyRegistered = errors.New("there is the other node registered with the inputted key")
)

type ControlPlane struct {
	controlPlaneRepository repository.IControlPlane
	caCertificate          cert.CACertificate
}

type IControlPlane interface {
	RegisterWorker(ctx context.Context, ip string, port int32, nodePublicKey cert.RAJDSPublicKey) (certificate cert.TLSCertificate, err error)
}

func ProvideControlPlane(controlPlaneRepository repository.IControlPlane, caCertificate cert.CACertificate) IControlPlane {
	return &ControlPlane{
		controlPlaneRepository: controlPlaneRepository,
		caCertificate:          caCertificate,
	}
}

func (s *ControlPlane) RegisterWorker(ctx context.Context, ip string, port int32, nodePublicKey cert.RAJDSPublicKey) (certificate cert.TLSCertificate, err error) {
	isExists, err := s.controlPlaneRepository.IsNodeAlreadyRegistered(ctx, nodePublicKey.GetSHA1Hash())
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
	err = s.controlPlaneRepository.RegisterWorkerNodeWithCertificate(ctx, signedCertificate)
	if err != nil {
		return nil, err
	}

	// Response the certificate back.
	return signedCertificate, nil
}
