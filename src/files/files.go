package files

import (
	"os"
	"path/filepath"
)

type File struct {
	Name    string
	Content []byte

	Path string
}

func NewFile(path string) (File, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return File{}, err
	}

	return File{
		Name:    filepath.Base(path),
		Content: bytes,

		Path: filepath.Clean(path),
	}, nil
}
