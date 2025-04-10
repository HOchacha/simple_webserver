package service_engine_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"webserver/Webserver/pkg/file_provider"
	"webserver/Webserver/pkg/service_engine"
)

// 임시 WebRoot 디렉토리 생성
func setupTempWebRoot(t *testing.T) string {
	dir := t.TempDir()
	indexPath := filepath.Join(dir, "index.html")
	err := os.WriteFile(indexPath, []byte("<html><body>Hello World</body></html>"), 0644)
	if err != nil {
		t.Fatalf("파일 생성 실패: %v", err)
	}
	return dir
}

func TestServiceEngine_GET(t *testing.T) {
	webRootPath := setupTempWebRoot(t)
	webRoot := file_provider.NewVirtualHostWebRoot(webRootPath)

	engine := service_engine.NewServiceEngine(webRoot, nil)

	req := httptest.NewRequest(http.MethodGet, "/index.html", nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("원하는 상태 코드 아님: got %d", resp.StatusCode)
	}

	if !strings.Contains(string(body), "Hello World") {
		t.Errorf("본문이 예상과 다름: %s", string(body))
	}
}

func TestServiceEngine_POST(t *testing.T) {
	webRootPath := setupTempWebRoot(t)
	webRoot := file_provider.NewVirtualHostWebRoot(webRootPath)

	engine := service_engine.NewServiceEngine(webRoot, nil)

	req := httptest.NewRequest(http.MethodPost, "/upload.txt", strings.NewReader("Uploaded content"))
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("POST 실패: 상태 코드 %d", resp.StatusCode)
	}

	// 실제로 파일이 생성되었는지 확인
	uploadedPath := filepath.Join(webRootPath, "upload.txt")
	content, err := os.ReadFile(uploadedPath)
	if err != nil {
		t.Fatalf("업로드 파일 읽기 실패: %v", err)
	}

	if string(content) != "Uploaded content" {
		t.Errorf("업로드 내용 불일치: %s", string(content))
	}
}
