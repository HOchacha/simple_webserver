package file_provider

import "os"

type FileLoader interface {
	Load(path string) ([]byte, error)
}

type DiskFileLoader struct {
}

func (fl *DiskFileLoader) Load(path string) ([]byte, error) {
	return os.ReadFile(path)
}
