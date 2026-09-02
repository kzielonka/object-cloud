package object

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
)

type diskFileSystem struct {
	dirPath string
}

func NewDiskFileSystem(path string) *diskFileSystem {
	return &diskFileSystem{
		dirPath: path,
	}
}

func (s *diskFileSystem) SaveFile(path string, data io.Reader) error {
	fullPath := filepath.Join(s.dirPath, path)
	outFile, _ := os.Create(fullPath)
	defer outFile.Close()
	// TODO: handle error
	io.Copy(outFile, data)
	return nil
}

func (s *diskFileSystem) OpenFile(path string) (io.Reader, error) {
	fullPath := filepath.Join(s.dirPath, path)
	file, err := os.Open(fullPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	defer file.Close()

	data, _ := io.ReadAll(file)

	// it should return reader and close method
	return bytes.NewReader(data), nil
}
