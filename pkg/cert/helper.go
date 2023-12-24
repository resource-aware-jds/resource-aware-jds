package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"net"
	"time"
)

func GeneratePublicAndPrivateKeyPair() (publicKeyData KeyData, privateKeyData KeyData, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)

	if err != nil {
		return nil, nil, err
	}

	return &RSAPublicKeyData{
			data: privateKey.PublicKey,
		}, &RSAPrivateKeyData{
			data: privateKey,
		}, nil
}

type CreateCertificateOptions struct {
	PublicKey            KeyData
	PrivateKey           KeyData
	Duration             time.Duration
	CertificateSubject   pkix.Name
	ParentTLSCertificate TLSCertificate
	IsCA                 bool
}

func CreateCertificate(c CreateCertificateOptions) (TLSCertificate, error) {
	if c.PublicKey == nil || c.PrivateKey == nil {
		return nil, fmt.Errorf("public key or private key is nil")
	}

	// Create the Certificate
	certificate := &x509.Certificate{
		SerialNumber:          big.NewInt(2019),
		Subject:               c.CertificateSubject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(c.Duration),
		IsCA:                  c.IsCA,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IPAddresses: []net.IP{
			net.IPv4(0, 0, 0, 0),
			net.IPv6unspecified,
		},
		DNSNames: []string{"localhost"},
	}

	var parentCertificateValue *x509.Certificate
	keyToSignCertificate := c.PrivateKey.GetRawKeyData()
	if c.ParentTLSCertificate != nil {
		keyToSignCertificate = c.ParentTLSCertificate.GetPrivateKey().GetRawKeyData()
		parentCertificateValue = c.ParentTLSCertificate.GetCertificate()
	} else {
		certificate.IsCA = true
		parentCertificateValue = certificate
	}

	certificateByte, err := x509.CreateCertificate(rand.Reader, certificate, parentCertificateValue, c.PublicKey.GetRawKeyData(), keyToSignCertificate)
	if err != nil {
		return nil, err
	}

	certificateParsed, err := x509.ParseCertificate(certificateByte)
	if err != nil {
		return nil, err
	}

	return &tlsCertificate{
		certificate:       certificateParsed,
		parentCertificate: c.ParentTLSCertificate,
		publicKey:         c.PublicKey,
		privateKey:        c.PrivateKey,
	}, nil
}

func IsSupportedPEMBlock(blockType string) bool {
	switch {
	case blockType == PEMPrivateKeyBlockType:
		return true
	case blockType == PEMCertBlockType:
		return true
	}
	return false
}
