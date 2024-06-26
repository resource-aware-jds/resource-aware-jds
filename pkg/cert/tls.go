package cert

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/file"
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

//go:generate mockgen -source=./tls.go -destination=./mock_cert/mock_tls.go -package=mock_cert

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

	ValidateSignature(underValidateCertificate *x509.Certificate) error

	GetNodeID() string
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

	return CreateCertificate(CreateCertificateOptions{
		PublicKey:            subjectPublicKey,
		ValidDuration:        validDuration,
		CertificateSubject:   certificateSubject,
		ParentTLSCertificate: t,
		IsCA:                 false,
		DNSName: []string{
			fmt.Sprintf("%s.%s", certificateSubject.SerialNumber, GetDefaultDomainName()),
		},
	})
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
	err := file.CreateFolderForFile(certificateFilePath)
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

	err = file.CreateFolderForFile(privateKeyFilePath)
	if err != nil {
		return err
	}
	err = os.WriteFile(privateKeyFilePath, privateKeyPEM.Bytes(), 0700)
	if err != nil {
		return err
	}

	return nil
}

func (t *tlsCertificate) ValidateSignature(underValidateCertificate *x509.Certificate) error {
	hash := sha256.New()
	hash.Write(underValidateCertificate.RawTBSCertificate)
	hashData := hash.Sum(nil)
	return rsa.VerifyPKCS1v15(t.publicKey.GetRawKeyData().(*rsa.PublicKey), crypto.SHA256, hashData, underValidateCertificate.Signature)
}

func (t *tlsCertificate) GetNodeID() string {
	return t.GetCertificate().Subject.SerialNumber
}
