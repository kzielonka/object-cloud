package object

import (
	"bytes"
	"fmt"
	"io"
)

type inMemoryFileSystem struct {
	savedFiles map[string][]byte
}

func NewInMemoryFileSystem() *inMemoryFileSystem {
	return &inMemoryFileSystem{
		savedFiles: make(map[string][]byte),
	}
}

func (s *inMemoryFileSystem) SaveFile(path string, data io.Reader) error {
	content, err := io.ReadAll(data)
	if err != nil {
		return fmt.Errorf("failed to read data: %w", err)
	}
	s.savedFiles[path] = content
	return nil
}

func (s *inMemoryFileSystem) OpenFile(path string) (io.Reader, error) {
	data, ok := s.savedFiles[path]
	if !ok {
		return nil, ErrNotFound
	}
	return bytes.NewReader(data), nil
}
