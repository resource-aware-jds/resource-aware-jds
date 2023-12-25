package cert

import (
	"bytes"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"net"
	"os"
	"time"
)

const (
	PEMPrivateKeyBlockType = "PRIVATE KEY"
	PEMCertBlockType       = "CERTIFICATE"
)

var (
	ErrInvalidPEMBlockType = errors.New("invalid PEM block type")
	ErrNoPrivateKey        = errors.New("no private key to sign the certificate")
)

type tlsCertificate struct {
	certificate       *x509.Certificate
	parentCertificate TLSCertificate
	privateKey        KeyData
	publicKey         KeyData
}

type TLSCertificate interface {
	IsCA() bool

	GetPublicKey() KeyData
	GetPrivateKey() KeyData
	GetCertificate() *x509.Certificate

	GetCACertificate() (*x509.Certificate, error)
	GetCertificateInPEM() ([]byte, error)
	GetCertificateChains(pemEncoded bool) [][]byte
	GetParentTLSCertificate() TLSCertificate

	CreateCertificateAndSign(certificateSubject pkix.Name, subjectPublicKey KeyData, validDuration time.Duration) (TLSCertificate, error)

	SaveCertificateToFile(certificateFilePath, privateKeyFilePath string) error
	GetCertificateSubjectSerialNumber() string
}

func ProvideTLSCertificate(certificateChain []*x509.Certificate, privateKey KeyData, isCA bool) (TLSCertificate, error) {
	parsedFirstCertificateInChain, err := ParsePublicKeyToKeyData(certificateChain[0].PublicKey)
	if err != nil {
		return nil, err
	}

	certificateChain[0].IsCA = isCA

	firstCertificate := &tlsCertificate{
		certificate: certificateChain[0],
		publicKey:   parsedFirstCertificateInChain,
		privateKey:  privateKey,
	}

	previousTLSCertificate := firstCertificate
	for i := 1; i < len(certificateChain); i++ {
		focusedCertificate := certificateChain[i]
		if focusedCertificate == nil {
			continue
		}

		parsedPublicKeyData, err := ParsePublicKeyToKeyData(focusedCertificate.PublicKey)
		if err != nil {
			return nil, err
		}

		if i == len(certificateChain)-1 {
			(*focusedCertificate).IsCA = true
		}

		latestTlSCertificate := &tlsCertificate{
			certificate: focusedCertificate,
			publicKey:   parsedPublicKeyData,
		}
		previousTLSCertificate.parentCertificate = latestTlSCertificate
		previousTLSCertificate = latestTlSCertificate
	}

	return firstCertificate, nil
}

func (t *tlsCertificate) IsCA() bool {
	return t.certificate.IsCA
}

func (t *tlsCertificate) GetCertificate() *x509.Certificate {
	return t.certificate
}

func (t *tlsCertificate) GetParentTLSCertificate() TLSCertificate {
	return t.parentCertificate
}

func (t *tlsCertificate) GetPublicKey() KeyData {
	return t.publicKey
}

func (t *tlsCertificate) GetPrivateKey() KeyData {
	return t.privateKey
}

func (t *tlsCertificate) GetCertificateInPEM() ([]byte, error) {
	// Encode the current focused TLS Certificate.
	certificatePEM := new(bytes.Buffer)
	err := pem.Encode(certificatePEM, &pem.Block{
		Type:  PEMCertBlockType,
		Bytes: t.GetCertificate().Raw,
	})

	return certificatePEM.Bytes(), err
}

func (t *tlsCertificate) GetCertificateChains(pemEncoded bool) [][]byte {
	// Call pop to reverse the certificate chain.
	certificateStack := make([][]byte, 0)

	var focusedTLSCertificate TLSCertificate
	focusedTLSCertificate = t
	for {

		certByte := focusedTLSCertificate.GetCertificate().Raw
		if pemEncoded {
			var err error
			certByte, err = focusedTLSCertificate.GetCertificateInPEM()
			if err != nil {
				continue
			}
		}

		certificateStack = append(certificateStack, certByte)

		if focusedTLSCertificate.GetParentTLSCertificate() == nil {
			break
		}
		focusedTLSCertificate = focusedTLSCertificate.GetParentTLSCertificate()
	}

	return certificateStack
}

func (t *tlsCertificate) CreateCertificateAndSign(certificateSubject pkix.Name, subjectPublicKey KeyData, validDuration time.Duration) (TLSCertificate, error) {
	if t.privateKey == nil {
		return nil, ErrNoPrivateKey
	}

	// Create the Certificate
	certificate := &x509.Certificate{
		SerialNumber:          big.NewInt(2019),
		Subject:               certificateSubject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(validDuration),
		IsCA:                  false,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IPAddresses: []net.IP{
			net.IPv4(0, 0, 0, 0),
			net.IPv6unspecified,
		},
		DNSNames: []string{"localhost"},
	}

	certificateByte, err := x509.CreateCertificate(rand.Reader, certificate, t.certificate, subjectPublicKey, t.privateKey)
	if err != nil {
		return nil, err
	}

	parsedCertificate, err := x509.ParseCertificate(certificateByte)
	if err != nil {
		return nil, err
	}

	return &tlsCertificate{
		certificate:       parsedCertificate,
		parentCertificate: t,
		publicKey:         subjectPublicKey,
	}, nil
}

func (t *tlsCertificate) GetCertificateSubjectSerialNumber() string {
	return t.certificate.Subject.SerialNumber
}

func (t *tlsCertificate) GetCACertificate() (*x509.Certificate, error) {
	var focusedTLSCertificate TLSCertificate

	focusedTLSCertificate = t
	for {
		if focusedTLSCertificate.IsCA() {
			return focusedTLSCertificate.GetCertificate(), nil
		}

		focusedTLSCertificate = t.GetParentTLSCertificate()
		if focusedTLSCertificate == nil {
			break
		}
	}

	return nil, errors.New("no CA Certificate in this TLS certificate chain")
}

func (t *tlsCertificate) SaveCertificateToFile(certificateFilePath, privateKeyFilePath string) error {
	// Encode the certificate into PEM format
	certificateBytes := t.GetCertificateChains(true)

	// Save the certificates into file
	err := createFolderForFile(certificateFilePath)
	if err != nil {
		return err
	}
	certificateByteJoin := bytes.Join(certificateBytes, []byte(""))
	err = os.WriteFile(certificateFilePath, certificateByteJoin, 0700)
	if err != nil {
		return err
	}

	// If no private key path or empty, not save the private key
	if t.privateKey == nil || privateKeyFilePath == "" {
		return nil
	}

	// Encode Private Key into PEM format
	parsedPrivateKey, err := t.privateKey.GetKeyX509Format()
	if err != nil {
		return err
	}

	privateKeyPEM := new(bytes.Buffer)
	err = pem.Encode(privateKeyPEM, &pem.Block{
		Type:  PEMPrivateKeyBlockType,
		Bytes: parsedPrivateKey,
	})
	if err != nil {
		return err
	}

	err = createFolderForFile(privateKeyFilePath)
	if err != nil {
		return err
	}
	err = os.WriteFile(privateKeyFilePath, privateKeyPEM.Bytes(), 0700)
	if err != nil {
		return err
	}

	return nil
}
