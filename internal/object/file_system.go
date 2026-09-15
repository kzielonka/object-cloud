package object

import (
	"errors"
	"io"
)

var ErrNotFound = errors.New("no object")
var ErrFileExists = errors.New("file exists")

type FileSystem interface {
	SaveFile(path string, data io.Reader) error
	OpenFile(path string) (io.ReadCloser, error)
	DeleteFile(path string) error
	RenameFile(oldPath string, newPath string) error
}
