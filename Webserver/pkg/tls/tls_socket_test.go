package tls_test

import (
	"crypto/tls"
	"net"
	"os"
	"testing"
	"time"
	"webserver/Webserver/pkg/filterchain"
	tlssock "webserver/Webserver/pkg/tls"
)

// mockFilter는 다음 필터가 제대로 호출되는지 확인
type mockFilter struct {
	called    bool
	isTLSConn bool
}

type mockConn struct{}

func (c *mockConn) Read(b []byte) (n int, err error)   { return 0, nil }
func (c *mockConn) Write(b []byte) (n int, err error)  { return len(b), nil }
func (c *mockConn) Close() error                       { return nil }
func (c *mockConn) LocalAddr() net.Addr                { return nil }
func (c *mockConn) RemoteAddr() net.Addr               { return nil }
func (c *mockConn) SetDeadline(t time.Time) error      { return nil }
func (c *mockConn) SetReadDeadline(t time.Time) error  { return nil }
func (c *mockConn) SetWriteDeadline(t time.Time) error { return nil }

// mockNextFilter: 다음 Filter 호출 여부를 체크
type mockNextFilter struct {
	Called   bool
	LastConn net.Conn
}

func (f *mockNextFilter) Init(config map[string]interface{}) error { return nil }
func (f *mockNextFilter) Handle(conn net.Conn) error {
	f.Called = true
	f.LastConn = conn
	return nil
}
func (f *mockNextFilter) SetNext(next filterchain.Filter) {}

func TestTLSSocketFilter_HandlesTLSAndCallsNext(t *testing.T) {
	certPEM, err := os.ReadFile("testdata/server.crt")
	if err != nil {
		t.Fatalf("failed to read cert file: %v", err)
	}
	keyPEM, err := os.ReadFile("testdata/server.key")
	if err != nil {
		t.Fatalf("failed to read key file: %v", err)
	}

	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		t.Fatalf("failed to load test cert/key: %v", err)
	}
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}

	// 연결된 net.Conn 쌍 생성
	clientConn, serverConn := net.Pipe()

	// 필터 준비
	filter := &tlssock.TLSSocketFilter{}
	filter.SetTLSConfig(tlsConfig)

	next := &mockNextFilter{}
	filter.SetNext(next)

	// 서버 쪽 핸들링 시작 (TLS Handshake 포함)
	done := make(chan error, 1)
	go func() {
		done <- filter.Handle(serverConn)
	}()

	// 클라이언트 쪽에서도 TLS Handshake 시도
	clientTLS := tls.Client(clientConn, &tls.Config{
		InsecureSkipVerify: true,
	})

	err = clientTLS.Handshake()
	if err != nil {
		t.Fatalf("TLS handshake failed: %v", err)
	}

	// 서버 핸들링 종료 대기
	select {
	case err := <-done:
		if err != nil {
			t.Logf("server Handle returned error (often OK in test): %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("server Handle timed out")
	}

	if !next.Called {
		t.Errorf("expected next filter to be called")
	}
	if _, ok := next.LastConn.(*tls.Conn); !ok {
		t.Errorf("expected a *tls.Conn to be passed to next filter")
	}
}
