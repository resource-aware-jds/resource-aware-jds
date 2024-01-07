package cert

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"github.com/resource-aware-jds/resource-aware-jds/config"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
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

func ProvideWorkerNodeTransportCertificate(workerCertificateConfig WorkerNodeTransportCertificateConfig, controlPlaneConfig config.ControlPlaneConfigModel) (TransportCertificate, error) {
	privateKeyData, err := LoadKeyFromFile(workerCertificateConfig.PrivateKeyFileLocation)
	if err != nil {
		response, privateKeyData, err := registerWorker(workerCertificateConfig, controlPlaneConfig)
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

func registerWorker(workerCertificateConfig WorkerNodeTransportCertificateConfig, controlPlaneConfig config.ControlPlaneConfigModel) (*proto.ComputeNodeRegistrationResponse, KeyData, error) {
	caCertificate, err := ProvideWorkerNodeCACertificate(WorkerNodeCACertificateConfig{
		CACertificateFilePath: workerCertificateConfig.CertificateFileLocation,
	})
	if err != nil {
		return nil, nil, err
	}

	grpcConn, err := grpc.ProvideRAJDSGrpcClient(grpc.ClientConfig{
		Target:        controlPlaneConfig.GRPCServerAddress,
		CACertificate: caCertificate,
	})
	if err != nil {
		return nil, nil, err
	}

	controlPlaneClient := proto.NewControlPlaneClient(grpcConn.GetConnection())
	publicKeyData, privateKeyData, err := GeneratePublicAndPrivateKeyPair()
	if err != nil {
		return nil, nil, err
	}

	result, err := controlPlaneClient.WorkerRegistration(context.Background(), &proto.ComputeNodeRegistrationRequest{
		Ip:            "1234",
		Port:          1234,
		NodePublicKey: x509.MarshalPKCS1PublicKey(publicKeyData.GetRawKeyData().(*rsa.PublicKey)),
	})
	if err != nil {
		return nil, nil, err
	}

	return result, privateKeyData, nil
}
