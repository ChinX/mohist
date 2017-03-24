package crypts

import (
	"bytes"
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"io"
)

// MD5 creates md5 string for an input key
func MD5(input string) string {
	h := md5.New()
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}

// SHA512 creates sha512 string for an input bytes
func SHA512(byteArr []byte) (string, error) {
	return SHA512Stream(bytes.NewReader(byteArr))
}

// SHA512 creates sha512 string for an input reader
func SHA512Stream(reader io.Reader) (string, error) {
	hashed := sha512.New()
	if _, err := io.Copy(hashed, reader); err != nil {
		return "", err
	}
	return hex.EncodeToString(hashed.Sum(nil)), nil
}

// GenerateRSAKeyPair generate a private key and a public key
func GenerateRSAKeyPair(bits int) ([]byte, []byte, error) {
	prvKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}

	prvBytes := x509.MarshalPKCS1PrivateKey(prvKey)
	pubBytes, err := x509.MarshalPKIXPublicKey(&prvKey.PublicKey)

	if err != nil {
		return nil, nil, err
	}

	prvBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: prvBytes}
	pubBlock := &pem.Block{Type: "RSA PUBLIC KEY", Bytes: pubBytes}

	return pem.EncodeToMemory(prvBlock), pem.EncodeToMemory(pubBlock), nil
}

// RSAEncrypt encrypts a content by a public key
func RSAEncrypt(keyBytes, contentBytes []byte) ([]byte, error) {
	pubKey, err := getPubKey(keyBytes)
	if err != nil {
		return nil, err
	}
	return rsa.EncryptPKCS1v15(rand.Reader, pubKey, contentBytes)
}

// RSADecrypt decrypts a content by a private key
func RSADecrypt(keyBytes, contentBytes []byte) ([]byte, error) {
	prvKey, err := getPrvKey(keyBytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, prvKey, contentBytes)
}

// SHA256Sign signs a content by a private key
func SHA256Sign(keyBytes, contentBytes []byte) ([]byte, error) {
	prvKey, err := getPrvKey(keyBytes)
	if err != nil {
		return nil, err
	}

	hashed := sha256.Sum256(contentBytes)
	return rsa.SignPKCS1v15(rand.Reader, prvKey, crypto.SHA256, hashed[:])
}

// SHA256Verify verifies if a content is valid by a signed data and a public key
func SHA256Verify(keyBytes, contentBytes, signBytes []byte) error {
	pubKey, err := getPubKey(keyBytes)
	if err != nil {
		return err
	}

	signStr := hex.EncodeToString(signBytes)
	newSignBytes, _ := hex.DecodeString(signStr)
	hashed := sha256.Sum256(contentBytes)
	return rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed[:], newSignBytes)
}

func getPrvKey(prvBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(prvBytes)
	if block == nil {
		return nil, errors.New("Fail to decode private key")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
func getPubKey(pubBytes []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pubBytes)
	if block == nil {
		return nil, errors.New("Fail to decode public key")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	pubKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("Fail get public key form pulic interface")
	}

	return pubKey, nil
}
