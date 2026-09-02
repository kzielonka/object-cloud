package object

import (
	"errors"
	"io"
)

var ErrNotFound = errors.New("no object")

type FileSystem interface {
	SaveFile(path string, data io.Reader) error
	OpenFile(path string) (io.Reader, error)
}
