package cert

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
)

type KeyType string

const (
	PKIPublicKeyType  KeyType = "pki-public-key"
	PKIPrivateKeyType KeyType = "pki-private-key"
)

type KeyData interface {
	GetKeyX509Format() ([]byte, error)
	GetSHA1Hash() (string, error)
	GetKeyType() KeyType
	GetRawKeyData() any
}

func ParsePrivateKeyToKeyData(unparsedPrivateKey []byte) (KeyData, error) {
	key, err := x509.ParsePKCS1PrivateKey(unparsedPrivateKey)
	if err != nil {
		return nil, err
	}

	return &RSAPrivateKeyData{
		data: key,
	}, nil
}

func ParsePublicKeyToKeyData(unparsedPublicKey any) (KeyData, error) {
	parsedPublicKey, ok := unparsedPublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("unknown public key type")
	}

	return &RSAPublicKeyData{
		data: parsedPublicKey,
	}, nil
}
