package tls_test

import (
	"os"
	"path/filepath"
	"testing"

	"webserver/Webserver/pkg/tls"
)

func TestTLSConfig_BuildTLSConfig(t *testing.T) {
	// 테스트용 PEM 파일 경로
	certPath := filepath.Join("testdata", "server.crt")
	keyPath := filepath.Join("testdata", "server.key")

	// 파일 존재 확인
	if _, err := os.Stat(certPath); err != nil {
		t.Fatalf("Missing certificate file: %v", err)
	}
	if _, err := os.Stat(keyPath); err != nil {
		t.Fatalf("Missing key file: %v", err)
	}

	// KeyProvider 설정
	kp := tls.NewKeyProvider(certPath, keyPath)
	kp.Loader = TestFileLoader{}

	tlsCfg := &tls.TLSConfig{
		KeyProvider: kp,
	}

	config, err := tlsCfg.BuildTLSConfig()
	if err != nil {
		t.Fatalf("BuildTLSConfig failed: %v", err)
	}

	if config == nil {
		t.Fatal("Expected non-nil *tls.Config")
	}

	if len(config.Certificates) == 0 {
		t.Fatal("Expected at least one certificate in tls.Config")
	}

	t.Logf("Successfully built tls.Config with certificate: %+v", config.Certificates[0])
}
