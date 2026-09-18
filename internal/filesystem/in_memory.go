package filesystem

import (
	"bytes"
	"errors"
	"io"
	"path/filepath"

	"github.com/kzielonka/object-cloud/internal/object"
)

type inMemoryFileSystem struct {
	savedFiles         map[string][]byte
	createdDirectories map[string]struct{}
}

func NewInMemory() *inMemoryFileSystem {
	return &inMemoryFileSystem{
		savedFiles:         make(map[string][]byte),
		createdDirectories: make(map[string]struct{}),
	}
}

func (s *inMemoryFileSystem) SaveFile(path string, data io.Reader) error {
	dir := filepath.Dir(path)
	if dir != "." && dir != "/" && !s.HasDirectory(dir) {
		return errors.New("directory does not exist")
	}

	content, err := io.ReadAll(data)
	if err != nil {
		return err
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
	dir := filepath.Dir(newPath)
	if dir != "." && dir != "/" && !s.HasDirectory(dir) {
		return object.ErrNotFound
	}
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

func (s *inMemoryFileSystem) CreateDir(path string) error {
	s.createdDirectories[path] = struct{}{}
	return nil
}

func (s *inMemoryFileSystem) HasDirectory(path string) bool {
	_, ok := s.createdDirectories[path]
	return ok
}
