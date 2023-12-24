package cert

import (
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
)

type RSAPublicKeyData struct {
	data rsa.PublicKey
}

func (r *RSAPublicKeyData) GetKeyType() KeyType {
	return PKIPublicKeyType
}

func (r *RSAPublicKeyData) GetSHA1Hash() (string, error) {
	sha1Hash := sha1.Sum(x509.MarshalPKCS1PublicKey(&r.data))
	return hex.EncodeToString(sha1Hash[:]), nil
}

func (r *RSAPublicKeyData) GetKeyX509Format() ([]byte, error) {
	return x509.MarshalPKCS1PublicKey(&r.data), nil
}

func (r *RSAPublicKeyData) GetRawKeyData() any {
	return r.data
}

type RSAPrivateKeyData struct {
	data *rsa.PrivateKey
}

func (r *RSAPrivateKeyData) GetKeyType() KeyType {
	return PKIPrivateKeyType
}

func (r *RSAPrivateKeyData) GetSHA1Hash() (string, error) {
	sha1Hash := sha1.Sum(x509.MarshalPKCS1PrivateKey(r.data))
	return hex.EncodeToString(sha1Hash[:]), nil
}

func (r *RSAPrivateKeyData) GetKeyX509Format() ([]byte, error) {
	return x509.MarshalPKCS1PrivateKey(r.data), nil
}

func (r *RSAPrivateKeyData) GetRawKeyData() any {
	return r.data
}
