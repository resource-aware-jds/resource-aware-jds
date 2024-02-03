package cert

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

func LoadCertificatesFromFile(certificateFilePath string) ([]*x509.Certificate, error) {
	certificateFile, err := os.ReadFile(certificateFilePath)
	if err != nil {
		return nil, err
	}

	return LoadCertificate(certificateFile)
}

func LoadKeyFromFile(privateKeyFilePath string) (KeyData, error) {
	privateKeyFile, err := os.ReadFile(privateKeyFilePath)
	if err != nil {
		return nil, err
	}

	block, rest := pem.Decode(privateKeyFile)
	if block == nil {
		return nil, nil
	}

	if !IsSupportedPEMBlock(block.Type) {
		return nil, fmt.Errorf("%e: %s", ErrInvalidPEMBlockType, block.Type)
	}

	privateKeyData, err := ParsePrivateKeyToKeyData(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failure reading private key from \"%s\": %s", privateKeyFile, err)
	}

	if len(rest) != 0 {
		logrus.Warn("The Private Key file contain more than one key which won't be loaded.")
	}

	return privateKeyData, nil
}

func LoadCertificate(pemCertificateData []byte) ([]*x509.Certificate, error) {
	certificates := make([]*x509.Certificate, 0)
	for {
		block, rest := pem.Decode(pemCertificateData)
		if block == nil {
			return nil, nil
		}

		if !IsSupportedPEMBlock(block.Type) {
			return nil, fmt.Errorf("%e: %s", ErrInvalidPEMBlockType, block.Type)
		}

		parsedCertificate, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}

		certificates = append(certificates, parsedCertificate)
		pemCertificateData = rest
		if rest == nil || (rest != nil && len(rest) == 0) {
			break
		}
	}

	return certificates, nil
}
