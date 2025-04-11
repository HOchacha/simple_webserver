package http_manager

import (
	"fmt"
	"net"
	"net/http"
)

type ConnResponseWriter struct {
	conn        net.Conn
	header      http.Header
	status      int
	wroteHeader bool
}

func NewConnResponseWriter(conn net.Conn) *ConnResponseWriter {
	return &ConnResponseWriter{
		conn:   conn,
		header: make(http.Header),
		status: http.StatusOK,
	}
}

func (w *ConnResponseWriter) Header() http.Header {
	return w.header
}

func (w *ConnResponseWriter) WriteHeader(statusCode int) {
	if w.wroteHeader {
		return
	}
	w.status = statusCode
	w.wroteHeader = true

	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, http.StatusText(statusCode))
	w.conn.Write([]byte(statusLine))
	w.header.Write(w.conn)
	w.conn.Write([]byte("\r\n")) // end of headers
}

func (w *ConnResponseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.conn.Write(b)
}
