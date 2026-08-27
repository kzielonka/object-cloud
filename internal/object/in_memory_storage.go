package object

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

var ErrNotFound = errors.New("no object")

type inMemoryStorage struct {
	savedFiles map[string][]byte
}

func InMemoryStorage() *inMemoryStorage {
	return &inMemoryStorage{
		savedFiles: make(map[string][]byte),
	}
}

func (s *inMemoryStorage) Save(key string, data io.Reader) error {
	content, err := io.ReadAll(data)
	if err != nil {
		return fmt.Errorf("failed to read data: %w", err)
	}
	s.savedFiles[key] = content
	return nil
}

func (s *inMemoryStorage) Load(key string) (io.Reader, error) {
	data, ok := s.savedFiles[key]
	if !ok {
		return nil, ErrNotFound
	}
	return bytes.NewReader(data), nil
}
