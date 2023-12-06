package cert

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"math/big"
	"net"
	"os"
	"strings"
	"time"
)

const (
	PEMPrivateKeyBlockType = "PRIVATE KEY"
	PEMCertBlockType       = "CERTIFICATE"
)

var (
	ErrInvalidPEMBlockType = errors.New("invalid PEM block type")
)

type tlsCertificate struct {
	isCA              bool
	certificate       *x509.Certificate
	parentCertificate TLSCertificate
	privateKey        crypto.PrivateKey
	publicKey         crypto.PublicKey
}

type TLSCertificate interface {
	IsCA() bool
	GetCertificate() *x509.Certificate
	GetParentCertificate() TLSCertificate
	GetPublicKey() crypto.PublicKey
	GetPrivateKey() crypto.PrivateKey
	GetCertificateChains(pemEncoded bool) [][]byte
}

type Config struct {
	CertificateFileLocation string
	PrivateKeyFileLocation  string
	CertificateSubject      pkix.Name
	ParentCertificate       TLSCertificate
	ValidDuration           time.Duration
}

func ProvideTLSCertificate(config Config) (TLSCertificate, error) {
	certificate := tlsCertificate{}

	// Try to load the Certificate from the file
	err := certificate.loadCertificateFromFile(config.CertificateFileLocation, config.PrivateKeyFileLocation)
	if err == nil {
		logrus.Info("Loaded Certificate from file: ", config.CertificateFileLocation, ":", config.PrivateKeyFileLocation)
		return &certificate, nil
	}

	logrus.Warn("Failed to load certificate from file with this error: ", err)

	// Create the new certificate instead.
	err = certificate.createCertificate(config.ValidDuration, config.CertificateSubject, config.ParentCertificate)
	if err != nil {
		logrus.Error("Failed to create new certificate with this error: ", err)
		return nil, err
	}

	// Save the created certificate to file
	err = certificate.saveCertificateToFile(config.CertificateFileLocation, config.PrivateKeyFileLocation)
	if err != nil {
		logrus.Error("Failed to save the created certificate with this error", err)
		return nil, err
	}
	return &certificate, nil
}

func (t *tlsCertificate) IsCA() bool {
	return t.isCA
}

func (t *tlsCertificate) GetCertificate() *x509.Certificate {
	return t.certificate
}

func (t *tlsCertificate) GetParentCertificate() TLSCertificate {
	return t.parentCertificate
}

func (t *tlsCertificate) GetPublicKey() crypto.PublicKey {
	return t.publicKey
}

func (t *tlsCertificate) GetPrivateKey() crypto.PrivateKey {
	return t.privateKey
}

func (t *tlsCertificate) GetCertificateChains(pemEncoded bool) [][]byte {
	// Call pop to reverse the certificate chain.
	certificateStack := make([][]byte, 0)

	var focusedTLSCertificate TLSCertificate
	focusedTLSCertificate = t
	for {

		certByte := focusedTLSCertificate.GetCertificate().Raw
		if pemEncoded {
			// Encode the current focused TLS Certificate.
			certificatePEM := new(bytes.Buffer)
			err := pem.Encode(certificatePEM, &pem.Block{
				Type:  PEMCertBlockType,
				Bytes: focusedTLSCertificate.GetCertificate().Raw,
			})

			if err != nil {
				continue
			}

			certByte = certificatePEM.Bytes()
		}

		certificateStack = append(certificateStack, certByte)

		if focusedTLSCertificate.GetParentCertificate() == nil {
			break
		}
		focusedTLSCertificate = focusedTLSCertificate.GetParentCertificate()
	}

	return certificateStack
}

// CreateCertificate is used to generate the Public and Private Key pair
// and generate the x509 certificate using the generated Ker pair.
func (t *tlsCertificate) createCertificate(duration time.Duration, certificateSubject pkix.Name, parentTLSCertificate TLSCertificate) error {
	// Create the Key pair.
	publicKey, privateKey, err := t.createPublicAndPrivateKeyPair()
	if err != nil {
		return err
	}

	t.privateKey = privateKey

	// Create the Certificate
	certificate := &x509.Certificate{
		SerialNumber:          big.NewInt(2019),
		Subject:               certificateSubject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(duration),
		IsCA:                  t.isCA,
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
	keyToSignCertificate := privateKey
	if parentTLSCertificate != nil {
		keyToSignCertificate = parentTLSCertificate.GetPrivateKey()
		parentCertificateValue = parentTLSCertificate.GetCertificate()
	} else {
		certificate.IsCA = true
		t.isCA = true
		parentCertificateValue = certificate
	}

	certificateByte, err := x509.CreateCertificate(rand.Reader, certificate, parentCertificateValue, publicKey, keyToSignCertificate)
	if err != nil {
		return err
	}

	certificateParsed, err := x509.ParseCertificate(certificateByte)
	if err != nil {
		return err
	}

	t.certificate = certificateParsed
	t.parentCertificate = parentTLSCertificate
	logrus.Info("")
	return nil
}

// createPublicAndPrivateKeyPair is used to create key used in the certificate procedure.
func (t *tlsCertificate) createPublicAndPrivateKeyPair() (crypto.PublicKey, crypto.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}
	return &privateKey.PublicKey, privateKey, nil
}

func (t *tlsCertificate) saveCertificateToFile(certificateFilePath, privateKeyFilePath string) error {
	// Encode the certificate into PEM format
	certificateBytes := t.GetCertificateChains(true)

	// Save the certificates into file
	certificateFilePathSplit := strings.Split(certificateFilePath, "/")
	certificateFilePathSplit = certificateFilePathSplit[0 : len(certificateFilePathSplit)-1]
	certificateFileLocation := strings.Join(certificateFilePathSplit, "/")

	err := os.MkdirAll(certificateFileLocation, os.ModePerm)
	if err != nil {
		return err
	}
	certificateByteJoin := bytes.Join(certificateBytes, []byte(""))
	err = os.WriteFile(certificateFilePath, certificateByteJoin, 0700)
	if err != nil {
		return err
	}

	// Encode Private Key into PEM format
	ecPrivateKey, err := x509.MarshalPKCS8PrivateKey(t.privateKey)
	if err != nil {
		return err
	}

	privateKeyPEM := new(bytes.Buffer)
	err = pem.Encode(privateKeyPEM, &pem.Block{
		Type:  PEMPrivateKeyBlockType,
		Bytes: ecPrivateKey,
	})
	if err != nil {
		return err
	}

	privateKeyFilePathSplit := strings.Split(privateKeyFilePath, "/")
	privateKeyFilePathSplit = privateKeyFilePathSplit[0 : len(privateKeyFilePathSplit)-1]
	privateKeyFileLocation := strings.Join(privateKeyFilePathSplit, "/")
	err = os.MkdirAll(privateKeyFileLocation, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.WriteFile(privateKeyFilePath, privateKeyPEM.Bytes(), 0700)
	if err != nil {
		return err
	}

	return nil
}

func (t *tlsCertificate) loadCertificateFromFile(certificateFilePath, privateKeyFilePath string) error {
	// Load the Private Key and store in the struct.
	privateKey, err := t.loadPrivateKeyFromFile(privateKeyFilePath)
	if err != nil {
		return err
	}

	parsedTLSCertificate, err := t.loadCertificateFromFilePath(certificateFilePath)
	if err != nil {
		return err
	}
	*t = *parsedTLSCertificate
	t.privateKey = privateKey

	// Validate the public and private key.
	privateKeyParsed, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return fmt.Errorf("private key is not rsa")
	}

	publicKeyParsed, ok := privateKeyParsed.Public().(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("PublicKey parsed failed")
	}

	if !publicKeyParsed.Equal(t.certificate.PublicKey) {
		return fmt.Errorf("no Private Key matched with the certificate")
	}

	return nil
}

func (t *tlsCertificate) loadPrivateKeyFromFile(privateKeyFilePath string) (crypto.PrivateKey, error) {
	privateKeyFile, err := os.ReadFile(privateKeyFilePath)
	if err != nil {
		return nil, err
	}

	block, rest := pem.Decode(privateKeyFile)
	if block == nil {
		return nil, nil
	}

	if !t.isSupportedPEM(block.Type) {
		return nil, fmt.Errorf("%e: %s", ErrInvalidPEMBlockType, block.Type)
	}

	privateKey, err := t.parsePrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failure reading private key from \"%s\": %s", privateKeyFile, err)
	}

	if rest != nil && len(rest) != 0 {
		logrus.Warn("The Private Key file contain more than one key which won't be loaded.")
	}

	return privateKey, nil
}

func (t *tlsCertificate) isSupportedPEM(blockType string) bool {
	switch {
	case blockType == PEMPrivateKeyBlockType:
		return true
	case blockType == PEMCertBlockType:
		return true
	}
	return false
}

func (t *tlsCertificate) loadCertificateFromFilePath(certificateFilePath string) (*tlsCertificate, error) {
	certificateFile, err := os.ReadFile(certificateFilePath)
	if err != nil {
		return nil, err
	}

	certificates := make([]*x509.Certificate, 0)
	for {
		block, rest := pem.Decode(certificateFile)
		if block == nil {
			return nil, nil
		}

		if !t.isSupportedPEM(block.Type) {
			return nil, fmt.Errorf("%e: %s", ErrInvalidPEMBlockType, block.Type)
		}

		parsedCertificate, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}

		certificates = append(certificates, parsedCertificate)
		certificateFile = rest
		if rest == nil || (rest != nil && len(rest) == 0) {
			break
		}
	}

	firstCertificate := &tlsCertificate{
		isCA:        false,
		certificate: certificates[0],
		publicKey:   certificates[0].PublicKey,
	}
	previousTLSCertificate := firstCertificate
	for i := 1; i < len(certificates); i++ {
		focusedCertificate := certificates[i]
		if focusedCertificate == nil {
			continue
		}

		latestTlSCertificate := &tlsCertificate{
			certificate:       focusedCertificate,
			publicKey:         focusedCertificate.PublicKey,
			parentCertificate: previousTLSCertificate,
			isCA:              i == len(certificates)-1,
		}

		previousTLSCertificate = latestTlSCertificate
	}
	return firstCertificate, nil
}

func (t *tlsCertificate) parsePrivateKey(der []byte) (crypto.PrivateKey, error) {
	// Parse Private Key by using PKCS8 standard.
	key, err := x509.ParsePKCS8PrivateKey(der)
	if err != nil {
		return nil, err
	}

	return key, nil
}
