package service_engine

import (
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"webserver/Webserver/pkg/file_provider"
)

type ServiceEngine struct {
	WebRoot file_provider.WebRoot
	CGI     *CGIEngine
}

func NewServiceEngine(root file_provider.WebRoot, cgi *CGIEngine) *ServiceEngine {
	return &ServiceEngine{
		WebRoot: root,
		CGI:     cgi,
	}
}

func (s *ServiceEngine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cleanPath := filepath.Clean(r.URL.Path)

	if s.isCGIScript(cleanPath) {
		if s.CGI != nil {
			s.CGI.HandleCGI(w, r, cleanPath)
			return
		}
		http.Error(w, "CGI not found", http.StatusNotImplemented)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.handleGet(w, r, cleanPath)
	case http.MethodPost:
		s.handlePost(w, r, cleanPath)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (s *ServiceEngine) handleGet(w http.ResponseWriter, r *http.Request, path string) {
	file, err := s.WebRoot.Open(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer file.(io.Closer).Close()

	if seeker, ok := file.(io.ReadSeeker); ok {
		http.ServeContent(w, r, filepath.Base(path), fileModTime(file), seeker)
	} else {
		// fallback: ReadSeeker 아님
		http.Error(w, "File does not support seeking", http.StatusInternalServerError)
	}
}

func (s *ServiceEngine) handlePost(w http.ResponseWriter, r *http.Request, path string) {
	if root, ok := s.WebRoot.(*file_provider.VirtualHostWebRoot); ok {
		fullPath := filepath.Join(root.RootPath(), filepath.Clean(path))

		file, err := os.Create(fullPath)
		if err != nil {
			http.Error(w, "Failed to create file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		_, err = io.Copy(file, r.Body)
		if err != nil {
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "WebRoot does not support writing", http.StatusForbidden)
	}
}

func (s *ServiceEngine) isCGIScript(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".php" || ext == ".cgi"
}

func fileModTime(file fs.File) time.Time {
	if stat, err := file.Stat(); err == nil {
		return stat.ModTime()
	}
	return time.Now()
}
