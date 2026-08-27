package object_test

import (
	"io"
	"strings"
	"testing"

	"github.com/kzielonka/object-cloud/internal/object"
)

func TestInMemoryStorage_ReturnErrorWhenKeyIsNotSet(t *testing.T) {
	store := object.InMemoryStorage()
	_, err := store.Load("key")
	if err == nil {
		t.Fatalf("expected error when key is not set, got nil")
	}
}

func TestInMemoryStorage_SetsKey(t *testing.T) {
	store := object.InMemoryStorage()
	initialData := "hello"

	err := store.Save("key", strings.NewReader(initialData))
	if err != nil {
		t.Fatalf("unexpected error saving: %v", err)
	}

	stream, err := store.Load("key")
	if err != nil {
		t.Fatalf("unexpected error loading: %v", err)
	}

	bytes, err := io.ReadAll(stream)
	if err != nil {
		t.Fatalf("unexpected error reading stream: %v", err)
	}

	if string(bytes) != initialData {
		t.Errorf("expected %q, got %q", initialData, string(bytes))
	}
}

func TestInMemoryStorage_SetsTwoDifferentKeys(t *testing.T) {
	store := object.InMemoryStorage()
	initialData1 := "hello 1"
	initialData2 := "hello 2"

	if err := store.Save("key1", strings.NewReader(initialData1)); err != nil {
		t.Fatalf("failed to save key1: %v", err)
	}
	if err := store.Save("key2", strings.NewReader(initialData2)); err != nil {
		t.Fatalf("failed to save key2: %v", err)
	}

	stream1, err := store.Load("key1")
	if err != nil {
		t.Fatalf("failed to load key1: %v", err)
	}
	bytes1, _ := io.ReadAll(stream1)

	stream2, err := store.Load("key2")
	if err != nil {
		t.Fatalf("failed to load key2: %v", err)
	}
	bytes2, _ := io.ReadAll(stream2)

	if string(bytes1) != initialData1 {
		t.Errorf("expected key1 to be %q, got %q", initialData1, string(bytes1))
	}
	if string(bytes2) != initialData2 {
		t.Errorf("expected key2 to be %q, got %q", initialData2, string(bytes2))
	}
}
