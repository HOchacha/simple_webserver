package file_provider

import (
	"io/fs"
	"os"
	"path/filepath"
)

type WebRoot interface {
	Open(path string) (fs.File, error)
}

type VirtualHostWebRoot struct {
	webRootPath string
}

func NewVirtualHostWebRoot(webRootPath string) *VirtualHostWebRoot {
	return &VirtualHostWebRoot{webRootPath: webRootPath}
}

// return file descriptor by enclosing with fs.File
func (v *VirtualHostWebRoot) Open(path string) (fs.File, error) {
	// remove relational path
	cleanPath := filepath.Clean(path)
	fullPath := filepath.Join(v.webRootPath, cleanPath)
	return os.Open(fullPath)
}

func (v *VirtualHostWebRoot) RootPath() string {
	return v.webRootPath
}
