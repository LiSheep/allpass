package helper

import (
	"io/ioutil"
	"encoding/pem"
	"crypto/x509"
	"crypto/rsa"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"crypto/hmac"
	"crypto/sha256"
)

type Credentials struct {
	publicKey []byte
	privateKey []byte
	secretBASE64 []byte

	rsaPubKey *rsa.PublicKey
	rsaPrvKey *rsa.PrivateKey
}

func NewCredentials(pubKeyPath, prvKeyPath string) (c *Credentials, err error) {
	c = new(Credentials)
	c.publicKey, err = ioutil.ReadFile(pubKeyPath)
	if err != nil {
		return
	}
	block, _ := pem.Decode(c.publicKey)
	if block == nil {
		return
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return
	}
	c.rsaPubKey = pubInterface.(*rsa.PublicKey)

	c.privateKey, err = ioutil.ReadFile(prvKeyPath)
	if err != nil {
		return
	}
	block, _ = pem.Decode(c.privateKey)
	if block == nil {
		return
	}
	c.rsaPrvKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return
	}
	c.secretBASE64 = make([]byte, base64.StdEncoding.EncodedLen(len(block.Bytes)))
	base64.StdEncoding.Encode(c.secretBASE64, block.Bytes)
	return
}

func (c *Credentials) Encode(data []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, c.rsaPubKey, data)
}

func (c *Credentials) Decode(data []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, c.rsaPrvKey, data)
}

func (c *Credentials) EncodeSHA512(data []byte) (result []byte) {
	h := sha512.New()
	h.Write(data)
	h.Write(c.secretBASE64)
	result = h.Sum(nil)
	return
}

func (c *Credentials) EncodeSHA256(data []byte) (result []byte) {
	h := sha256.New()
	h.Write(data)
	h.Write(c.secretBASE64)
	result = h.Sum(nil)
	return
}

func (c *Credentials) EncodeHmacSha512(data []byte) ([] byte) {
	mac := hmac.New(sha512.New, c.secretBASE64)
	mac.Write(data)
	return mac.Sum(nil)
}

func (c *Credentials) EncodeHmacSha256(data []byte) ([] byte) {
	mac := hmac.New(sha256.New, c.secretBASE64)
	mac.Write(data)
	return mac.Sum(nil)
}