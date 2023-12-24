package cert

import (
	"crypto/rand"
	"crypto/rsa"
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
