package object

import (
	"encoding/hex"
	"errors"
	"io"
	"path/filepath"
)

var StoreError = errors.New("store error")

// Store defines the public interface for uploading and downloading objects
type Store interface {
	Upload(key string, data io.Reader) error
	Download(key string) (io.Reader, error)
}

type defaultStore struct {
	fs     FileSystem
	hasher KeyHasher
	dir    string
}

func NewStore(opts ...Option) (Store, error) {
	s := &defaultStore{}
	s.hasher = NewSHA256KeyHasher()
	for _, opt := range opts {
		opt(s)
	}
	if s.fs == nil {
		return nil, errors.New("fs must be set")
	}
	if s.dir == "" {
		return nil, errors.New("dir must be set")
	}
	return s, nil
}

func (s *defaultStore) Upload(key string, data io.Reader) error {
	path := s.pathFor(key)
	err := s.fs.SaveFile(path, data)
	if err != nil {
		return StoreError
	}
	return nil
}

func (s *defaultStore) Download(key string) (io.Reader, error) {
	path := s.pathFor(key)
	data, err := s.fs.OpenFile(path)
	if err != nil {
		return nil, StoreError
	}
	return data, nil
}

func (s *defaultStore) pathFor(key string) string {
	return filepath.Join(s.dir, hex.EncodeToString(s.hasher.Hash(key)))
}

type Option func(*defaultStore)

func WithFileSystem(fs FileSystem) Option {
	return func(s *defaultStore) {
		s.fs = fs
	}
}

func WithHasher(hasher KeyHasher) Option {
	return func(s *defaultStore) {
		if hasher != nil {
			s.hasher = hasher
		}
	}
}

func WithDir(dir string) Option {
	return func(s *defaultStore) {
		s.dir = dir
	}
}
