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

func (s *inMemoryFileSystem) DeleteFile(path string) error {
	_, ok := s.savedFiles[path]
	if !ok {
		return object.ErrNotFound
	}
	delete(s.savedFiles, path)
	return nil
}

// RenameFile moves oldPath to newPath without overwriting an existing newPath.
// Note: In V2, consider atomic overwrites or explicit object versioning / immutability.
func (s *inMemoryFileSystem) RenameFile(oldPath string, newPath string) error {
	_, ok := s.savedFiles[oldPath]
	if !ok {
		return object.ErrNotFound
	}
	_, ok = s.savedFiles[newPath]
	if ok {
		return object.ErrFileExists
	}

	s.savedFiles[newPath] = s.savedFiles[oldPath]
	delete(s.savedFiles, oldPath)

	return nil
}
