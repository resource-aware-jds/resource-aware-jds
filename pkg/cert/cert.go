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
	return ProvideTLSCertificate(Config{
		CertificateFileLocation: caCertificateConfig.CertificateFileLocation,
		PrivateKeyFileLocation:  caCertificateConfig.PrivateKeyFileLocation,
		CertificateSubject: pkix.Name{
			CommonName: "Resource Aware Job Distribution CA",
		},
		ValidDuration:     24 * 365 * 10 * time.Hour,
		ParentCertificate: nil,
	})
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
	return ProvideTLSCertificate(Config{
		CertificateFileLocation: transportCertificateConfig.CertificateFileLocation,
		PrivateKeyFileLocation:  transportCertificateConfig.PrivateKeyFileLocation,
		CertificateSubject: pkix.Name{
			CommonName:   transportCertificateConfig.CommonName,
			SerialNumber: transportCertificateConfig.NodeID,
		},
		ValidDuration:     transportCertificateConfig.ValidDuration,
		ParentCertificate: caCertificate,
	})
}
