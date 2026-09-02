package object_test

import (
	"io"
	"strings"
	"testing"

	"github.com/kzielonka/object-cloud/internal/object"
)

type FileSystemFactory func(t *testing.T) object.FileSystem

func RunStorageContractTests(t *testing.T, factoryFs FileSystemFactory) {
	t.Run("return error when path is not set", func(t *testing.T) {
		fs := factoryFs(t)
		_, err := fs.OpenFile("file-id")
		if err == nil {
			t.Fatalf("expected error when path is not set, got nil")
		}
	})

	t.Run("SetsPath", func(t *testing.T) {
		fs := factoryFs(t)
		initialData := "hello"

		err := fs.SaveFile("file-id", strings.NewReader(initialData))
		if err != nil {
			t.Fatalf("unexpected error saving: %v", err)
		}

		stream, err := fs.OpenFile("file-id")
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
	})

	t.Run("SetsTwoDifferentPaths", func(t *testing.T) {
		fs := factoryFs(t)
		initialData1 := "hello 1"
		initialData2 := "hello 2"

		if err := fs.SaveFile("file-id-1", strings.NewReader(initialData1)); err != nil {
			t.Fatalf("failed to save path1: %v", err)
		}
		if err := fs.SaveFile("file-id-2", strings.NewReader(initialData2)); err != nil {
			t.Fatalf("failed to save path2: %v", err)
		}

		stream1, err := fs.OpenFile("file-id-1")
		if err != nil {
			t.Fatalf("failed to open path1: %v", err)
		}
		bytes1, _ := io.ReadAll(stream1)

		stream2, err := fs.OpenFile("file-id-2")
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
	})
}
