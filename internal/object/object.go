package object

import "io"

type Store struct {
	fs FileSystem
}

func NewStore(fs FileSystem) *Store {
	return &Store{fs: fs}
}

func (s *Store) Upload(key string, data io.Reader) error {
	return s.fs.SaveFile(key, data)
}

func (s *Store) Download(key string) (io.Reader, error) {
	return s.fs.OpenFile(key)
}
