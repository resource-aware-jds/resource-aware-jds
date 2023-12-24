package cert

import (
	"crypto/x509/pkix"
	"time"
)

type CACertificate TLSCertificate

type CACertificateConfig struct {
	CertificateFileLocation string
	PrivateKeyFileLocation  string
}

func ProvideCACertificate(caCertificateConfig CACertificateConfig) (CACertificate, error) {
	privateKeyData, err := LoadKeyFromFile(caCertificateConfig.PrivateKeyFileLocation)
	if err != nil {
		// If Load Private Key failed, Create the new certificate instead. (For control plane)
		publicKey, privateKey, err := GeneratePublicAndPrivateKeyPair()
		if err != nil {
			return nil, err
		}

		certificate, err := CreateCertificate(CreateCertificateOptions{
			PublicKey:  publicKey,
			PrivateKey: privateKey,
			CertificateSubject: pkix.Name{
				CommonName: "Resource Aware Job Distribution CA",
			},
			ValidDuration:        24 * 365 * 10 * time.Hour,
			ParentTLSCertificate: nil,
			IsCA:                 true,
		})
		if err != nil {
			return nil, err
		}

		err = certificate.SaveCertificateToFile(caCertificateConfig.CertificateFileLocation, caCertificateConfig.PrivateKeyFileLocation)
		if err != nil {
			return nil, err
		}

		return certificate, nil
	}

	certificateChain, err := LoadCertificatesFromFile(caCertificateConfig.CertificateFileLocation)
	if err != nil {
		return nil, err
	}

	return ProvideTLSCertificate(certificateChain, privateKeyData)
}

type TransportCertificate TLSCertificate

type TransportCertificateConfig struct {
	CertificateFileLocation string
	PrivateKeyFileLocation  string
	ValidDuration           time.Duration
	CommonName              string
	NodeID                  string
}

func ProvideTransportCertificate(transportCertificateConfig TransportCertificateConfig, caCertificate CACertificate) (TransportCertificate, error) {
	privateKeyData, err := LoadKeyFromFile(transportCertificateConfig.PrivateKeyFileLocation)
	if err != nil {
		// If Load Private Key failed, Create the new certificate instead. (For control plane)
		publicKey, privateKey, err := GeneratePublicAndPrivateKeyPair()
		if err != nil {
			return nil, err
		}

		certificate, err := CreateCertificate(CreateCertificateOptions{
			PublicKey:  publicKey,
			PrivateKey: privateKey,
			CertificateSubject: pkix.Name{
				CommonName:   transportCertificateConfig.CommonName,
				SerialNumber: transportCertificateConfig.NodeID,
			},
			ValidDuration:        transportCertificateConfig.ValidDuration,
			ParentTLSCertificate: caCertificate,
		})
		if err != nil {
			return nil, err
		}

		err = certificate.SaveCertificateToFile(transportCertificateConfig.CertificateFileLocation, transportCertificateConfig.PrivateKeyFileLocation)
		if err != nil {
			return nil, err
		}

		return certificate, nil
	}

	certificateChain, err := LoadCertificatesFromFile(transportCertificateConfig.CertificateFileLocation)
	if err != nil {
		return nil, err
	}

	return ProvideTLSCertificate(certificateChain, privateKeyData)
}
