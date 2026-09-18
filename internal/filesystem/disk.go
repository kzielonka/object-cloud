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

func (s *diskFileSystem) DeleteFile(path string) error {
	fullPath := filepath.Join(s.dirPath, path)
	err := os.Remove(fullPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return object.ErrNotFound
		}
		return err
	}
	return nil
}

// RenameFile moves oldPath to newPath without overwriting an existing newPath.
// Note: The check-then-rename prevents accidental overwrites, but is not fully
// atomic against concurrent writes (TOCTOU race). In V2, consider atomic
// overwrites or explicit object versioning / immutability.
func (s *diskFileSystem) RenameFile(oldPath string, newPath string) error {
	fullOldPath := filepath.Join(s.dirPath, oldPath)
	fullNewPath := filepath.Join(s.dirPath, newPath)
	_, err := os.Stat(fullNewPath)
	if err == nil {
		return object.ErrFileExists
	}
	err = os.Rename(fullOldPath, fullNewPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return object.ErrNotFound
		}
		return err
	}
	return nil
}

func (s *diskFileSystem) CreateDir(path string) error {
	fullPath := filepath.Join(s.dirPath, path)
	return os.MkdirAll(fullPath, 0755)
}
