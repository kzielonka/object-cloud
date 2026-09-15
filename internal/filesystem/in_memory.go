package filesystem

import (
	"bytes"
	"fmt"
	"io"

	"github.com/kzielonka/object-cloud/internal/object"
)

type inMemoryFileSystem struct {
	savedFiles map[string][]byte
}

func NewInMemory() *inMemoryFileSystem {
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

func (s *inMemoryFileSystem) OpenFile(path string) (io.ReadCloser, error) {
	data, ok := s.savedFiles[path]
	if !ok {
		return nil, object.ErrNotFound
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}
