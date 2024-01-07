package cert

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
)

type WorkerNodeCACertificate TLSCertificate

type WorkerNodeCACertificateConfig struct {
	CACertificateFilePath string
}

func ProvideWorkerNodeCACertificate(config WorkerNodeCACertificateConfig) (WorkerNodeCACertificate, error) {
	certificateChain, err := LoadCertificatesFromFile(config.CACertificateFilePath)
	if err != nil {
		return nil, err
	}

	return ProvideTLSCertificate(certificateChain, nil, true)
}

type WorkerNodeTransportCertificateConfig struct {
	CertificateFileLocation string
	PrivateKeyFileLocation  string
}

func ProvideWorkerNodeTransportCertificate(workerCertificateConfig WorkerNodeTransportCertificateConfig, controlPLaneClient proto.ControlPlaneClient) (TransportCertificate, error) {
	privateKeyData, err := LoadKeyFromFile(workerCertificateConfig.PrivateKeyFileLocation)
	if err != nil {
		response, privateKeyData, err := registerWorker(controlPLaneClient)
		if err != nil {
			return nil, err
		}
		certificate, err := LoadCertificate(response.Certificate)
		if err != nil {
			return nil, err
		}
		provideTLSCertificate, err := ProvideTLSCertificate(certificate, privateKeyData, false)
		if err != nil {
			return nil, err
		}
		err = provideTLSCertificate.SaveCertificateToFile(workerCertificateConfig.CertificateFileLocation, workerCertificateConfig.PrivateKeyFileLocation)
		if err != nil {
			return nil, err
		}
		return provideTLSCertificate, nil
	}

	certificateChain, err := LoadCertificatesFromFile(workerCertificateConfig.CertificateFileLocation)
	if err != nil {
		return nil, err
	}

	return ProvideTLSCertificate(certificateChain, privateKeyData, false)
}

func registerWorker(controlPlaneClient proto.ControlPlaneClient) (*proto.ComputeNodeRegistrationResponse, KeyData, error) {
	publicKeyData, privateKeyData, err := GeneratePublicAndPrivateKeyPair()
	result, err := controlPlaneClient.WorkerRegistration(context.Background(), &proto.ComputeNodeRegistrationRequest{
		Port:          1234,
		NodePublicKey: x509.MarshalPKCS1PublicKey(publicKeyData.GetRawKeyData().(*rsa.PublicKey)),
	})
	if err != nil {
		return nil, nil, err
	}

	return result, privateKeyData, nil
}
