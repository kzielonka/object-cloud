package object

import (
	"errors"
	"io"
)

var ErrNotFound = errors.New("no object")

type FileSystem interface {
	SaveFile(path string, data io.Reader) error
	OpenFile(path string) (io.ReadCloser, error)
	DeleteFile(path string) error
	RenameFile(oldPath string, newPath string) error
}
