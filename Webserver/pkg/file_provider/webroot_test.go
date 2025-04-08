package file_provider_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"
	"webserver/Webserver/pkg/file_provider"
)

func setupTestFile(t *testing.T, dir, filename, content string) string {
	t.Helper()

	fullPath := filepath.Join(dir, filename)
	err := os.MkdirAll(filepath.Dir(fullPath), 0755)
	if err != nil {
		t.Fatalf("failed to create dir: %v", err)
	}

	err = os.WriteFile(fullPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	return fullPath
}

func TestVirtualHostWebRoot_Open_ValidFile(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestFile(t, tmpDir, "test.txt", "Hello, WebRoot!")

	v := file_provider.NewVirtualHostWebRoot(tmpDir)
	f, err := v.Open("test.txt")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	if string(data) != "Hello, WebRoot!" {
		t.Errorf("unexpected file content: %s", data)
	}
}

func TestVirtualHostWebRoot_Open_FileNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	v := file_provider.NewVirtualHostWebRoot(tmpDir)

	_, err := v.Open("notfound.txt")
	if err == nil {
		t.Errorf("expected error for non-existent file, got nil")
	}
}

func TestVirtualHostWebRoot_Open_PathTraversal(t *testing.T) {
	tmpDir := t.TempDir()

	// 시도: 루트 밖의 파일에 접근 (비정상)
	v := file_provider.NewVirtualHostWebRoot(tmpDir)
	_, err := v.Open("../etc/passwd")
	if err == nil {
		t.Errorf("expected error for path traversal, got nil")
	}
}
