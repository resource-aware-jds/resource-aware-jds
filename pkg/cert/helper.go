package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"net"
	"os"
	"strings"
	"time"
)

func GeneratePublicAndPrivateKeyPair() (publicKeyData KeyData, privateKeyData KeyData, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)

	if err != nil {
		return nil, nil, err
	}

	return &RSAPublicKeyData{
			data: &privateKey.PublicKey,
		}, &RSAPrivateKeyData{
			data: privateKey,
		}, nil
}

type CreateCertificateOptions struct {
	PublicKey            KeyData
	PrivateKey           KeyData
	ValidDuration        time.Duration
	CertificateSubject   pkix.Name
	ParentTLSCertificate TLSCertificate
	IsCA                 bool
	DNSName              []string
}

func CreateCertificate(c CreateCertificateOptions) (TLSCertificate, error) {
	if c.PublicKey == nil {
		return nil, fmt.Errorf("public key is nil")
	}

	// Create the Certificate
	certificate := &x509.Certificate{
		SerialNumber:          big.NewInt(2019),
		Subject:               c.CertificateSubject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(c.ValidDuration),
		IsCA:                  c.IsCA,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IPAddresses: []net.IP{
			net.IPv4(127, 0, 0, 1),
			net.IPv6loopback,
			net.IPv4zero,
			net.IPv6unspecified,
		},
		DNSNames: c.DNSName,
	}

	var parentCertificateValue *x509.Certificate
	var keyToSignCertificate any
	if c.ParentTLSCertificate != nil {
		keyToSignCertificate = c.ParentTLSCertificate.GetPrivateKey().GetRawKeyData()
		parentCertificateValue = c.ParentTLSCertificate.GetCertificate()
	} else {
		keyToSignCertificate = c.PrivateKey.GetRawKeyData()
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

func createFolderForFile(filePath string) error {
	filePathSplit := strings.Split(filePath, "/")
	filePathSplit = filePathSplit[0 : len(filePathSplit)-1]
	pathWithoutFileJoined := strings.Join(filePathSplit, "/")
	return os.MkdirAll(pathWithoutFileJoined, os.ModePerm)
}

func GetNodeIDFromCertificate(cert *x509.Certificate) string {
	return cert.Subject.SerialNumber
}
