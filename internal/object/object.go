package object

import "io"

type Storage interface {
    Save(key string, data io.Reader) error
}

type Store struct {
    storage Storage
}

func NewStore(storage Storage) *Store {
    return &Store{storage: storage}
}

func (s *Store) Upload(key string, data io.Reader) error {
    return s.storage.Save(key, data)
}
