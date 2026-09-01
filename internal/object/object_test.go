package object_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/kzielonka/object-cloud/internal/object"
)

func TestStore_UploadAndDownload(t *testing.T) {
	// Arrange: Set up our dependencies
	fs := object.InMemoryFileSystem()
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


