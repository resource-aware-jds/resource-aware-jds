package cert

import (
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"errors"
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
