package object_test

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/kzielonka/object-cloud/internal/filesystem"
	"github.com/kzielonka/object-cloud/internal/object"
)

func TestStore_UploadAndDownload(t *testing.T) {
	// Arrange: Set up our dependencies
	store, err := object.NewStore(
		object.WithFileSystem(filesystem.NewInMemory()),
		object.WithDir("/test"),
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	testKey := "pets/dog-123.jpg"
	testContent := []byte("fake image content")
	reader := bytes.NewReader(testContent)

	// Act: Execute upload
	err = store.Upload(testKey, reader)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Act: Execute download
	downloadData, err := store.Download(testKey)
	if downloadData != nil {
		defer downloadData.Close()
	}
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

func TestStore_FailedUploadAtomicity(t *testing.T) {
	// Arrange: Set up our dependencies
	store, err := object.NewStore(
		object.WithFileSystem(filesystem.NewInMemory()),
		object.WithDir("/test"),
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	testKey := "pets/dog-123.jpg"
	reader := io.MultiReader(
		strings.NewReader("some data"),
		iotest.ErrReader(errors.New("some error")),
	)

	// Act: Execute upload
	err = store.Upload(testKey, reader)
	if err == nil {
		t.Fatalf("expected error when upload fails, got nil")
	}
	if !errors.Is(err, object.StoreError) {
		t.Fatalf("expected StoreError, got %v", err)
	}

	// Act: Execute download
	downloadData, err := store.Download(testKey)
	if downloadData != nil {
		defer downloadData.Close()
	}
	if err == nil {
		t.Fatalf("expected error when download fails, got nil")
	}
	if !errors.Is(err, object.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for aborted upload %q, got %v", testKey, err)
	}
}

type fakeFileSystem struct{}

func (fs *fakeFileSystem) SaveFile(path string, data io.Reader) error {
	return errors.New("save error")
}

func (fs *fakeFileSystem) OpenFile(path string) (io.ReadCloser, error) {
	return nil, errors.New("open error")
}

func (fs *fakeFileSystem) DeleteFile(path string) error {
	return errors.New("delete error")
}

func (fs *fakeFileSystem) RenameFile(oldPath string, newPath string) error {
	return errors.New("rename error")
}

func (fs *fakeFileSystem) CreateDir(path string) error {
	return errors.New("create dir error")
}

func TestStore_UploadErrorTranslation(t *testing.T) {
	// Arrange: Set up our dependencies
	store, err := object.NewStore(
		object.WithFileSystem(&fakeFileSystem{}),
		object.WithDir("/test"),
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	testKey := "pets/dog-123.jpg"
	testContent := []byte("fake image content")
	reader := bytes.NewReader(testContent)

	err = store.Upload(testKey, reader)

	if !errors.Is(err, object.StoreError) {
		t.Errorf("expected StoreError, got %s", err)
	}
}

func TestStore_DownloadErrorTranslation(t *testing.T) {
	// Arrange: Set up our dependencies
	store, err := object.NewStore(
		object.WithFileSystem(&fakeFileSystem{}),
		object.WithDir("/test"),
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	testKey := "pets/dog-123.jpg"

	_, err = store.Download(testKey)

	if !errors.Is(err, object.StoreError) {
		t.Errorf("expected StoreError, got %s", err)
	}
}

type trackingReadCloser struct {
	closed bool
}

func (r *trackingReadCloser) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (r *trackingReadCloser) Close() error {
	r.closed = true
	return nil
}

type errorWithFileFileSystem struct {
	file io.ReadCloser
	err  error
}

func (fs *errorWithFileFileSystem) SaveFile(path string, data io.Reader) error {
	return nil
}

func (fs *errorWithFileFileSystem) OpenFile(path string) (io.ReadCloser, error) {
	return fs.file, fs.err
}

func (fs *errorWithFileFileSystem) DeleteFile(path string) error {
	return nil
}

func (fs *errorWithFileFileSystem) RenameFile(oldPath string, newPath string) error {
	return nil
}

func (fs *errorWithFileFileSystem) CreateDir(path string) error {
	return nil
}

func TestStore_DownloadClosesFileOnError(t *testing.T) {
	trackingFile := &trackingReadCloser{}
	fakeFS := &errorWithFileFileSystem{
		file: trackingFile,
		err:  errors.New("open error with non-nil file"),
	}

	store, err := object.NewStore(
		object.WithFileSystem(fakeFS),
		object.WithDir("/test"),
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = store.Download("pets/dog-123.jpg")
	if !errors.Is(err, object.StoreError) {
		t.Fatalf("expected StoreError, got %v", err)
	}

	if !trackingFile.closed {
		t.Errorf("expected file to be closed when OpenFile returns an error with non-nil data")
	}
}

func TestStore_UploadCreatesDirectoriesLazily(t *testing.T) {
	fs := filesystem.NewInMemory()
	store, err := object.NewStore(
		object.WithFileSystem(fs),
		object.WithDir("/test"),
	)
	if err != nil {
		t.Fatalf("expected store, got %v", err)
	}
	if fs.HasDirectory("/test/tmp") {
		t.Fatalf("expected no tmp directory yet")
	}
	if fs.HasDirectory("/test/data") {
		t.Fatalf("expected no data directory yet")
	}

	// First upload: triggers lazy directory creation
	err = store.Upload("test-key", strings.NewReader("some data"))
	if err != nil {
		t.Fatalf("expected successful upload, got %v", err)
	}
	if !fs.HasDirectory("/test/tmp") {
		t.Fatalf("expected tmp directory to be created")
	}
	if !fs.HasDirectory("/test/data") {
		t.Fatalf("expected data directory to be created")
	}

	// Subsequent upload: succeeds without issue when directories already exist
	err = store.Upload("test-key-2", strings.NewReader("other data"))
	if err != nil {
		t.Fatalf("expected subsequent upload to succeed, got %v", err)
	}
}
