package tls_test

import (
	"os"
	"path/filepath"
	"testing"
	"webserver/Webserver/pkg/tls"
)

type TestFileLoader struct{}

func (l TestFileLoader) Load(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func TestKeyProvider_LoadCertificate(t *testing.T) {
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

	kp := tls.NewKeyProvider(certPath, keyPath)
	kp.Loader = TestFileLoader{}

	err := kp.LoadCertificate()
	if err != nil {
		t.Fatalf("LoadCertificate failed: %v", err)
	}

	tlsCert, err := kp.GetCertificate()
	if err != nil {
		t.Fatalf("GetCertificate returned error: %v", err)
	}

	if tlsCert == nil || len(tlsCert.Certificate) == 0 {
		t.Fatalf("Invalid tls.Certificate returned")
	}
}
