package object

import "io"

type FileSystem interface {
	SaveFile(path string, data io.Reader) error
	OpenFile(path string) (io.Reader, error)
}
