package cert

import (
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
)

type RAJDSPublicKey struct {
	publicKey *rsa.PublicKey
}

func ParseToRAJDSPublicKey(publicKey any) (RAJDSPublicKey, error) {
	rsaParsedPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return RAJDSPublicKey{}, errors.New("failed to parsed rsa public key")
	}

	return RAJDSPublicKey{
		publicKey: rsaParsedPublicKey,
	}, nil
}

func (r *RAJDSPublicKey) UnmarshalBSON(data []byte) error {
	publicKey, err := x509.ParsePKCS1PublicKey(data)
	if err != nil {
		return err
	}

	r.publicKey = publicKey
	return nil
}

func (r *RAJDSPublicKey) MarshalBSON() ([]byte, error) {
	return x509.MarshalPKCS1PublicKey(r.publicKey), nil
}

func (r *RAJDSPublicKey) GetSHA1Hash() string {
	sha1Hash := sha1.Sum(x509.MarshalPKCS1PublicKey(r.publicKey))
	return hex.EncodeToString(sha1Hash[:])
}

func (r *RAJDSPublicKey) GetPublicKey() *rsa.PublicKey {
	return r.publicKey
}

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
	parsedPublicKey, ok := unparsedPublicKey.(rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("unknown public key type")
	}

	return &RSAPublicKeyData{
		data: parsedPublicKey,
	}, nil
}
