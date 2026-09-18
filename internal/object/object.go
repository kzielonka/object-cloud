// Package object provides domain-level abstractions and primitives
// for uploading and downloading objects backed by a pluggable FileSystem.
package object

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"path/filepath"
)

// StoreError indicates an error occurred during a store operation such as
// uploading, downloading, or reading/writing to the underlying filesystem.
var StoreError = errors.New("store error")

// Store defines the public interface for uploading and downloading objects.
// Keys are sanitized and hashed before being stored in the underlying FileSystem.
type Store interface {
	// Upload streams data into storage under the given logical key.
	Upload(key string, data io.Reader) error

	// Download returns a stream of the object stored under key.
	// The caller is responsible for closing the returned io.ReadCloser.
	Download(key string) (io.ReadCloser, error)
}

type defaultStore struct {
	fs     FileSystem
	hasher KeyHasher
	dir    string
}

// NewStore creates a new Store configured with the provided functional options.
// Both WithFileSystem and WithDir must be specified.
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
	tmpPath := s.tmpPathFor(key)
	defer s.fs.DeleteFile(tmpPath)

	err := s.fs.SaveFile(tmpPath, data)
	if err != nil {
		if err := s.fs.CreateDir(s.tmpPath()); err != nil {
			return StoreError
		}
		if err := s.fs.SaveFile(tmpPath, data); err != nil {
			return StoreError
		}
	}

	path := s.pathFor(key)
	err = s.fs.RenameFile(tmpPath, path)
	if err != nil {
		if err := s.fs.CreateDir(s.dataPath()); err != nil {
			return StoreError
		}
		if err := s.fs.RenameFile(tmpPath, path); err != nil {
			return StoreError
		}
	}
	return nil
}

func (s *defaultStore) Download(key string) (io.ReadCloser, error) {
	path := s.pathFor(key)
	data, err := s.fs.OpenFile(path)
	if err != nil {
		if data != nil {
			data.Close()
		}
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, StoreError
	}
	return data, nil
}

func (s *defaultStore) dataPath() string {
	return filepath.Join(s.dir, "data")
}

func (s *defaultStore) tmpPath() string {
	return filepath.Join(s.dir, "tmp")
}

func (s *defaultStore) pathFor(key string) string {
	return filepath.Join(s.dataPath(), hex.EncodeToString(s.hasher.Hash(key)))
}

// tmpPathFor generates a unique temporary path for staging an upload inside s.tmpPath().
func (s *defaultStore) tmpPathFor(key string) string {
	var suffix [8]byte
	_, _ = rand.Read(suffix[:])
	return filepath.Join(s.tmpPath(), hex.EncodeToString(s.hasher.Hash(key))+"."+hex.EncodeToString(suffix[:]))
}

// Option configures a Store instance.
type Option func(*defaultStore)

// WithFileSystem configures the underlying FileSystem adapter used by the Store.
func WithFileSystem(fs FileSystem) Option {
	return func(s *defaultStore) {
		s.fs = fs
	}
}

// WithHasher configures a custom KeyHasher used to hash logical keys into safe filenames.
// If not specified, NewStore defaults to using a SHA-256 key hasher.
func WithHasher(hasher KeyHasher) Option {
	return func(s *defaultStore) {
		if hasher != nil {
			s.hasher = hasher
		}
	}
}

// WithDir sets the base storage directory for the Store.
func WithDir(dir string) Option {
	return func(s *defaultStore) {
		s.dir = dir
	}
}
