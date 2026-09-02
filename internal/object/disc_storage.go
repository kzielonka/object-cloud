package object

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type discFileSystem struct {
	dirPath string
}

func DiscFileSystem(path string) *discFileSystem {
	return &discFileSystem{
		dirPath: path,
	}
}

func (s *discFileSystem) SaveFile(path string, data io.Reader) error {
	fullPath := filepath.Join(s.dirPath, path)
	outFile, _ := os.Create(fullPath)
	defer outFile.Close()
	// TODO: handle error
	io.Copy(outFile, data)
	return nil
}

func (s *discFileSystem) OpenFile(path string) (io.Reader, error) {
	fullPath := filepath.Join(s.dirPath, path)
	file, err := os.Open(fullPath)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, ErrNotFound
	}
	defer file.Close()

	data, _ := io.ReadAll(file)

	// it should return reader and close method
	return bytes.NewReader(data), nil
}
