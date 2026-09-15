package filesystem

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/kzielonka/object-cloud/internal/object"
)

type diskFileSystem struct {
	dirPath string
}

func NewDisk(path string) *diskFileSystem {
	return &diskFileSystem{
		dirPath: path,
	}
}

func (s *diskFileSystem) SaveFile(path string, data io.Reader) error {
	fullPath := filepath.Join(s.dirPath, path)
	outFile, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, data); err != nil {
		return err
	}
	return nil
}

func (s *diskFileSystem) OpenFile(path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(s.dirPath, path)
	file, err := os.Open(fullPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, object.ErrNotFound
		}
		return nil, err
	}

	return file, nil
}
