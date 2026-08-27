package object_test

import (
	"bytes"
	"github.com/kzielonka/object-cloud/internal/object"
	"io"
	"testing"
)

type SpyFileSystem struct {
	SaveWasCalled bool
	SavedPath     string
	SavedData     []byte
}

func (s *SpyFileSystem) SaveFile(path string, data io.Reader) error {
	s.SaveWasCalled = true
	s.SavedPath = path

	if data != nil {
		s.SavedData, _ = io.ReadAll(data)
	}
	return nil
}

func (s *SpyFileSystem) OpenFile(path string) (io.Reader, error) {
	s.SavedPath = path
	return bytes.NewReader(s.SavedData), nil
}

func TestStore_Upload(t *testing.T) {
	// Arrange: Set up our dependencies
	spy := &SpyFileSystem{}
	store := object.NewStore(spy)

	testKey := "pets/dog-123.jpg"
	testContent := []byte("fake image content")
	reader := bytes.NewReader(testContent)

	// Act: Execute the method we are testing
	err := store.Upload(testKey, reader)

	// Assert: Verify the behavior
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !spy.SaveWasCalled {
		t.Errorf("expected FileSystem.SaveFile to be called, but it was not")
	}

	if spy.SavedPath != testKey {
		t.Errorf("expected path %q, got %q", testKey, spy.SavedPath)
	}

	if string(spy.SavedData) != string(testContent) {
		t.Errorf("expected data %q, got %q", string(testContent), string(spy.SavedData))
	}
}

func TestStore_Download(t *testing.T) {
	testKey := "pets/dog-123.jpg"
	testContent := []byte("fake image content")
	spy := &SpyFileSystem{SavedPath: testKey, SavedData: testContent}
	store := object.NewStore(spy)

	reader, err := store.Download(testKey)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("unexpected error reading data: %v", err)
	}

	if string(data) != string(testContent) {
		t.Errorf("expected data %q, got %q", string(testContent), string(data))
	}
}
