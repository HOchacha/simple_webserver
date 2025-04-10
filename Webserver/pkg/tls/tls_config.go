package tls

import "crypto/tls"

type TLSConfig struct {
	KeyProvider *KeyProvider
}

func (cfg *TLSConfig) BuildTLSConfig() (*tls.Config, error) {
	cert, err := cfg.KeyProvider.GetCertificate()
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{*cert},
		MinVersion:   tls.VersionTLS12,
	}, nil
}
