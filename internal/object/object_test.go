package object_test

import (
	"bytes"
	"errors"
	"github.com/kzielonka/object-cloud/internal/object"
	"io"
	"testing"
)

func TestStore_UploadAndDownload(t *testing.T) {
	// Arrange: Set up our dependencies
	fs := object.NewInMemoryFileSystem()
	store := object.NewStore(fs)

	testKey := "pets/dog-123.jpg"
	testContent := []byte("fake image content")
	reader := bytes.NewReader(testContent)

	// Act: Execute upload
	err := store.Upload(testKey, reader)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Act: Execute download
	downloadData, err := store.Download(testKey)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Assert: Verify data
	data, err := io.ReadAll(downloadData)
	if err != nil {
		t.Fatalf("unexpected error reading data: %v", err)
	}

	if !bytes.Equal(data, testContent) {
		t.Errorf("expected data %q, got %q", testContent, data)
	}
}

type fakeFileSystem struct{}

func (fs *fakeFileSystem) SaveFile(path string, data io.Reader) error {
	return errors.New("save error")
}

func (fs *fakeFileSystem) OpenFile(path string) (io.Reader, error) {
	return nil, errors.New("open error")

}

func TestStore_UploadErrorTranslation(t *testing.T) {
	// Arrange: Set up our dependencies
	fs := &fakeFileSystem{}
	store := object.NewStore(fs)

	testKey := "pets/dog-123.jpg"
	testContent := []byte("fake image content")
	reader := bytes.NewReader(testContent)

	err := store.Upload(testKey, reader)

	if !errors.Is(err, object.StoreError) {
		t.Errorf("expected StoreError, got %s", err)
	}
}

func TestStore_DownloadErrorTranslation(t *testing.T) {
	// Arrange: Set up our dependencies
	fs := &fakeFileSystem{}
	store := object.NewStore(fs)

	testKey := "pets/dog-123.jpg"

	_, err := store.Download(testKey)

	if !errors.Is(err, object.StoreError) {
		t.Errorf("expected StoreError, got %s", err)
	}
}
