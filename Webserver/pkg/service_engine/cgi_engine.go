package service_engine

import (
	"bufio"
	"fmt"
	"github.com/tomasen/fcgi_client"
	"io"
	"net"
	"net/http"
	"path/filepath"
	"strings"
)

type CGIEngine struct {
	ScriptRoot    string
	FastCGISocket string // e.g. "127.0.0.1:9000" or "/run/php/php-fpm.sock"
	IsUnixSocket  bool   // true = unix socket, false = tcp
}

func NewCGIEngine(scriptRoot, fastcgiSocket string, isUnix bool) *CGIEngine {
	return &CGIEngine{
		ScriptRoot:    scriptRoot,
		FastCGISocket: fastcgiSocket,
		IsUnixSocket:  isUnix,
	}
}

func (c *CGIEngine) HandleCGI(w http.ResponseWriter, r *http.Request, scriptPath string) {
	scriptFullPath := filepath.Join(c.ScriptRoot, filepath.Clean(scriptPath))

	var conn *fcgiclient.FCGIClient
	var err error
	if c.IsUnixSocket {
		conn, err = fcgiclient.Dial("unix", c.FastCGISocket)
	} else {
		conn, err = fcgiclient.Dial("tcp", c.FastCGISocket)
	}
	if err != nil {
		http.Error(w, "FastCGI 연결 실패: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	params := map[string]string{
		"SCRIPT_FILENAME":   scriptFullPath,
		"REQUEST_METHOD":    r.Method,
		"QUERY_STRING":      r.URL.RawQuery,
		"CONTENT_TYPE":      r.Header.Get("Content-Type"),
		"CONTENT_LENGTH":    fmt.Sprintf("%d", r.ContentLength),
		"SCRIPT_NAME":       scriptPath,
		"REQUEST_URI":       r.RequestURI,
		"DOCUMENT_ROOT":     c.ScriptRoot,
		"GATEWAY_INTERFACE": "CGI/1.1",
		"SERVER_SOFTWARE":   "GoServiceEngine/1.0",
		"SERVER_PROTOCOL":   r.Proto,
		"REMOTE_ADDR":       parseRemoteAddr(r.RemoteAddr),
	}

	// 헤더를 FastCGI에 전달
	for key, vals := range r.Header {
		headerKey := "HTTP_" + strings.ReplaceAll(strings.ToUpper(key), "-", "_")
		params[headerKey] = strings.Join(vals, ",")
	}

	var stdin io.Reader
	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		stdin = r.Body
	}

	resp, err := conn.Request(params, stdin)
	if err != nil {
		http.Error(w, "FastCGI 요청 실패: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	br := bufio.NewReader(resp.Body)
	headers := http.Header{}
	for {
		line, err := br.ReadString('\n')
		if err != nil {
			http.Error(w, "FastCGI 응답 오류", http.StatusInternalServerError)
			return
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		headers.Add(key, val)
	}

	if statusLine := headers.Get("Status"); statusLine != "" {
		headers.Del("Status")
		statusParts := strings.SplitN(statusLine, " ", 2)
		code := 200
		if len(statusParts) > 0 {
			fmt.Sscanf(statusParts[0], "%d", &code)
		}
		w.WriteHeader(code)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	for k, vv := range headers {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}

	io.Copy(w, br)
}

func parseRemoteAddr(addr string) string {
	if host, _, err := net.SplitHostPort(addr); err == nil {
		return host
	}
	return addr
}
