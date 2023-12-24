package service

import (
	"context"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/resource-aware-jds/resource-aware-jds/config"
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
	config                 config.ControlPlaneConfigModel
}

type IControlPlane interface {
	RegisterWorker(ctx context.Context, ip string, port int32, nodePublicKey cert.KeyData) (certificate cert.TLSCertificate, err error)
}

func ProvideControlPlane(controlPlaneRepository repository.IControlPlane, caCertificate cert.CACertificate, config config.ControlPlaneConfigModel) IControlPlane {
	return &ControlPlane{
		controlPlaneRepository: controlPlaneRepository,
		caCertificate:          caCertificate,
		config:                 config,
	}
}

func (s *ControlPlane) RegisterWorker(ctx context.Context, ip string, port int32, nodePublicKey cert.KeyData) (certificate cert.TLSCertificate, err error) {
	hashedPublicKey, err := nodePublicKey.GetSHA1Hash()
	if err != nil {
		return nil, err
	}

	isExists, err := s.controlPlaneRepository.IsNodeAlreadyRegistered(ctx, hashedPublicKey)
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
	err = s.controlPlaneRepository.RegisterWorkerNodeWithCertificate(ctx, ip, port, signedCertificate)
	if err != nil {
		return nil, err
	}

	err = signedCertificate.SaveCertificateToFile(fmt.Sprintf("%s/%s.pem", s.config.ClientCertificateStoragePath, signedCertificate.GetCertificate().Subject.SerialNumber), "")
	if err != nil {
		return nil, err
	}

	// Response the certificate back.
	return signedCertificate, nil
}
