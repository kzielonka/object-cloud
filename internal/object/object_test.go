package object_test

import (
	"bytes"
	"github.com/kzielonka/object-cloud/internal/object"
	"io"
	"testing"
)

type SpyStorage struct {
	SaveWasCalled bool
	SavedKey      string
	SavedData     []byte
}

func (s *SpyStorage) Save(key string, data io.Reader) error {
	s.SaveWasCalled = true
	s.SavedKey = key

	if data != nil {
		s.SavedData, _ = io.ReadAll(data)
	}
	return nil
}

func TestStore_Upload(t *testing.T) {
	// Arrange: Set up our dependencies
	spy := &SpyStorage{}
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
		t.Errorf("expected Storage.Save to be called, but it was not")
	}

	if spy.SavedKey != testKey {
		t.Errorf("expected key %q, got %q", testKey, spy.SavedKey)
	}

	if string(spy.SavedData) != string(testContent) {
		t.Errorf("expected data %q, got %q", string(testContent), string(spy.SavedData))
	}
}
