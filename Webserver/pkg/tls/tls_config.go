package tls

import "crypto/tls"

type TLSConfig struct {
	KeyProvider *KeyProvider
}

func NewTLSConfigBuilder(kp *KeyProvider) *TLSConfig {
	return &TLSConfig{
		KeyProvider: kp,
	}
}

func (cfg *TLSConfig) BuildDefaultTLSConfig() (*tls.Config, error) {
	cert, err := cfg.KeyProvider.GetCertificate()
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{*cert},
		MinVersion:   tls.VersionTLS12,
	}, nil
}
