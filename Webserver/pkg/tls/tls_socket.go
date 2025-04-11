package tls

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"webserver/Webserver/pkg/types"
)

type TLSSocket struct {
	cSocket net.Conn
	reader  io.Reader
	writer  io.Writer
}

type TLSSocketFilter struct {
	config *tls.Config
	next   types.Filter
}

func (f *TLSSocketFilter) Init(config map[string]interface{}) error {
	// read certPath and keyPath for loading them via KeyProvider
	certPath := config["cert_file"].(string)
	keyPath := config["key_file"].(string)

	kp := NewKeyProvider(certPath, keyPath)

	tlsCfg := NewTLSConfigBuilder(kp)
	var err error
	f.config, err = tlsCfg.BuildDefaultTLSConfig()
	if err != nil {
		return fmt.Errorf("failed to build TLS config: %w", err)
	}
	return nil
}

func (f *TLSSocketFilter) Handle(conn net.Conn) error {
	tlsConn := tls.Server(conn, f.config)
	if err := tlsConn.Handshake(); err != nil {
		log.Printf("[TLSFilter] TLS Handshake failed: %v", err)
		return err
	}
	// 다음 필터로 전달
	return f.next.Handle(tlsConn)
}

func (f *TLSSocketFilter) SetNext(next types.Filter) {
	f.next = next
}

func (f *TLSSocketFilter) SetTLSConfig(cfg *tls.Config) {
	f.config = cfg
}

// Legacy Codes
func BuildTLSSocket(conn net.Conn) *TLSSocket {
	return &TLSSocket{
		cSocket: conn,
		reader:  bufio.NewReader(conn),
	}
}

func handleTLSConn(rawConn net.Conn, config *tls.Config) (*tls.Conn, error) {
	tlsConn := tls.Server(rawConn, config)

	err := tlsConn.Handshake()
	if err != nil {
		log.Printf("TLS Handshake failed: %v", err)
		return nil, err
	}
	return tlsConn, nil
}

func isTLSClientHello(data []byte) bool {
	return len(data) >= 3 && data[0] == 0x16 && data[1] == 0x03 && (data[2] >= 0x01 && data[2] <= 0x04)
}

func (t *TLSSocket) ReadRequest() (*http.Request, error) {
	return http.ReadRequest(bufio.NewReader(t.cSocket))
}

func (t *TLSSocket) WriteResponse(resp *http.Response) error {
	return resp.Write(t.cSocket)
}

func (t *TLSSocket) Close() error {
	return t.cSocket.Close()
}
