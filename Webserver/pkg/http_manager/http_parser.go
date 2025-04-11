package http_manager

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"webserver/Webserver/pkg/file_provider"
	"webserver/Webserver/pkg/service_engine"
	"webserver/Webserver/pkg/types"
)

type ServiceEngine interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type HTTPManager struct {
	engine ServiceEngine
}

func NewHTTPManager() *HTTPManager {
	return &HTTPManager{}
}

func (m *HTTPManager) SetEngine(engine ServiceEngine) {
	m.engine = engine
}

func (m *HTTPManager) Init(config map[string]interface{}) error {
	// 예시: config["services"] = map[string]interface{}{"localhost": map[string]interface{}{ "webRoot": "/path/to/root" }}

	services, ok := config["services"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid or missing 'services' config")
	}

	for _, raw := range services {
		serviceConf, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		webRootPath, ok := serviceConf["webRoot"].(string)
		if !ok {
			continue
		}

		root := file_provider.NewVirtualHostWebRoot(webRootPath)
		engine := service_engine.NewServiceEngine(root, nil)
		m.SetEngine(engine)
	}

	return nil
}

func (m *HTTPManager) Handle(conn net.Conn) error {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	req, err := http.ReadRequest(reader)
	if err != nil {
		log.Println("Failed to read HTTP request:", err)
		return err
	}

	writer := NewConnResponseWriter(conn)

	if m.engine == nil {
		http.Error(writer, "Service engine not available", http.StatusInternalServerError)
		return nil
	}

	m.engine.ServeHTTP(writer, req)
	return nil
}

func (m *HTTPManager) SetNext(_ types.Filter) {
	// HTTPManager is the end of the chain; does not pass to next
}
