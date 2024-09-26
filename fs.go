package fluentfga

import (
	"os"
	"path/filepath"
)

type WriteFS interface {
	WriteFile(string, []byte) error
}

func NewWriteFS(dir string) WriteFS {
	return &writeFS{dir: dir}
}

type writeFS struct {
	dir string
}

func (w *writeFS) WriteFile(name string, contents []byte) error {
	path := filepath.Join(w.dir, name)
	return os.WriteFile(path, contents, 0644)
}
