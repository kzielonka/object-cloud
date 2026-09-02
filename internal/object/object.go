package object

import (
	"errors"
	"io"
)

var StoreError = errors.New("store error")

type Store struct {
	fs FileSystem
}

func NewStore(fs FileSystem) *Store {
	return &Store{fs: fs}
}

func (s *Store) Upload(key string, data io.Reader) error {
	// TODO: key should be converted to single file name (hashed)
	err := s.fs.SaveFile(key, data)
	if err != nil {
		return StoreError
	}
	return nil
}

func (s *Store) Download(key string) (io.Reader, error) {
	data, err := s.fs.OpenFile(key)
	if err != nil {
		return nil, StoreError
	}
	return data, nil
}
