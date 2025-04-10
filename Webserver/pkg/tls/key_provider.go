package tls

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"webserver/Webserver/pkg/file_provider"
)

type keyProvider interface {
	GetCertificate() (*tls.Certificate, error)
}

type KeyProvider struct {
	certPath string
	keyPath  string
	pubCert  *x509.Certificate
	privKey  *rsa.PrivateKey
	Loader   file_provider.FileLoader
}

func NewKeyProvider(certPath, keyPath string) *KeyProvider {
	return &KeyProvider{
		certPath: certPath,
		keyPath:  keyPath,
		Loader:   file_provider.NewDiskFileLoader(),
	}
}

func (kp *KeyProvider) LoadCertificate() error {
	var err error
	var buf []byte
	buf, err = kp.Loader.Load(kp.certPath)
	if err != nil {
		return err
	}
	kp.pubCert, err = ParseCertFromPEM(buf)

	buf, err = kp.Loader.Load(kp.keyPath)
	if err != nil {
		return err
	}
	kp.privKey, err = ParsePrivateKeyFromPEM(buf)

	return nil
}

func (kp *KeyProvider) GetCertificate() (*tls.Certificate, error) {
	if kp.pubCert == nil || kp.privKey == nil {
		if err := kp.LoadCertificate(); err != nil {
			return nil, err
		}
	}

	cert := tls.Certificate{
		Certificate: [][]byte{kp.pubCert.Raw},
		PrivateKey:  kp.privKey,
	}
	return &cert, nil
}

func ParseCertFromPEM(certPEM []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(certPEM)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, errors.New("failed to decode PEM block containing certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

func ParsePrivateKeyFromPEM(keyPEM []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(keyPEM)
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing key")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}
