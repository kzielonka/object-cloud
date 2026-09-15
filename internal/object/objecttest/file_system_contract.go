package objecttest

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/kzielonka/object-cloud/internal/object"
)

type FileSystemFactory func(t *testing.T) object.FileSystem

func RunFileSystemContract(t *testing.T, newFS FileSystemFactory) {
	t.Run("returns ErrNotFound when file does not exist", func(t *testing.T) {
		fs := newFS(t)
		_, err := fs.OpenFile("file-id")
		if err == nil {
			t.Fatalf("expected error when path is not set, got nil")
		}
		if !errors.Is(err, object.ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("returns error when reader fails during save", func(t *testing.T) {
		fs := newFS(t)
		failingReader := iotest.ErrReader(errors.New("network stream dropped"))

		err := fs.SaveFile("broken-file", failingReader)
		if err == nil {
			t.Fatalf("expected error when reading fails, got nil")
		}
	})

	t.Run("saves and reads back file content", func(t *testing.T) {
		fs := newFS(t)
		initialData := "hello"

		err := fs.SaveFile("file-id", strings.NewReader(initialData))
		if err != nil {
			t.Fatalf("unexpected error saving: %v", err)
		}

		stream, err := fs.OpenFile("file-id")
		defer stream.Close()
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

	t.Run("isolates multiple files with different paths", func(t *testing.T) {
		fs := newFS(t)
		initialData1 := "hello 1"
		initialData2 := "hello 2"

		if err := fs.SaveFile("file-id-1", strings.NewReader(initialData1)); err != nil {
			t.Fatalf("failed to save path1: %v", err)
		}
		if err := fs.SaveFile("file-id-2", strings.NewReader(initialData2)); err != nil {
			t.Fatalf("failed to save path2: %v", err)
		}

		stream1, err := fs.OpenFile("file-id-1")
		defer stream1.Close()
		if err != nil {
			t.Fatalf("failed to open path1: %v", err)
		}
		bytes1, _ := io.ReadAll(stream1)

		stream2, err := fs.OpenFile("file-id-2")
		defer stream2.Close()
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

	t.Run("deletes existing file", func(t *testing.T) {
		fs := newFS(t)
		someData := "hello"

		err := fs.SaveFile("file-id", strings.NewReader(someData))
		if err != nil {
			t.Fatalf("unexpected error saving: %v", err)
		}

		err = fs.DeleteFile("file-id")
		if err != nil {
			t.Fatalf("unexpected error deleting: %v", err)
		}

		stream, err := fs.OpenFile("file-id")
		if stream != nil {
			defer stream.Close()
			t.Fatalf("deleted file stream should be nil")
		}
		if !errors.Is(err, object.ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("returns ErrNotFound when deleting non-existent file", func(t *testing.T) {
		fs := newFS(t)

		err := fs.DeleteFile("file-id")
		if !errors.Is(err, object.ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("returns ErrNotFound when renaming non-existent file", func(t *testing.T) {
		fs := newFS(t)

		err := fs.RenameFile("file-1-id", "file-2-id")
		if !errors.Is(err, object.ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("returns ErrFileExists when new file exists", func(t *testing.T) {
		fs := newFS(t)

		originalFile := "file-1-id"
		destFile := "file-2-id"

		err := fs.SaveFile(originalFile, strings.NewReader("some-data"))
		if err != nil {
			t.Fatalf("unexpected error saving: %v", err)
		}


		err = fs.SaveFile(destFile, strings.NewReader("some-data"))
		if err != nil {
			t.Fatalf("unexpected error saving: %v", err)
		}

		err = fs.RenameFile(originalFile, destFile)
		if !errors.Is(err, object.ErrFileExists) {
			t.Fatalf("expected ErrFileExists, got %v", err)
		}
	})

	t.Run("renames file", func(t *testing.T) {
		fs := newFS(t)

		data := "some-data-1234"
		originalFile := "file-1-id"
		destFile := "file-2-id"

		err := fs.SaveFile(originalFile, strings.NewReader(data))
		if err != nil {
			t.Fatalf("unexpected error saving: %v", err)
		}

		err = fs.RenameFile(originalFile, destFile)
		if err != nil {
			t.Fatalf("unexpected error renaming: %v", err)
		}

		stream, err := fs.OpenFile(destFile)
		if stream != nil {
			defer stream.Close()
		}
		if err != nil {
			t.Fatalf("unexpected error opening destination file: %v", err)
		}
		bytes, err := io.ReadAll(stream)
		if err != nil {
			t.Fatalf("unexpected error reading stream: %v", err)
		}
		if string(bytes) != data {
			t.Errorf("expected renamed data to be %q, got %q", data, string(bytes))
		}

		// Verify original file is gone
		origStream, err := fs.OpenFile(originalFile)
		if origStream != nil {
			defer origStream.Close()
			t.Fatalf("expected original file to no longer exist after rename")
		}
		if !errors.Is(err, object.ErrNotFound) {
			t.Errorf("expected ErrNotFound for original file after rename, got %v", err)
		}
	})
}
