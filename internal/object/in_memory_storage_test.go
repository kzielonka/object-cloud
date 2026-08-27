package object_test

import (
	"io"
	"strings"
	"testing"

	"github.com/kzielonka/object-cloud/internal/object"
)

func TestInMemoryFileSystem_ReturnErrorWhenPathIsNotSet(t *testing.T) {
	fs := object.InMemoryFileSystem()
	_, err := fs.OpenFile("some/file.txt")
	if err == nil {
		t.Fatalf("expected error when path is not set, got nil")
	}
}

func TestInMemoryFileSystem_SetsPath(t *testing.T) {
	fs := object.InMemoryFileSystem()
	initialData := "hello"

	err := fs.SaveFile("some/file.txt", strings.NewReader(initialData))
	if err != nil {
		t.Fatalf("unexpected error saving: %v", err)
	}

	stream, err := fs.OpenFile("some/file.txt")
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

func TestInMemoryFileSystem_SetsTwoDifferentPaths(t *testing.T) {
	fs := object.InMemoryFileSystem()
	initialData1 := "hello 1"
	initialData2 := "hello 2"

	if err := fs.SaveFile("path1/file1.txt", strings.NewReader(initialData1)); err != nil {
		t.Fatalf("failed to save path1: %v", err)
	}
	if err := fs.SaveFile("path2/file2.txt", strings.NewReader(initialData2)); err != nil {
		t.Fatalf("failed to save path2: %v", err)
	}

	stream1, err := fs.OpenFile("path1/file1.txt")
	if err != nil {
		t.Fatalf("failed to open path1: %v", err)
	}
	bytes1, _ := io.ReadAll(stream1)

	stream2, err := fs.OpenFile("path2/file2.txt")
	if err != nil {
		t.Fatalf("failed to open path2: %v", err)
	}
	bytes2, _ := io.ReadAll(stream2)

	if string(bytes1) != initialData1 {
		t.Errorf("expected path1 to be %q, got %q", initialData1, string(bytes1))
	}
	if string(bytes2) != initialData2 {
		t.Errorf("expected path2 to be %q, got %q", initialData2, string(bytes2))
	}
}
